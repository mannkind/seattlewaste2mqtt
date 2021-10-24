using System;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using SeattleWaste.DataAccess;
using SeattleWaste.Models.Source;

namespace SeattleWasteTest.DataAccess;

[TestClass]
public class SourceDAOTest
{
    [TestMethod]
    public async Task FetchOneAsync()
    {
        var tests = new[]
        {
                new {
                    Address = "2133 N 61ST ST",
                },
                new {
                    Address = "7022 12th Ave NW",
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
            var result = await dao.FetchOneAsync(address, DateTime.Today);
            Assert.IsNotNull(result);
        }
    }
}

public class LocalSourceDAO : SourceDAO
{
    public LocalSourceDAO(ILogger<LocalSourceDAO> logger, IHttpClientFactory httpClientFactory) :
        base(logger, httpClientFactory)
    {
    }

    public Task<Response> FetchOneAsync(string address, DateTime today, CancellationToken cancellationToken = default)
    {
        return base.FetchAsync(address, today, cancellationToken);
    }
}
