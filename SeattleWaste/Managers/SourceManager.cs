using System;
using System.Linq;
using System.Threading.Channels;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using SeattleWaste.Models.Shared;
using SeattleWaste.Models.SourceManager;
using TwoMQTT.Core.DataAccess;
using TwoMQTT.Core.Managers;

namespace SeattleWaste.Managers
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceManager : APIPollingManager<SlugMapping, FetchResponse, object, Resource, Command>
    {
        public SourceManager(ILogger<SourceManager> logger, IOptions<Models.Shared.Opts> sharedOpts,
            IOptions<Models.SourceManager.Opts> opts, ChannelWriter<Resource> outgoing, ChannelReader<Command> incoming,
            ISourceDAO<SlugMapping, Command, FetchResponse, object> sourceDAO) :
            base(logger, outgoing, incoming, sharedOpts.Value.Resources, opts.Value.PollingInterval, sourceDAO,
                SourceSettings(sharedOpts.Value, opts.Value))
        {
        }

        /// <inheritdoc />
        protected override Resource MapResponse(FetchResponse src) =>
            new Resource
            {
                Address = src.Address,
                Start = src.Start,
                Garbage = src.Garbage,
                Recycling = src.Recycling,
                FoodAndYardWaste = src.FoodAndYardWaste,
            };

        private static string SourceSettings(Models.Shared.Opts sharedOpts, Models.SourceManager.Opts opts) =>
            $"PollingInterval: {opts.PollingInterval}\n" +
            $"Resources: {string.Join(',', sharedOpts.Resources.Select(x => $"{x.Address}:{x.Slug}"))}\n" +
            $"";
    }
}
