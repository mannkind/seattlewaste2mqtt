using System;

namespace SeattleWaste.Models.Source
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public record Response
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Address { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public DateTime Start { get; init; } = DateTime.MinValue;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Garbage { get; init; } = false;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool Recycling { get; init; } = false;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool FoodAndYardWaste { get; init; } = false;
    }
}