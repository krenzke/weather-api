# Weather API

This is a simple golang web server that acts as a proxy for the [PirateWeather API](https://pirateweather.net/en/latest/). But rather than searching by latitude/longitude like the pirate weather api requires, it accepts a zipcode instead. That zipcode then gets geocoded via a Google API to turn it into a latitude/longitude. Nothing fancy, but it makes it a little easier when a user only has to type in a zipcode and not a pair of coordinates.

## Development

To run locally

```sh
go run *.go
```

## Building

```sh
go build *.go
```