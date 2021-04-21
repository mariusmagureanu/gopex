Synopsis
========

Gopex is a small project meant as a client towards Pexip while exposing a rest api to manage requests against Pexip.

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
