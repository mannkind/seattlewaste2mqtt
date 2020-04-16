FROM mcr.microsoft.com/dotnet/core/sdk:3.1 as build
WORKDIR /src
COPY . .
RUN if [ ! -d output/$BUILDPLATFORM ]; then dotnet build -o output/$BUILDPLATFORM -c Release SeattleWaste; fi \
    && cp -r output/$BUILDPLATFORM archoutput

FROM mcr.microsoft.com/dotnet/core/runtime:3.1 AS runtime
COPY --from=build /src/archoutput .
ENTRYPOINT ["./SeattleWaste"]
