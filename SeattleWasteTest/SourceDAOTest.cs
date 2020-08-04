using System;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using SeattleWaste.DataAccess;
using SeattleWaste.Models.Source;

namespace SeattleWasteTest
{
    [TestClass]
    public class SourceDAOTest
    {
        [TestMethod]
        public async Task FetchOneAsync()
        {
            var testAddress = "2133 N 61ST ST";
            var tests = new[]
            {
                new {
                    Address = testAddress,
                    Today = new DateTime(2020, 6, 13),
                    Start =  ((DateTimeOffset)new DateTime(2020, 6, 12)).ToUnixTimeSeconds(),
                    Expected = new DateTime(2020, 6, 19)
                },
                new {
                    Address = testAddress,
                    Today = new DateTime(2020, 6, 26),
                    Start = ((DateTimeOffset)new DateTime(2020, 6, 26)).ToUnixTimeSeconds(),
                    Expected = new DateTime(2020, 7, 3)
                }
            };

            var logger = new Mock<ILogger<LocalSourceDAO>>();
            var httpClientFactory = new Mock<IHttpClientFactory>();
            httpClientFactory
                .Setup(x => x.CreateClient(It.IsAny<string>()))
                .Returns(new HttpClient());

            var dao = new LocalSourceDAO(logger.Object, httpClientFactory.Object);

            foreach (var test in tests)
            {
                var address = test.Address;
                var today = ((DateTimeOffset)test.Today).ToUnixTimeSeconds();
                var start = test.Start;
                var result = await dao.FetchOneAsync(address, today, start);
                Assert.AreEqual(test.Expected, result.Start);
            }
        }
    }

    public class LocalSourceDAO : SourceDAO
    {
        public LocalSourceDAO(ILogger<LocalSourceDAO> logger, IHttpClientFactory httpClientFactory) :
            base(logger, httpClientFactory)
        {
        }

        public Task<FetchResponse> FetchOneAsync(string address, long today, long start, CancellationToken cancellationToken = default)
        {
            return base.FetchAsync(address, today, start, cancellationToken);
        }
    }
}
