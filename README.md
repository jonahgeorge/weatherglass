Weatherglass
===

Weatherglass is an open-source Google Analytics alternative that is easy to self-host.

## Developing 

```sh
brew install go forego


go get ./...
forego run 'go run *.go'
```

## Deploying

```sh
./script/deploy
```

## Creating a new migration

```sh
./bin/migrate create -ext sql -dir migrations add_device_id_to_events
```
