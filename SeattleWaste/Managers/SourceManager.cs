using System.Threading.Channels;
using System.Linq;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using System;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core.Managers;
using SeattleWaste.Models.SourceManager;
using TwoMQTT.Core.DataAccess;

namespace SeattleWaste.Managers
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceManager : HTTPPollingManager<SlugMapping, FetchResponse, object, Resource, Command>
    {
        public SourceManager(ILogger<SourceManager> logger, IOptions<Models.Shared.Opts> sharedOpts, IOptions<Models.SourceManager.Opts> opts, ChannelWriter<Resource> outgoing, ChannelReader<Command> incoming, IHTTPSourceDAO<SlugMapping, Command, FetchResponse, object> sourceDAO) :
            base(logger, outgoing, incoming, sharedOpts.Value.Resources, opts.Value.PollingInterval, sourceDAO)
        {
            this.Opts = opts.Value;
            this.SharedOpts = sharedOpts.Value;
        }

        /// <inheritdoc />
        protected override void LogSettings() =>
            this.Logger.LogInformation(
                $"PollingInterval: {this.Opts.PollingInterval}\n" +
                $"Resources: {string.Join(',', this.SharedOpts.Resources.Select(x => $"{x.Address}:{x.Slug}"))}\n" +
                $""
            );

        /// <inheritdoc />
        protected override Resource MapResponse(FetchResponse src) =>
            new Resource
            {
                Address = src.Address,
                Start = src.Start,
                Garbage = src.Garbage,
                Recycling = src.Recycling,
                FoodAndYardWaste = src.FoodAndYardWaste,
                Status = src.Start.Subtract(DateTime.Now) <= this.Opts.AlertWithin,
            };

        /// <summary>
        /// The options for the source.
        /// </summary>
        private readonly Models.SourceManager.Opts Opts;

        /// <summary>
        /// The options that are shared.
        /// </summary>
        private readonly Models.Shared.Opts SharedOpts;
    }
}
