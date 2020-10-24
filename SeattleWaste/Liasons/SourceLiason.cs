using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using SeattleWaste.DataAccess;
using SeattleWaste.Models.Options;
using SeattleWaste.Models.Shared;
using SeattleWaste.Models.Source;
using TwoMQTT.Interfaces;
using TwoMQTT.Liasons;

namespace SeattleWaste.Liasons
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceLiason : SourceLiasonBase<Resource, object, SlugMapping, ISourceDAO, SharedOpts>, ISourceLiason<Resource, object>
    {
        public SourceLiason(ILogger<SourceLiason> logger, ISourceDAO sourceDAO,
            IOptions<SourceOpts> opts, IOptions<SharedOpts> sharedOpts) :
            base(logger, sourceDAO, sharedOpts)
        {
            this.Logger.LogInformation(
                "PollingInterval: {pollingInterval}\n" +
                "Resources: {@resources}\n" +
                "",
                opts.Value.PollingInterval,
                sharedOpts.Value.Resources
            );
        }

        /// <inheritdoc />
        protected override async Task<Resource?> FetchOneAsync(SlugMapping key, CancellationToken cancellationToken)
        {
            var result = await this.SourceDAO.FetchOneAsync(key, cancellationToken);
            return result switch
            {
                Response => new Resource
                {
                    Address = result.Address,
                    Start = result.Start,
                    Garbage = result.Garbage,
                    Recycling = result.Recycling,
                    FoodAndYardWaste = result.FoodAndYardWaste,
                },
                _ => null,
            };
        }
    }
}
