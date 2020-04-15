using System;

namespace SeattleWaste.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public class Resource
    {
        public string Address { get; set; } = string.Empty;
        public DateTime Start { get; set; } = DateTime.MaxValue;
        public bool Garbage { get; set; } = false;
        public bool Recycling { get; set; } = false;
        public bool FoodAndYardWaste { get; set; } = false;
        public bool Status { get; set; } = false;

        public override string ToString() => $"Date: {this.Start}, G: {this.Garbage}, R: {this.Recycling}, F: {this.FoodAndYardWaste}";

        public static Resource From(SourceManager.Response obj, TimeSpan alertWithin) => 
            new Resource
            {
                Address = obj.Address,
                Start = obj.Start,
                Garbage = obj.Garbage,
                Recycling = obj.Recycling,
                FoodAndYardWaste = obj.FoodAndYardWaste,
                Status = obj.Start.Subtract(DateTime.Now) <= alertWithin,
            };
    }
}
