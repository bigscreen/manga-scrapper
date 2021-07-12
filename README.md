# Manga-Scrapper (Mangajack)

Web service which provides non-official (pirate) manga content.

## Requirements

- Go 1.16
- Chrome webkit

## Project Setup

Clone this repo inside `$GOPATH/src/github.com/bigscreen/`, then set up the project by running the following commands:

```shell
make copy-config
make compile
```

## Running Tests

Run the following command:

```shell
make test
```

or run the following command to show coverage per package:

```shell
make test-cov
```

## Running Service

Run the following command to start this service *(ensure the setup commands have been executed)*:

```shell
make run-server
```

or can be run manually:

```shell
make compile
./out/mangajack start
```
