using System.Collections.Generic;
using SeattleWaste.Models.Shared;
using TwoMQTT.Core.Interfaces;

namespace SeattleWaste.Models.Options
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public class SharedOpts : ISharedOpts<SlugMapping>
    {
        public const string Section = "SeattleWaste";

        /// <summary>
        /// 
        /// </summary>
        /// <typeparam name="SlugMapping"></typeparam>
        /// <returns></returns>
        public List<SlugMapping> Resources { get; set; } = new List<SlugMapping>();
    }
}
