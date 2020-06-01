namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared key info => slug mapping across the application
    /// </summary>
    public class SlugMapping
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Address { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Slug { get; set; } = string.Empty;

        /// <inheritdoc />
        public override string ToString() => $"Address: {this.Address}, Slug: {this.Slug}";
    }
}
