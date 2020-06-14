using System;

namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public class Resource
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
        public DateTime Start { get; set; } = DateTime.MaxValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Garbage { get; set; } = false;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Recycling { get; set; } = false;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool FoodAndYardWaste { get; set; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Date: {this.Start}, G: {this.Garbage}, R: {this.Recycling}, F: {this.FoodAndYardWaste}";
    }
}
