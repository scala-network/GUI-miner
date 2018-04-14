// Package main implements a quickly hacked together information
// gatherer to collect information about the Stellite blockchain
// for display on stellite.live
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"bitbucket.org/iliveit/other-projects/stellite-api/src/runner"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	// This is how SQL drivers are imported, golint complains if
	// you don't have an explanation comment here
	_ "github.com/go-sql-driver/mysql"
)

var isRunning uint32
var logger *log.Entry
var waitGroup sync.WaitGroup

// Config represents the configuration for the application
type Config struct {
	LogFormat              string `split_words:"true"`
	LogLevel               string `split_words:"true"`
	DatabaseEndpoint       string `split_words:"true"`
	DatabaseName           string `split_words:"true"`
	DatabaseUsername       string `split_words:"true"`
	DatabasePassword       string `split_words:"true"`
	XtlDaemonEndpoint      string `split_words:"true"`
	PoolSleepSeconds       int    `split_words:"true"`
	ExchangeSleepSeconds   int    `split_words:"true"`
	BlockchainSleepSeconds int    `split_words:"true"`
}

var poolRunner runner.Pools
var exchangeRunner runner.Exchanges
var blockchainRunner runner.Blockchain

func main() {

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "Jan 02 15:04:05",
	})
	if strings.ToLower(config.LogFormat) == "text" {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "Jan 02 15:04:05",
		})
	}
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
	logger = log.WithFields(log.Fields{
		"service": "stellite-api",
	})

	// Setup signal handlers
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	// Remember, on Linux, syscall.SIGKILL can't be caught
	//runWaitGroup.Add(1)
	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGHUP:
			log.Warn("SIGHUP received from OS")
			//handle Reload
		case syscall.SIGINT:
			log.Warn("SIGINT received from OS")
			stop()
		case syscall.SIGTERM:
			log.Warn("SIGTERM received from OS")
			stop()
		}
	}()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseEndpoint,
		config.DatabaseName,
		"charset=utf8&parseTime=True")
	/*
		exchangeRunner := runner.Exchanges{
			ConnectionString: connectionString,
			SleepSeconds:     config.ExchangeSleepSeconds,
			Logger: logger.WithFields(log.Fields{
				"routine": "exchange_sync",
			}),
		}

		blockchainRunner := runner.Blockchain{
			DaemonEndpoint:   config.XtlDaemonEndpoint,
			ConnectionString: connectionString,
			SleepSeconds:     config.BlockchainSleepSeconds,
			Logger: logger.WithFields(log.Fields{
				"routine": "blockchain_sync",
			}),
		}
	*/

	poolRunner = runner.Pools{
		ConnectionString: connectionString,
		SleepSeconds:     config.PoolSleepSeconds,
		Logger: logger.WithFields(log.Fields{
			"routine": "pool_sync",
		}),
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := poolRunner.Run()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	waitGroup.Wait()
	logger.Info("Shutdown")
}

// stop everything!
func stop() {
	poolRunner.Stop()
	//exchangeRunner.Stop()
	//blockchainRunner.Stop()
}
