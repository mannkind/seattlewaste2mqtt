﻿using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Options;
using SeattleWaste.DataAccess;
using SeattleWaste.Liasons;
using SeattleWaste.Models.Shared;
using TwoMQTT;
using TwoMQTT.Extensions;
using TwoMQTT.Interfaces;
using TwoMQTT.Managers;


namespace SeattleWaste
{
    class Program : ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            await p.ExecuteAsync(args);
        }

        /// <inheritdoc />
        protected override IDictionary<string, string> EnvironmentDefaults()
        {
            var sep = "__";
            var section = Models.Options.MQTTOpts.Section.Replace(":", sep);
            var sectsep = $"{section}{sep}";

            return new Dictionary<string, string>
            {
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.TopicPrefix)}", Models.Options.MQTTOpts.TopicPrefixDefault },
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.DiscoveryName)}", Models.Options.MQTTOpts.DiscoveryNameDefault },
            };
        }

        /// <inheritdoc />
        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO>();

            return services
                .ConfigureOpts<Models.Options.SharedOpts>(hostContext, Models.Options.SharedOpts.Section)
                .ConfigureOpts<Models.Options.SourceOpts>(hostContext, Models.Options.SourceOpts.Section)
                .ConfigureOpts<TwoMQTT.Models.MQTTManagerOptions>(hostContext, Models.Options.MQTTOpts.Section)
                .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.PollingInterval)} is required for {nameof(ThrottleManager)}.");
                    }

                    return new ThrottleManager(opts.Value.PollingInterval);
                })
                .AddSingleton<ISourceDAO, SourceDAO>();
        }
    }
}
