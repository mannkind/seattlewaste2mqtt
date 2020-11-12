using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using SeattleWaste.Models.Options;
using SeattleWaste.Models.Shared;
using TwoMQTT;
using TwoMQTT.Interfaces;
using TwoMQTT.Liasons;
using TwoMQTT.Models;
using TwoMQTT.Utils;

namespace SeattleWaste.Liasons
{
    /// <summary>
    /// An class representing a managed way to interact with MQTT.
    /// </summary>
    public class MQTTLiason : MQTTLiasonBase<Resource, object, SlugMapping, SharedOpts>, IMQTTLiason<Resource, object>
    {
        /// <summary>
        /// 
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="generator"></param>
        /// <param name="sharedOpts"></param>
        public MQTTLiason(ILogger<MQTTLiason> logger, IMQTTGenerator generator, IOptions<SharedOpts> sharedOpts) :
            base(logger, generator, sharedOpts)
        {
        }

        /// <inheritdoc />
        public IEnumerable<(string topic, string payload)> MapData(Resource input)
        {
            var results = new List<(string, string)>();
            var slug = this.Questions
                .Where(x => x.Address == input.Address)
                .Select(x => x.Slug)
                .FirstOrDefault() ?? string.Empty;

            if (string.IsNullOrEmpty(slug))
            {
                this.Logger.LogDebug("Unable to find slug for {address}", input.Address);
                return results;
            }

            this.Logger.LogDebug("Found slug {slug} for incoming data for {address}", slug, input.Address);
            results.AddRange(new[]
                {
                    (this.Generator.StateTopic(slug, nameof(Resource.Start)), input.Start.ToShortDateString()),
                    (this.Generator.StateTopic(slug, nameof(Resource.Garbage)), this.Generator.BooleanOnOff(input.Garbage)),
                    (this.Generator.StateTopic(slug, nameof(Resource.Recycling)), this.Generator.BooleanOnOff(input.Recycling)),
                    (this.Generator.StateTopic(slug, nameof(Resource.FoodAndYardWaste)), this.Generator.BooleanOnOff(input.FoodAndYardWaste)),
                }
            );

            return results;
        }

        /// <inheritdoc />
        public IEnumerable<(string slug, string sensor, string type, MQTTDiscovery discovery)> Discoveries()
        {
            var discoveries = new List<(string, string, string, MQTTDiscovery)>();
            var assembly = Assembly.GetAssembly(typeof(Program))?.GetName() ?? new AssemblyName();
            var mapping = new[]
            {
                new { Sensor = nameof(Resource.Start), Type = Const.SENSOR },
                new { Sensor = nameof(Resource.Garbage), Type = Const.BINARY_SENSOR },
                new { Sensor = nameof(Resource.Recycling), Type = Const.BINARY_SENSOR },
                new { Sensor = nameof(Resource.FoodAndYardWaste), Type = Const.BINARY_SENSOR },
            };

            foreach (var input in this.Questions)
            {
                foreach (var map in mapping)
                {
                    var discovery = this.Generator.BuildDiscovery(input.Slug, map.Sensor, assembly, false);
                    discoveries.Add((input.Slug, map.Sensor, map.Type, discovery));
                }
            }

            return discoveries;
        }
    }
}