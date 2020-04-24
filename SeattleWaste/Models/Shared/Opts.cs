using System.Collections.Generic;
using System.Linq;

namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public class Opts
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
