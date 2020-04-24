using System;

namespace SeattleWaste.Models.SourceManager
{
    /// <summary>
    /// The source options
    /// </summary>
    public class Opts
    {
        public const string Section = "SeattleWaste";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan AlertWithin { get; set; } = new TimeSpan(24, 0, 0);

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(8, 3, 31);
    }
}
