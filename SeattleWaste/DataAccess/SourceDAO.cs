using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using SeattleWaste.Models.Shared;
using SeattleWaste.Models.Source;
using TwoMQTT.Core.Interfaces;

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
                var todayTimeStamp = ((DateTimeOffset)DateTime.Today).ToUnixTimeSeconds();
                return await this.FetchAsync(key.Address, todayTimeStamp, 0L, cancellationToken);
            }
            catch (Exception e)
            {
                var msg = e switch
                {
                    HttpRequestException => "Unable to fetch from the Seattle Waste API",
                    JsonException => "Unable to deserialize response from the Seattle Waste API",
                    _ => "Unable to send to the Seattle Waste API"
                };
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
        protected async Task<Response?> FetchAsync(string address,
            long todayTimeStamp,
            long lastTimeStamp,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug("Started finding {address} from Seattle Waste", address);
            var apiCalls = 0;

            // Limit the number of times we'll hit the source before giving up
            while (lastTimeStamp <= todayTimeStamp && apiCalls <= MAX_API_CALLS)
            {
                this.Logger.LogDebug("{apiCalls} iteration;  timestamp {lastTimeStamp}", apiCalls, lastTimeStamp);
                var collections = await this.FetchAllAsync(address, lastTimeStamp, cancellationToken);
                foreach (var collection in collections)
                {
                    lastTimeStamp = ((DateTimeOffset)collection.Start).ToUnixTimeSeconds();
                    if (lastTimeStamp <= todayTimeStamp)
                    {
                        continue;
                    }

                    var addressedCollection = collection with
                    {
                        Address = address,
                    };

                    return addressedCollection;
                }

                apiCalls += 1;
            }

            return null;
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

        /// <summary>
        /// Fetch all records from the source.
        /// </summary>
        private async Task<IEnumerable<Response>> FetchAllAsync(string address, long start,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug("Started finding collection days for {address} @ {start} from Seattle Waste", address, start);
            var baseUrl = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays";
            var query = $"pApp=CC&pAddress={WebUtility.UrlEncode(address)}&start={start}";
            var resp = await this.Client.GetAsync($"{baseUrl}?{query}", cancellationToken);
            resp.EnsureSuccessStatusCode();
            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<List<Response>>(content);
            this.Logger.LogDebug("Finished finding collection days for {address} @ {start} from Seattle Waste", address, start);

            return obj;
        }
    }
}
