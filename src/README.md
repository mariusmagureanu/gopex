Api-gw
======

Api-gw is a web server which exposes a rest interface. 

This rest api is essentially a wrapper on top of the pexip api. The idea behind is to control conferences and participants.

See Pexip's client api [here](https://docs.pexip.com/api_client/api_rest.htm).

Build and run
-------------

```sh
$ make build
$ ./bin/monitor
```

Usage options
-------------

|       Arg           |Description                                             
|---------------------|-------------------------------------------------------------------------|
|**db-host**		      |Host that points to a PostgreSQL instance                            |     								  
|**db-max-conn-lifetime** |PostgreSQL maximum connection lifetime (default 10m0s)               |
|**db-max-idle-conns**    |PostgreSQL maximum idle connections (default 5)                      |                          
|**db-max-open-conns**    |PostgreSQL maximum open connections (default 20)                     |
|**db-name**			  |Default PostgreSQL database name (default "pexmon")                  |
|**db-port**			  |PostgreSQL port (default 5432)                                       |
|**db-user**              |Default PostgreSQL database user                                     | 	
|**db-pwd** 			  |Password for the PostgreSQL database user                            |
|**host** 				  |Set the Monitor's host. (default "0.0.0.0")   			            |
|**port** 				  |Monitor port. (default 8088)    							            |
|**tls-port**			  |Monitor https port. (default 8443) 						            |
|**log-level**            |Logging level [quiet|debug|info|warning|error] (default "debug")     |
|**pexip-node**           |Pexip node address (default "https://test-join.dev.kinlycloud.net")  |
|**pexip-max-cons**       |Maximum open connections against a Pexip node (default 100)          |
|**pexip-timeout**        |Default timeout for the http client talking with Pexip (default 5s)  |
|**pexip-token-refresh**  |Interval for refreshing Pexip tokens (default 1m0s)                  |
|**sqlite**               |Path to an sqlite database (mutually exclusive with *db-host*)       |
|**V**    				  |Show current version and exit                                        |
|**h**    				  |Show this help and exit                                              |
