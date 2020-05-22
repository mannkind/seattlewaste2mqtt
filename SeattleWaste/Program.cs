using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using SeattleWaste.DataAccess;
using SeattleWaste.Managers;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core;
using TwoMQTT.Core.DataAccess;
using TwoMQTT.Core.Extensions;


namespace SeattleWaste
{
    class Program : ConsoleProgram<Resource, Command, SourceManager, SinkManager>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            await p.ExecuteAsync(args);
        }

        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>>();

            return services
                .ConfigureOpts<Models.Shared.Opts>(hostContext, Models.Shared.Opts.Section)
                .ConfigureOpts<Models.SourceManager.Opts>(hostContext, Models.SourceManager.Opts.Section)
                .ConfigureOpts<Models.SinkManager.Opts>(hostContext, Models.SinkManager.Opts.Section)
                .AddTransient<ISourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>, SourceDAO>();
        }

        [Obsolete("Remove in the near future.")]
        private void MapOldEnvVariables()
        {
            var found = false;
            var foundOld = new List<string>();
            var mappings = new[]
            {
                new { Src = "SEATTLEWASTE_ADDRESS", Dst = "SEATTLEWASTE__RESOURCES", CanMap = true, Strip = "",  Sep = ":" },
                new { Src = "SEATTLEWASTE_ALERTWITHIN", Dst = "SEATTLEWASTE__ALERTWITHIN", CanMap = false, Strip = "", Sep = "" },
                new { Src = "SEATTLEWASTE_LOOKUPINTERVAL", Dst = "SEATTLEWASTE__POLLINGINTERVAL", CanMap = false, Strip = "", Sep = "" },
                new { Src = "MQTT_TOPICPREFIX", Dst = "SEATTLEWASTE__MQTT__TOPICPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERY", Dst = "SEATTLEWASTE__MQTT__DISCOVERYENABLED", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYPREFIX", Dst = "SEATTLEWASTE__MQTT__DISCOVERYPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYNAME", Dst = "SEATTLEWASTE__MQTT__DISCOVERYNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_BROKER", Dst = "SEATTLEWASTE__MQTT__BROKER", CanMap = true, Strip = "tcp://", Sep = "" },
                new { Src = "MQTT_USERNAME", Dst = "SEATTLEWASTE__MQTT__USERNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_PASSWORD", Dst = "SEATTLEWASTE__MQTT__PASSWORD", CanMap = true, Strip = "", Sep = "" },
            };

            foreach (var mapping in mappings)
            {
                var old = Environment.GetEnvironmentVariable(mapping.Src);
                if (string.IsNullOrEmpty(old))
                {
                    continue;
                }

                found = true;
                foundOld.Add($"{mapping.Src} => {mapping.Dst}");

                if (!mapping.CanMap)
                {
                    continue;
                }

                // Strip junk where possible
                if (!string.IsNullOrEmpty(mapping.Strip))
                {
                    old = old.Replace(mapping.Strip, string.Empty);
                }

                // Simple
                if (string.IsNullOrEmpty(mapping.Sep))
                {
                    Environment.SetEnvironmentVariable(mapping.Dst, old);
                }
                // Complex
                else
                {
                    var resourceSlugs = old.Split(",");
                    var i = 0;
                    foreach (var resourceSlug in resourceSlugs)
                    {
                        var parts = resourceSlug.Split(mapping.Sep);
                        var id = parts.Length >= 1 ? parts[0] : string.Empty;
                        var slug = parts.Length >= 2 ? parts[1] : string.Empty;
                        var idEnv = $"{mapping.Dst}__{i}__Address";
                        var slugEnv = $"{mapping.Dst}__{i}__Slug";
                        Environment.SetEnvironmentVariable(idEnv, id);
                        Environment.SetEnvironmentVariable(slugEnv, slug);
                    }
                }

            }


            if (found)
            {
                var loggerFactory = LoggerFactory.Create(builder => { builder.AddConsole(); });
                var logger = loggerFactory.CreateLogger<Program>();
                logger.LogWarning("Found old environment variables.");
                logger.LogWarning($"Please migrate to the new environment variables: {(string.Join(", ", foundOld))}");
                Thread.Sleep(5000);
            }
        }
    }
}
