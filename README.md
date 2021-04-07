Synopsis
========

Gopex is an exercise project meant for a better understanding on how [pexmon](https://bitbucket.org/kinlydev/pexmon/src/master/) works towards Pexip and [pex_portal](https://bitbucket.org/kinlydev/pex_portal/src/master/).

Currently still WIP.

Repo structure
--------------

```
├── pkg
│   ├── dbl - Client lib for database ops (gorm based) [postgres|sqlite].
│   ├── ds - Common datastructures that are to be persisted within dbl.
│   ├── errors - Common custom errors package.
│   ├── log - Custom logger package.
├── src
│   ├── api-gw - Web server which exposes a rest interface.
│   ├── pexip - Client components which is responsable for talking to Pexip.
│   ├── main.go - Application entry point.

```
