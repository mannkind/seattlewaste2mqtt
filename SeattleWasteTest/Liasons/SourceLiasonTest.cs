using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using SeattleWaste.DataAccess;
using SeattleWaste.Liasons;
using SeattleWaste.Models.Options;
using SeattleWaste.Models.Shared;

namespace SeattleWasteTest.Liasons
{
    [TestClass]
    public class SourceLiasonTest
    {
        [TestMethod]
        public async Task FetchAllAsyncTest()
        {
            var tests = new[] {
                new {
                    Q = new SlugMapping { Address = BasicAddress, Slug = BasicSlug },
                    Expected = new { Address = BasicAddress, Start = BasicStart }
                },
            };

            foreach (var test in tests)
            {
                var logger = new Mock<ILogger<SourceLiason>>();
                var sourceDAO = new Mock<ISourceDAO>();
                var opts = Options.Create(new SourceOpts());
                var sharedOpts = Options.Create(new SharedOpts
                {
                    Resources = new[] { test.Q }.ToList(),
                });

                sourceDAO.Setup(x => x.FetchOneAsync(test.Q, It.IsAny<CancellationToken>()))
                     .ReturnsAsync(new SeattleWaste.Models.Source.FetchResponse
                     {
                         Address = test.Expected.Address,
                         Start = test.Expected.Start,
                     });

                var sourceLiason = new SourceLiason(logger.Object, sourceDAO.Object, opts, sharedOpts);
                await foreach (var result in sourceLiason.FetchAllAsync())
                {
                    Assert.AreEqual(test.Expected.Address, result.Address);
                    Assert.AreEqual(test.Expected.Start, result.Start);
                }
            }
        }

        private static string BasicSlug = "totallyaslug";
        private static DateTime BasicStart = DateTime.Now;
        private static string BasicAddress = "15873525";
    }
}
