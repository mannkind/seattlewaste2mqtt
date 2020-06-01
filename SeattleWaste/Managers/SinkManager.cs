using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core;
using TwoMQTT.Core.Managers;

namespace SeattleWaste
{
    /// <summary>
    /// An class representing a managed way to interact with a sink.
    /// </summary>
    public class SinkManager : MQTTManager<SlugMapping, Resource, Command>
    {
        /// <summary>
        /// Initializes a new instance of the SinkManager class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="sharedOpts"></param>
        /// <param name="opts"></param>
        /// <param name="incomingData"></param>
        /// <param name="outgoingCommand"></param>
        /// <returns></returns>
        public SinkManager(ILogger<SinkManager> logger, IOptions<Opts> sharedOpts,
            IOptions<Models.SinkManager.Opts> opts, ChannelReader<Resource> incomingData,
            ChannelWriter<Command> outgoingCommand) :
            base(logger, opts, incomingData, outgoingCommand, sharedOpts.Value.Resources, string.Empty)
        {
        }

        /// <inheritdoc />
        protected override async Task HandleIncomingDataAsync(Resource input,
             CancellationToken cancellationToken = default)
        {
            var slug = this.Questions
                .Where(x => x.Address == input.Address)
                .Select(x => x.Slug)
                .FirstOrDefault() ?? string.Empty;

            this.Logger.LogDebug($"Found slug {slug} for incoming data for {input.Address}");
            if (string.IsNullOrEmpty(slug))
            {
                this.Logger.LogDebug($"Unable to find slug for {input.Address}");
                return;
            }

            this.Logger.LogDebug($"Started publishing data for slug {slug}");
            await Task.WhenAll(
                this.PublishAsync(
                    this.StateTopic(slug, nameof(Resource.Start)), input.Start.ToShortDateString(),
                    cancellationToken
                ),
                this.PublishAsync(
                    this.StateTopic(slug, nameof(Resource.Garbage)), this.BooleanOnOff(input.Garbage),
                    cancellationToken
                ),
                this.PublishAsync(
                    this.StateTopic(slug, nameof(Resource.Recycling)), this.BooleanOnOff(input.Recycling),
                    cancellationToken
                ),
                this.PublishAsync(
                    this.StateTopic(slug, nameof(Resource.FoodAndYardWaste)), this.BooleanOnOff(input.FoodAndYardWaste),
                    cancellationToken
                ),
                this.PublishAsync(
                    this.StateTopic(slug, nameof(Resource.Status)), this.BooleanOnOff(input.Status),
                    cancellationToken
                )
            );
            this.Logger.LogDebug($"Finished publishing data for slug {slug}");
        }

        /// <inheritdoc />
        protected override async Task HandleDiscoveryAsync(CancellationToken cancellationToken = default)
        {
            if (!this.Opts.DiscoveryEnabled)
            {
                return;
            }

            var tasks = new List<Task>();
            var assembly = Assembly.GetAssembly(typeof(Program))?.GetName() ?? new AssemblyName();
            var mapping = new[]
            {
                new { Sensor = nameof(Resource.Start), Type = Const.SENSOR },
                new { Sensor = nameof(Resource.Garbage), Type = Const.BINARY_SENSOR },
                new { Sensor = nameof(Resource.Recycling), Type = Const.BINARY_SENSOR },
                new { Sensor = nameof(Resource.FoodAndYardWaste), Type = Const.BINARY_SENSOR },
                new { Sensor = nameof(Resource.Status), Type = Const.BINARY_SENSOR },
            };

            foreach (var input in this.Questions)
            {
                foreach (var map in mapping)
                {
                    var discovery = this.BuildDiscovery(input.Slug, map.Sensor, assembly, false);
                    tasks.Add(this.PublishDiscoveryAsync(input.Slug, map.Sensor, map.Type, discovery, cancellationToken));
                }
            }

            await Task.WhenAll(tasks);
        }
    }
}