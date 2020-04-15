namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared key info => slug mapping across the application
    /// </summary>
    public class SlugMapping
    {
        public string Address { get; set; } = string.Empty;
        public string Slug { get; set; } = string.Empty;

        public override string ToString() => $"Address: {this.Address}, Slug: {this.Slug}";
    }
}
