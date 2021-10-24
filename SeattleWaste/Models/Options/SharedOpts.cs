using System.Collections.Generic;
using SeattleWaste.Models.Shared;
using TwoMQTT.Interfaces;

namespace SeattleWaste.Models.Options;

/// <summary>
/// The shared options across the application
/// </summary>
public record SharedOpts : ISharedOpts<SlugMapping>
{
    public const string Section = "SeattleWaste";

    /// <summary>
    /// 
    /// </summary>
    /// <typeparam name="SlugMapping"></typeparam>
    /// <returns></returns>
    public List<SlugMapping> Resources { get; init; } = new();
}
