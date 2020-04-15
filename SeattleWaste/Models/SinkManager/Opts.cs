using TwoMQTT.Core.Models;

namespace SeattleWaste.Models.SinkManager
{
    /// <summary>
    /// The sink options
    /// </summary>
    public class Opts : MQTTManagerOptions
    {
        public const string Section = "SeattleWaste:Sink";

        public Opts()
        {
            this.TopicPrefix = "home/seattle_waste";
            this.DiscoveryName = "seattle_waste";
        }
    }
}