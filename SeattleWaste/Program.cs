using System.Threading.Tasks;
using System.Threading.Channels;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core;

namespace SeattleWaste
{
    class Program
    {
        static async Task Main(string[] args)
        {
            if (AppVersion.PrintVersion<Program>(args))
            {
                return;
            }

            var builder = new HostBuilder()
                .ConfigureAppConfiguration((hostingContext, config) =>
                {
                    config.AddJsonFile("appsettings.json", optional: true);
                    config.AddEnvironmentVariables();
                    config.AddCommandLine(args);
                })
                .ConfigureServices((hostContext, services) =>
                {
                    services.AddOptions();
                    services.Configure<Models.Shared.Opts>(hostContext.Configuration.GetSection(Models.Shared.Opts.Section));
                    services.Configure<Models.SourceManager.Opts>(hostContext.Configuration.GetSection(Models.SourceManager.Opts.Section));
                    services.Configure<Models.SinkManager.Opts>(hostContext.Configuration.GetSection(Models.SinkManager.Opts.Section));

                    var dataComms = Channel.CreateUnbounded<Resource>();
                    var commandComms = Channel.CreateUnbounded<Command>();
                    services.AddSingleton<ChannelReader<Resource>>(x => dataComms.Reader);
                    services.AddSingleton<ChannelWriter<Resource>>(x => dataComms.Writer);
                    services.AddSingleton<ChannelReader<Command>>(x => commandComms.Reader);
                    services.AddSingleton<ChannelWriter<Command>>(x => commandComms.Writer);
                    services.AddHttpClient<Source>();
                    services.AddHostedService<Source>();
                    services.AddHostedService<Sink>();
                })
                .ConfigureLogging((hostingContext, logging) => {
                    logging.AddConsole();
                });

            await builder.RunConsoleAsync();
        }
    }
}