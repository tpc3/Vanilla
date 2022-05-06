# Vanilla
[![Go Report Card](https://goreportcard.com/badge/github.com/tpc3/vanilla)](https://goreportcard.com/report/github.com/tpc3/vanilla)
[![Docker Image CI](https://github.com/tpc3/Vanilla/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/Vanilla/actions/workflows/docker-image.yml)
<!-- [![Go](https://github.com/tpc3/Vanilla/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/Vanilla/actions/workflows/go.yml) -->

Discord Bot to collect statistics of custom emojis.

[japanese](https://github.com/tpc3/Vanilla/README-ja.md)

## Use
### Simple
1. ~~Download binary from [Releases](https://github.com/tpc3/Vanilla/releases)~~
    - ~~May binary is in `artifact.zip`~~
    - ~~Want latest? You can download from [Actions](https://github.com/tpc3/Vanilla/actions/workflows/go.yml)~~
    - Sorry, CI which make binary is currently stopped. Use docker or build yourself.
1. [Download config.yaml](https://raw.githubusercontent.com/tpc3/Vanilla/master/config.yaml)
1. Enter your token to config.yaml
1. `./Vanilla`
1. use `emoji.sync` command in your guild

### Docker
1. [Download config.yaml](https://raw.githubusercontent.com/tpc3/Vanilla/master/config.yaml)
1. Enter your token to config.yaml
1. `docker run --rm -it -v $(PWD):/data ghcr.io/tpc3/vanilla`
1. use `emoji.sync` command in your guild

## Build
1. Clone this repository
1. `go build`
### required
- git
- golang
- gcc

## Contribute
Any contribute is welcome.  
You can use Issue and Pull Requests.  