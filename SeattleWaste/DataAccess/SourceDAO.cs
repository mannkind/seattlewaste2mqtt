using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Newtonsoft.Json;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core.DataAccess;

namespace SeattleWaste.DataAccess
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceDAO : HTTPSourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>
    {
        /// <summary>
        /// Initializes a new instance of the SourceDAO class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="opts"></param>
        /// <param name="httpClientFactory"></param>
        /// <returns></returns>
        public SourceDAO(ILogger<SourceDAO> logger, IOptions<Models.SourceManager.Opts> opts,
            IHttpClientFactory httpClientFactory) :
            base(logger, httpClientFactory)
        {
        }

        /// <inheritdoc />
        public override async Task<Models.SourceManager.FetchResponse?> FetchOneAsync(SlugMapping key,
            CancellationToken cancellationToken = default)
        {
            try
            {
                return await this.FetchAsync(key.Address, cancellationToken);
            }
            catch (Exception e)
            {
                var msg = e is HttpRequestException ? "Unable to fetch from the Seattle Waste API" :
                          e is JsonException ? "Unable to deserialize response from the Seattle Waste API" :
                          "Unable to send to the Seattle Waste API";
                this.Logger.LogError(msg, e);
                return null;
            }
        }

        /// <summary>
        /// Fetch one response from the source.
        /// </summary>
        /// <param name="address"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        private async Task<Models.SourceManager.FetchResponse?> FetchAsync(string address,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug($"Started finding {address} from Seattle Waste");
            var apiCalls = 0;
            var lastTimeStamp = 0L;
            var todayTimeStamp = ((DateTimeOffset)DateTime.Today).ToUnixTimeSeconds();

            // Limit the number of times we'll hit the source before giving up
            while (lastTimeStamp < todayTimeStamp && apiCalls <= MAX_API_CALLS)
            {
                this.Logger.LogDebug($"{apiCalls} iteration;  timestamp {lastTimeStamp}");
                var collections = await this.FetchAllAsync(address, lastTimeStamp, cancellationToken);
                foreach (var collection in collections)
                {
                    lastTimeStamp = ((DateTimeOffset)collection.Start).ToUnixTimeSeconds();
                    if (lastTimeStamp <= todayTimeStamp)
                    {
                        continue;
                    }

                    collection.Address = address;
                    return collection;
                }

                apiCalls += 1;
            }

            return null;
        }

        /// <summary>
        /// Fetch all records from the source.
        /// </summary>
        private async Task<IEnumerable<Models.SourceManager.FetchResponse>> FetchAllAsync(string address, long start,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug($"Started finding collection days for {address} @ {start} from Seattle Waste");
            var baseUrl = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays";
            var query = $"pApp=CC&pAddress={WebUtility.UrlEncode(address)}&start={start}";
            var resp = await this.Client.GetAsync($"{baseUrl}?{query}", cancellationToken);
            resp.EnsureSuccessStatusCode();
            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<List<Models.SourceManager.FetchResponse>>(content);
            this.Logger.LogDebug($"Finished finding collection days for {address} @ {start} from Seattle Waste");

            return obj;
        }

        private const int MAX_API_CALLS = 5;
    }
}
