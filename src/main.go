package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mariusmagureanu/gopex/api-gw/mux"
	"github.com/mariusmagureanu/gopex/pexip"
	"github.com/mariusmagureanu/gopex/pkg/dbl"
	"github.com/mariusmagureanu/gopex/pkg/log"
)

var (
	commandLine  = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port         = commandLine.Int("port", 8088, "Monitor port.")
	httpsPort    = commandLine.Int("tls-port", 8443, "Monitor https port.")
	host         = commandLine.String("host", "0.0.0.0", "Set the Monitor's host.")
	versionFlag  = commandLine.Bool("V", false, "Show version and exit")
	logLevelFlag = commandLine.String("log-level", "debug", "Logging level [quiet|debug|info|warning|error]")

	postgresHostFlag        = commandLine.String("db-host", "localhost", "PostgreSQL host")
	postgresPortFlag        = commandLine.Int("db-port", 5432, "PostgreSQL port")
	postgresUserFlag        = commandLine.String("db-user", "foo", "Default PostgreSQL database user")
	postgresUserPwdFlag     = commandLine.String("db-pwd", "", "Password for the PostgreSQL database user")
	postgresDBFlag          = commandLine.String("db-name", "pexmon", "Default PostgreSQL database name")
	postgresMaxIdleConns    = commandLine.Int("db-max-idle-conns", 5, "PostgreSQL maximum idle connections")
	postgresMaxOpenConns    = commandLine.Int("db-max-open-conns", 20, "PostgreSQL maximum open connections")
	postgresMaxConnLifetime = commandLine.Duration("db-max-conn-lifetime", 10*time.Minute, "PostgreSQL maximum connection lifetime")

	pexipNode                 = commandLine.String("pexip-node", "https://test-join.dev.kinlycloud.net", "Pexip node address")
	pexipClientTimeout        = commandLine.Duration("pexip-timeout", 5*time.Second, "Default timeout for the http client talking with Pexip")
	pexipMaxConns             = commandLine.Int("pexip-max-conns", 100, "Maximum open connections against a Pexip node")
	pexipTokenRefreshInterval = commandLine.Duration("pexip-token-refresh", 60*time.Second, "Interval for refreshing Pexip tokens")

	sqliteDBFlag = commandLine.String("sqlite", "", "Path to an sqlite database")

	version  = "N/A"
	revision = "N/A"

	err error

	db = dbl.DAO{}
)

func initLogger() {
	log.InitNewLogger(os.Stdout, log.ErrorLevel)
	log.SetLogLevel(log.GetLogLevelID(*logLevelFlag))
}

func startHTTPServer() error {
	router, err := mux.InitMux()

	if err != nil {
		return err
	}

	server := &http.Server{Handler: router}
	server.IdleTimeout = 1 * time.Minute
	server.ReadHeaderTimeout = 1 * time.Minute

	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		return err
	}

	log.InfoSync("gopex about to start listening on", listener.Addr())

	return server.Serve(listener)
}

func initPostgresql() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		*postgresHostFlag, *postgresUserFlag, *postgresUserPwdFlag, *postgresDBFlag, *postgresPortFlag)

	log.Debug("connecting to postgresql with", dsn)
	err := db.InitPostgres(dsn, *postgresMaxIdleConns, *postgresMaxOpenConns, *postgresMaxConnLifetime)

	return err
}

func initSqlite() error {
	log.Debug("connecting to sqlite with", *sqliteDBFlag)

	err := db.InitSqlite(*sqliteDBFlag)

	return err
}

func main() {

	commandLine.Usage = func() {
		fmt.Println("Pexip Monitor usage:")
		commandLine.PrintDefaults()
	}

	if err = commandLine.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *versionFlag {
		fmt.Println("Version:  " + version)
		fmt.Println("Revision: " + revision)
		os.Exit(0)
	}

	initLogger()

	//TODO: uncomment this if you actually have a database
	/*
			if *sqliteDBFlag != "" {
				err = initSqlite()
			} else {
				err = initPostgresql()
			}


		if err != nil {
			log.ErrorSync(err)
			os.Exit(1)
		}
	*/

	err = pexip.InitPexipClient(*pexipNode, *pexipClientTimeout, *pexipMaxConns, *pexipMaxConns, *pexipTokenRefreshInterval)

	if err != nil {
		log.ErrorSync(err)
		os.Exit(1)
	}

	err = startHTTPServer()

	if err != nil {
		log.ErrorSync(err)
		os.Exit(1)
	}
}
