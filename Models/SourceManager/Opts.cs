using System;

namespace SeattleWaste.Models.SourceManager
{
    /// <summary>
    /// The source options
    /// </summary>
    public class Opts
    {
        public const string Section = "SeattleWaste:Source";

        public TimeSpan AlertWithin { get; set; } = new TimeSpan(24, 0, 0);
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(8, 3, 31);

        public override string ToString() => $"Alert Within: {this.AlertWithin}, Polling Interval: {this.PollingInterval}";
    }
}
