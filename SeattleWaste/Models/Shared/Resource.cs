using System;

namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public record Resource
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
        public DateTime Start { get; init; } = DateTime.MaxValue;

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
