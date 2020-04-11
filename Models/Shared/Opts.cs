using System.Collections.Generic;

namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public class Opts
    {
        public const string Section = "SeattleWaste:Shared";

        public List<SlugMapping> Resources { get; set; } = new List<SlugMapping>();
    }
}
