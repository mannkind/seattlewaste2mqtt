using TwoMQTT.Models;

namespace SeattleWaste.Models.Options;

/// <summary>
/// The sink options
/// </summary>
public record MQTTOpts : MQTTManagerOptions
{
    public const string Section = "SeattleWaste:MQTT";
    public const string TopicPrefixDefault = "home/seattle_waste";
    public const string DiscoveryNameDefault = "seattle_waste";
}
