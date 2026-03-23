# Reverse Geocode Service

## Build

```bash
    docker compose build
    docker compose up
```

## Usage

```bash
    curl -X GET "http://localhost:8080/api/v1/geocode?lat=13.701142&lng=-89.2244492"
    # {"country_code":"SV","city":"San Salvador","latitude":13.68935,"longitude":-89.18718,"population":525990,"state":"San Salvador Department","country":"El Salvador"}
```