using System.Collections.Generic;
using System.Threading.Channels;
using System.Threading.Tasks;
using System.Threading;
using System.Linq;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Newtonsoft.Json;
using System;
using System.Net.Http;
using SeattleWaste.Models.Shared;
using System.Net;
using TwoMQTT.Core.Managers;

namespace SeattleWaste
{
    public class Source : HTTPManager<SlugMapping, Resource, Command>
    {
        public Source(ILogger<Source> logger, IOptions<Opts> sharedOpts, IOptions<Models.SourceManager.Opts> opts, ChannelWriter<Resource> outgoing, ChannelReader<Command> incoming, IHttpClientFactory httpClientFactory) :
            base(logger, outgoing, incoming, httpClientFactory.CreateClient())
        {
            this.opts = opts.Value;
            this.sharedOpts = sharedOpts.Value;
        }
        protected readonly Models.SourceManager.Opts opts;
        protected readonly Opts sharedOpts;

        /// <inheritdoc />
        protected override void LogSettings()
        {
            var resources = string.Join(",",
                this.sharedOpts.Resources.Select(x => $"{x.Address}:{x.Slug}")
            );

            this.logger.LogInformation(
                $"AlertWithin:           {this.opts.AlertWithin}\n" +
                $"PollingInterval:       {this.opts.PollingInterval}\n" +
                $"Resources:             {resources}\n" +
                $""
            );
        }

        /// <inheritdoc />
        protected override async Task PollAsync(CancellationToken cancellationToken = default)
        {
            this.logger.LogInformation("Polling");

            var tasks = new List<Task<Models.SourceManager.Response>>();
            foreach (var key in this.sharedOpts.Resources)
            {
                this.logger.LogInformation($"Looking up {key}");
                tasks.Add(this.FetchOneAsync(key, cancellationToken));
            }

            var results = await Task.WhenAll(tasks);
            foreach (var result in results.Where(x => x.Ok))
            {
                this.logger.LogInformation($"Found {result}");
                await this.outgoing.WriteAsync(Resource.From(result, this.opts.AlertWithin), cancellationToken);
            }
        }

        /// <inheritdoc />
        protected override Task DelayAsync(CancellationToken cancellationToken = default) =>
            Task.Delay(this.opts.PollingInterval, cancellationToken);

        /// <summary>
        /// Fetch one record from the source
        /// </summary>
        private async Task<Models.SourceManager.Response> FetchOneAsync(SlugMapping key, CancellationToken cancellationToken = default)
        {
            var apiCalls = 0;
            var lastTimeStamp = 0L;
            var todayTimeStamp = ((DateTimeOffset)DateTime.Today).ToUnixTimeSeconds();

            while (lastTimeStamp < todayTimeStamp && apiCalls <= MAX_API_CALLS)
            {
                var collections = await this.FetchAllAsync(key, lastTimeStamp, cancellationToken);
                foreach (var collection in collections)
                {
                    lastTimeStamp = ((DateTimeOffset)collection.Start).ToUnixTimeSeconds();
                    if (todayTimeStamp < lastTimeStamp)
                    {
                        continue;
                    }
                    
                    collection.Address = key.Address;
                    collection.Ok = true;
                    return collection;
                }

                apiCalls += 1;
            }

            return new Models.SourceManager.Response();
        }


        /// <summary>
        /// Fetch one record from the source
        /// </summary>
        private async Task<IEnumerable<Models.SourceManager.Response>> FetchAllAsync(SlugMapping key, long start, CancellationToken cancellationToken = default)
        {
            var baseUrl = "https://www.seattle.gov/UTIL/WARP/CollectionCalendar/GetCollectionDays";
            var query = $"pApp=CC&pAddress={WebUtility.UrlEncode(key.Address)}&start={start}";
            var resp = await this.client.GetAsync($"{baseUrl}?{query}", cancellationToken);
            if (!resp.IsSuccessStatusCode)
            {
                return new List<Models.SourceManager.Response>();
            }

            var content = await resp.Content.ReadAsStringAsync();
            var obj = JsonConvert.DeserializeObject<List<Models.SourceManager.Response>>(content);

            return obj;
        }

        private const int MAX_API_CALLS = 5;
    }
}
