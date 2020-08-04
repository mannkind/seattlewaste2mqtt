using System;

namespace SeattleWaste.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public class SourceOpts
    {
        public const string Section = "SeattleWaste";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(8, 3, 31);
    }
}
