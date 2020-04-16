FROM mcr.microsoft.com/dotnet/core/sdk:3.1 as build
WORKDIR /src
COPY . .
RUN if [ ! -d output/`uname -m` ]; then dotnet build -o output/`uname -m` -c Release SeattleWaste; fi \
    && cp -r output/`uname -m` archoutput

FROM mcr.microsoft.com/dotnet/core/runtime:3.1 AS runtime
COPY --from=build /src/archoutput .
ENTRYPOINT ["./SeattleWaste"]
