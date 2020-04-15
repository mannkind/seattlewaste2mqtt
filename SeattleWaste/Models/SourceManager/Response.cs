using System;

namespace SeattleWaste.Models.SourceManager
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public class Response
    {
        public string Address { get; set; } = string.Empty;
        public DateTime Start { get; set; } = DateTime.MaxValue;
        public bool Garbage { get; set; } = false;
        public bool Recycling { get; set; } = false;
        public bool FoodAndYardWaste { get; set; } = false;
        public bool Ok {get; set;} = false;

        public override string ToString() => $"Date: {this.Start}, G: {this.Garbage}, R: {this.Recycling}, F: {this.FoodAndYardWaste}";
    }
}