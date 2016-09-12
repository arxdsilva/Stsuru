[![Build Status](https://travis-ci.org/ArxdSilva/Stsuru.svg?branch=master)](https://travis-ci.org/ArxdSilva/Stsuru)
[![Go Report Card](https://goreportcard.com/badge/github.com/arxdsilva/Stsuru)](https://goreportcard.com/badge/github.com/arxdsilva/Stsuru)
[![codebeat badge](https://codebeat.co/badges/2ffb3187-79c2-4589-a383-86da64440e64)](https://codebeat.co/projects/github-com-arxdsilva-stsuru)


# Stsuru

- **Simple** link 'shortener';
- Written in [Go](http://golang.org);

## Introduction
Implementation of a simple link shortener in Golang. Intended to `hash` & `display` sortened URL's in a 'pure' Golang's server. It uses Gorilla's mux to handle server requests.

## Instalation
### Go Get
The easiest way is to install with go get (**needed Golang 1.7 or later installed**):
```shell
$ go get -u github.com/arxdsilva/Stsuru
```

## Usage
### Tests
```shell
$ go test
```
### Start server
```shell
$ go run main.go
```
Open your **browser** and type:
[`localhost:8080`](http://localhost:8080/)

## Contributing
Any help is appreciated! Use [Issue tracker](https://github.com/arxdsilva/stsuru/issues) to report **any** problem or fill feature requests.

Interested to help development? Fork the project and submit a [Pull Request](https://github.com/arxdsilva/stsuru/pulls).

## LICENSE
Check our [MIT](https://github.com/ArxdSilva/Stsuru/blob/master/LICENSE) license file for more info.

## Credits
Made by **[@arxdsilva](https://twitter.com/arxdsilva)** with great help of [Tsuru team](https://github.com/tsuru/tsuru)!

## Extra
Want to implement a free PaaS & that is Open source? Check [Tsuru](https://github.com/tsuru/tsuru)!
