using System;

namespace SeattleWaste.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public record SourceOpts
    {
        public const string Section = "SeattleWaste";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; init; } = new(8, 3, 31);
    }
}
