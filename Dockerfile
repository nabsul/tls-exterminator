FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build-env
WORKDIR /App

COPY . ./
RUN dotnet restore TlsExterminator.csproj
RUN dotnet publish -c Release -o out TlsExterminator.csproj

FROM mcr.microsoft.com/dotnet/aspnet:9.0
WORKDIR /App
COPY --from=build-env /App/out .
ENTRYPOINT ["dotnet", "TlsExterminator.dll"]