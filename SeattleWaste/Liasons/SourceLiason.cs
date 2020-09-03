using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using SeattleWaste.DataAccess;
using SeattleWaste.Models.Options;
using SeattleWaste.Models.Shared;
using SeattleWaste.Models.Source;
using TwoMQTT.Core.Interfaces;
using TwoMQTT.Core.Liasons;

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
            var resp = result != null ? this.MapData(result) : null;
            return resp;
        }

        /// <summary>
        /// Map the source response to a shared response representation.
        /// </summary>
        /// <param name="src"></param>
        /// <returns></returns>
        private Resource MapData(Response src) =>
            new Resource
            {
                Address = src.Address,
                Start = src.Start,
                Garbage = src.Garbage,
                Recycling = src.Recycling,
                FoodAndYardWaste = src.FoodAndYardWaste,
            };
    }
}
