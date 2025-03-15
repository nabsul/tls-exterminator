FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build-env
WORKDIR /App

COPY . ./
RUN apt-get update
RUN apt-get install -y clang zlib1g-dev
RUN dotnet restore TlsExterminator.csproj
RUN dotnet publish -r linux-x64 -c Release -o out TlsExterminator.csproj

FROM mcr.microsoft.com/dotnet/aspnet:9.0
WORKDIR /App
COPY --from=build-env /App/out .
ENTRYPOINT ["./TlsExterminator"]
