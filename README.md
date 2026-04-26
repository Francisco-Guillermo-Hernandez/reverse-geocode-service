# Reverse Geocode Service

## Build for Docker

```bash
    docker compose build
    docker compose up
```

## Usage

```bash
    curl -X GET "http://localhost:8080/geocode?lat=13.701142&lng=-89.2244492"
    # {"country_code":"SV","city":"San Salvador","latitude":13.68935,"longitude":-89.18718,"population":525990,"state":"San Salvador Department","country":"El Salvador"}

    curl -X GET "http://localhost:8080/search-city?cityName=San%20Salvador&countryCode=SV"
    # {"country_code":"SV","city":"San Salvador","latitude":13.68935,"longitude":-89.18718,"population":525990,"state":"San Salvador Department","country":"El Salvador"}
```

## Run local
```bash
    make serve
```

## Build for Lambda

```bash
    make build-function
```

## Deploy
```bash
    cd cloud
    cdk bootstrap --profile Developer
    cdk deploy --profile Developer
```

## Using Lambda

```bash
curl -X GET "https://your-api-id.execute-api.your-region.amazonaws.com/geocode?lat=13.701142&lng=-89.2244492"
curl -X GET "https://your-api-id.execute-api.your-region.amazonaws.com/search-city?cityName=San%20Salvador&countryCode=SV"
```