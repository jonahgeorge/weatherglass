# Weatherglass

Weatherglass is an open-source Google Analytics alternative that is easy to self-host.

## Development

```sh
# Copy env example and modify as needed
cp .env.example .env

# Install gin
go get github.com/gin-gonic/gin

# Run application
gin run weatherglass.go
```

## Creating a new migration

```sh
./bin/migrate create -ext sql -dir migrations add_device_id_to_events
```
