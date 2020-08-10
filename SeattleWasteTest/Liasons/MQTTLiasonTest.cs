using System.Linq;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using TwoMQTT.Core.Utils;
using SeattleWaste.Liasons;
using SeattleWaste.Models.Options;
using SeattleWaste.Models.Shared;
using System;

namespace SeattleWasteTest.Liasons
{
    [TestClass]
    public class MQTTLiasonTest
    {
        [TestMethod]
        public void MapDataTest()
        {
            var tests = new[] {
                new {
                    Q = new SlugMapping { Address = BasicAddress, Slug = BasicSlug },
                    Resource = new Resource { Address = BasicAddress, Start = BasicStart },
                    Expected = new { Address = BasicAddress, Start = BasicStart.ToShortDateString(), Slug = BasicSlug, Found = true }
                },
                new {
                    Q = new SlugMapping { Address = BasicAddress, Slug = BasicSlug },
                    Resource = new Resource { Address = $"{BasicAddress}-fake" , Start = BasicStart },
                    Expected = new { Address = string.Empty, Start = DateTime.MinValue.ToShortDateString(), Slug = string.Empty, Found = false }
                },
            };

            foreach (var test in tests)
            {
                var logger = new Mock<ILogger<MQTTLiason>>();
                var generator = new Mock<IMQTTGenerator>();
                var sharedOpts = Options.Create(new SharedOpts
                {
                    Resources = new[] { test.Q }.ToList(),
                });

                generator.Setup(x => x.BuildDiscovery(It.IsAny<string>(), It.IsAny<string>(), It.IsAny<System.Reflection.AssemblyName>(), false))
                    .Returns(new TwoMQTT.Core.Models.MQTTDiscovery());
                generator.Setup(x => x.StateTopic(test.Q.Slug, nameof(Resource.Start)))
                    .Returns($"totes/{test.Q.Slug}/topic/{nameof(Resource.Start)}");

                var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
                var results = mqttLiason.MapData(test.Resource);
                var actual = results.FirstOrDefault();

                Assert.AreEqual(test.Expected.Found, results.Any(), "The mapping should exist if found.");
                if (test.Expected.Found)
                {
                    Assert.IsTrue(actual.topic.Contains(test.Expected.Slug), "The topic should contain the expected Address.");
                    Assert.AreEqual(test.Expected.Start, actual.payload, "The payload be the expected Start.");
                }
            }
        }

        [TestMethod]
        public void DiscoveriesTest()
        {
            var tests = new[] {
                new {
                    Q = new SlugMapping { Address = BasicAddress, Slug = BasicSlug },
                    Resource = new Resource { Address = BasicAddress, Start = BasicStart },
                    Expected = new { Address = BasicAddress, Start = BasicStart, Slug = BasicSlug }
                },
            };

            foreach (var test in tests)
            {
                var logger = new Mock<ILogger<MQTTLiason>>();
                var generator = new Mock<IMQTTGenerator>();
                var sharedOpts = Options.Create(new SharedOpts
                {
                    Resources = new[] { test.Q }.ToList(),
                });

                generator.Setup(x => x.BuildDiscovery(test.Q.Slug, nameof(Resource.Start), It.IsAny<System.Reflection.AssemblyName>(), false))
                    .Returns(new TwoMQTT.Core.Models.MQTTDiscovery());

                var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
                var results = mqttLiason.Discoveries();
                var result = results.FirstOrDefault();

                Assert.IsNotNull(result, "A discovery should exist.");
            }
        }

        private static string BasicSlug = "totallyaslug";
        private static DateTime BasicStart = DateTime.Now;
        private static string BasicAddress = "55881 SE 42nd St.";
    }
}
