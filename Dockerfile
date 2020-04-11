FROM mcr.microsoft.com/dotnet/core/sdk:3.1-alpine as build
WORKDIR /src
COPY . .
RUN dotnet build -c release -o output

FROM mcr.microsoft.com/dotnet/core/runtime:3.1-alpine AS runtime
COPY --from=build /src/output .
ENTRYPOINT ["./SeattleWaste"]
