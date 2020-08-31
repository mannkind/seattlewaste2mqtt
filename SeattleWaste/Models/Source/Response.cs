using System;

namespace SeattleWaste.Models.Source
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public class Response
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
        public DateTime Start { get; set; } = DateTime.MinValue;

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