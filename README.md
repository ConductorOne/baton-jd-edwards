# `baton-jd-edwards` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-jd-edwards.svg)](https://pkg.go.dev/github.com/conductorone/baton-jd-edwards) ![main ci](https://github.com/conductorone/baton-jd-edwards/actions/workflows/main.yaml/badge.svg)

`baton-jd-edwards` is a connector for JD Edwards EnterpriseOne built using the 
[Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the 
JD Edwards EnterpriseOne Application Interface Services (AIS) Server REST APIs 
to sync data about users and roles. Check out 
[Baton](https://github.com/conductorone/baton) to learn more about the project 
in general.

# Getting Started

## Prerequisites

1. JD Edwards EnterpriseOne environment configured with an AIS Server before you can use the AIS Server REST APIs. More info [here](https://docs.oracle.com/cd/E53430_01/EOIIS/toc.htm).
2. AIS server url, JD Edwards username and password. 

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-jd-edwards

BATON_USERNAME=jdeUsername BATON_PASSWORD=jdePassword BATON_AIS_URL=https://your_ais_server:port baton-jd-edwards
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_USERNAME=jdeUsername BATON_PASSWORD=jdePassword BATON_AIS_URL=https://your_ais_server:port baton-jd-edwards ghcr.io/conductorone/baton-jd-edwards:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-jd-edwards/cmd/baton-jd-edwards@main

BATON_USERNAME=jdeUsername BATON_PASSWORD=jdePassword BATON_AIS_URL=https://your_ais_server:port baton-jd-edwards
baton resources
```

# Data Model

`baton-jd-edwards` will pull down information about the following JD Edwards resources:

- Users
- Roles

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-jd-edwards` Command Line Usage

```
baton-jd-edwards

Usage:
  baton-jd-edwards [flags]
  baton-jd-edwards [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --ais-url string         required: Your JD Edwards AIS Server REST API url. Provided url should contain port. (e.g: https://your_ais_server:port). ($BATON_AIS_URL)
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --env string             Environment to use for login. If not specified, the default environment configured for the AIS Server will be used. ($BATON_ENV)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-jd-edwards
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --password string        required: JD Edwards EnterpriseOne password. ($BATON_PASSWORD)
  -p, --provisioning           This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync         This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing              This must be set to enable ticketing support ($BATON_TICKETING)
      --username string        required: JD Edwards EnterpriseOne username. ($BATON_USERNAME)
  -v, --version                version for baton-jd-edwards

Use "baton-jd-edwards [command] --help" for more information about a command.
```