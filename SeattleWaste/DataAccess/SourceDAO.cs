using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using SeattleWaste.Models.Shared;
using SeattleWaste.Models.Source;
using TwoMQTT.Interfaces;

namespace SeattleWaste.DataAccess
{
    public interface ISourceDAO : ISourceDAO<SlugMapping, Response, object, object>
    {
    }

    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceDAO : ISourceDAO
    {
        /// <summary>
        /// Initializes a new instance of the SourceDAO class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="httpClientFactory"></param>
        /// <returns></returns>
        public SourceDAO(ILogger<SourceDAO> logger, IHttpClientFactory httpClientFactory)
        {
            this.Logger = logger;
            this.Client = httpClientFactory.CreateClient();
        }

        /// <summary>
        /// Fetch one response from the source.
        /// </summary>
        /// <param name="key"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        public async Task<Response?> FetchOneAsync(SlugMapping key,
            CancellationToken cancellationToken = default)
        {
            try
            {
                return await this.FetchAsync(key.Address, DateTime.Today, cancellationToken);
            }
            catch (Exception e)
            {
                var msg = e switch
                {
                    HttpRequestException => "Unable to fetch from the Seattle Waste API",
                    JsonException => "Unable to deserialize response from the Seattle Waste API",
                    _ => "Unable to send to the Seattle Waste API"
                };
                this.Logger.LogError(msg + "; {exception}", e);
                return null;
            }
        }

        /// <summary>
        /// Fetch one response from the source.
        /// </summary>
        /// <remarks>
        /// c.2021 the Seattle collection page changed APIs.
        /// Similar in implementation to this work https://github.com/mampfes/hacs_waste_collection_schedule/pull/51
        /// </remarks>
        /// <param name="address"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        protected async Task<Response?> FetchAsync(string address, DateTime today,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug("Started finding {address} from Seattle Waste", address);

            var accountNumberTask = this.FetchAccountAsync(address, cancellationToken);
            var tokenTask = this.FetchTokenAsync(cancellationToken);

            await Task.WhenAll(accountNumberTask, tokenTask);

            var accountNumber = await accountNumberTask;
            if (string.IsNullOrEmpty(accountNumber))
            {
                this.Logger.LogError("Unable to determine accountNumber {address}", address);
                return null;
            }

            var token = await tokenTask;
            if (string.IsNullOrEmpty(token))
            {
                this.Logger.LogError("Unable to obtain a token");
                return null;
            }

            var resp = await FetchCollectionAsync(accountNumber, token, today, cancellationToken);
            if (resp == null)
            {
                this.Logger.LogError("Unable to parse the solid waste calendar for {address}", address);
                return null;
            }

            return resp with
            {
                Address = address,
            };
        }

        /// <summary>
        /// The logger used internally.
        /// </summary>
        private readonly ILogger<SourceDAO> Logger;

        /// <summary>
        /// The client used to access the source.
        /// </summary>
        private readonly HttpClient Client;

        /// <summary>
        /// The max number of times to hit the API
        /// </summary>
        private const int MAX_API_CALLS = 5;

        private async Task<string?> FetchAccountAsync(string address, CancellationToken cancellationToken = default)
        {
            var premCode = await PremiseCodeAsync(address, cancellationToken);
            if (string.IsNullOrEmpty(premCode))
            {
                this.Logger.LogError("Unable to determine premCode for {address}", address);
                return null;
            }

            var accountNumber = await AccountCodeAsync(premCode, cancellationToken);
            if (string.IsNullOrEmpty(accountNumber))
            {
                this.Logger.LogError("Unable to determine accountNumber for {address}", address);
                return null;
            }

            return accountNumber;

            async Task<string?> PremiseCodeAsync(string address, CancellationToken cancellationToken = default)
            {
                var url = "https://myutilities.seattle.gov/rest/serviceorder/findaddress";
                var payload = new StringContent(
                    JsonConvert.SerializeObject(new
                    { address = new { addressLine1 = address, } }),
                    Encoding.UTF8, "application/json"
                );
                var resp = await PostAsync<SVCFindAddResp>(url, payload, cancellationToken);
                return resp?.Addresses?.FirstOrDefault()?.PremCode;
            }

            async Task<string?> AccountCodeAsync(string premCode, CancellationToken cancellationToken = default)
            {
                var url = "https://myutilities.seattle.gov/rest/serviceorder/findAccount";
                var payload = new StringContent(
                    JsonConvert.SerializeObject(new
                    { address = new { premCode = premCode, } }),
                    Encoding.UTF8, "application/json"
                );
                var resp = await PostAsync<SVCFindAcctResp>(url, payload, cancellationToken);
                return resp?.Account.AccountNumber;
            }
        }

        private async Task<string?> FetchTokenAsync(CancellationToken cancellationToken = default)
        {
            var url = "https://myutilities.seattle.gov/rest/auth/token";
            var payload = new FormUrlEncodedContent(
                new List<KeyValuePair<string?, string?>>
                {
                    new KeyValuePair<string?, string?>("grant_type", "password"),
                    new KeyValuePair<string?, string?>("username", "guest"),
                    new KeyValuePair<string?, string?>("password", "guest"),
                }
            );
            var resp = await PostAsync<SVCTokenResp>(url, payload, cancellationToken);
            return resp?.AccessToken;
        }

        private async Task<Response?> FetchCollectionAsync(string accountNumber, string token, DateTime today, CancellationToken cancellationToken = default)
        {
            var services = await SummaryAsync(accountNumber, token, cancellationToken);
            var solidWaste = await CollectionAsync(accountNumber, token, services, cancellationToken);
            
            // Can't do anything without the calendar
            if (solidWaste?.Calendar == null) {
                this.Logger.LogError("Unable to locate solid waste calendar; bailing");
                return null;
            }

            var garbageSvc = services.FirstOrDefault(x => x.Description == "Garbage") ?? new();
            var recyclingSvc = services.FirstOrDefault(x => x.Description == "Recycle") ?? new();
            var faywSvc = services.FirstOrDefault(x => x.Description == "Food/Yard Waste") ?? new();
            var calendar = solidWaste.Calendar;
            calendar.TryGetValue(garbageSvc.ServicePointId, out var garbageCal);
            calendar.TryGetValue(recyclingSvc.ServicePointId, out var recyclingCal);
            calendar.TryGetValue(faywSvc.ServicePointId, out var faywCal);

            // Base everything on the garbage calendar
            if (garbageCal == null || !garbageCal.Any())
            {
                this.Logger.LogError("Unable to discover garbage service from solid waste response; bailing");
                return null;
            }

            // Find the next date; these appear to be ordered
            var garbageDate = garbageCal.FirstOrDefault(x => DateTime.Parse(x) >= today);
            if (string.IsNullOrEmpty(garbageDate))
            {
                this.Logger.LogError("Unable to discover latest garabge collection date from solid waste response; bailing");
                return null;
            }

            return new Response
            {
                Garbage = true,
                Recycling = recyclingCal?.Any(x => x == garbageDate) ?? false,
                FoodAndYardWaste = faywCal?.Any(x => x == garbageDate) ?? false,
                Start = DateTime.Parse(garbageDate),
            };

            async Task<IEnumerable<SVCSummaryRespSeviceSummaryService>> SummaryAsync(string accountNumber, string token, CancellationToken cancellationToken = default)
            {
                var request = new HttpRequestMessage(HttpMethod.Post, "https://myutilities.seattle.gov/rest/account/swsummary");
                var payload = new StringContent(
                    JsonConvert.SerializeObject(new
                    {
                        customerId = "guest",
                        accountContext = new
                        {
                            accountNumber = accountNumber,
                        }
                    }),
                    Encoding.UTF8, "application/json"
                );
                request.Headers.Add("Authorization", $"Bearer {token}");
                request.Content = payload;

                var resp = await this.Client.SendAsync(request, cancellationToken);
                resp.EnsureSuccessStatusCode();
                var content = await resp.Content.ReadAsStringAsync();
                var obj = JsonConvert.DeserializeObject<SVCSummaryResp>(content);

                return obj?.AccountSummaryType.ServiceSummaries.FirstOrDefault()?.Services ?? new();
            }

            async Task<SVCSolidWasteResp> CollectionAsync(string accountNumber, string token, IEnumerable<SVCSummaryRespSeviceSummaryService> services, CancellationToken cancellationToken = default)
            {
                var request = new HttpRequestMessage(HttpMethod.Post, "https://myutilities.seattle.gov/rest/solidwastecalendar");
                var payload = new StringContent(
                    JsonConvert.SerializeObject(new
                    {
                        customerId = "guest",
                        accountContext = new
                        {
                            accountNumber = accountNumber,
                            companyCd = "SPU",
                        },
                        servicePoints = services.Select(x => x.ServicePointId),
                    }),
                    Encoding.UTF8, "application/json"
                );
                request.Headers.Add("Authorization", $"Bearer {token}");
                request.Content = payload;

                var resp = await this.Client.SendAsync(request, cancellationToken);
                resp.EnsureSuccessStatusCode();
                var content = await resp.Content.ReadAsStringAsync();
                var obj = JsonConvert.DeserializeObject<SVCSolidWasteResp>(content);
                return obj;
            }
        }

        private async Task<T?> PostAsync<T>(string url, HttpContent payload, CancellationToken cancellationToken = default)
        {
            var resp = await this.Client.PostAsync(url, payload, cancellationToken);
            resp.EnsureSuccessStatusCode();
            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<T>(content);
            return obj;
        }
    }

    class SVCFindAddResp
    {
        [JsonProperty("address")]
        public List<SVCFindAddyAddress> Addresses { get; set; } = new List<SVCFindAddyAddress>();
    }

    class SVCFindAddyAddress 
    {
        [JsonProperty("premCode")]
        public string PremCode { get; set; } = string.Empty;
    }

    public class SVCFindAcctResp
    {
        [JsonProperty("account")]
        public SVCFindAcctAccount Account { get; set; } = new();
    }

    public class SVCFindAcctAccount 
    {
        [JsonProperty("accountNumber")]
        public string AccountNumber { get; set; } = string.Empty;
    }

    public class SVCTokenResp
    {
        [JsonProperty("access_token")]
        public string AccessToken { get; set; } = string.Empty;
    }

    public class SVCSummaryResp
    {
        [JsonProperty("accountSummaryType")]
        public SVCSummaryRespSummary AccountSummaryType { get; set; } = new();
    }

    public class SVCSummaryRespSummary
    {
        [JsonProperty("persondId")]
        public string Person { get; set; } = string.Empty;

        [JsonProperty("companyCd")]
        public string Company { get; set; } = string.Empty;

        [JsonProperty("swServices")]
        public List<SVCSummaryRespSeviceSummary> ServiceSummaries { get; set; } = new();
    }

    public class SVCSummaryRespSeviceSummary
    {
        [JsonProperty("services")]
        public List<SVCSummaryRespSeviceSummaryService> Services { get; set; } = new();
    }

    public class SVCSummaryRespSeviceSummaryService
    {
        [JsonProperty("description")]
        public string Description { get; set; } = string.Empty;

        [JsonProperty("servicePointId")]
        public string ServicePointId { get; set; } = string.Empty;
    }

    public class SVCSolidWasteResp
    {
        [JsonProperty("calendar")]
        public Dictionary<string, List<string>> Calendar { get; set; } = new();
    }
}
