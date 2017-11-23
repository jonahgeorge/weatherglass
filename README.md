# Weatherglass

Weatherglass is an open-source Google Analytics alternative that is easy to self-host.

## Development

```sh
source .env

go run weatherglass.go
```

## Creating a new migration

```sh
./bin/migrate create -ext sql -dir migrations add_device_id_to_events
```
