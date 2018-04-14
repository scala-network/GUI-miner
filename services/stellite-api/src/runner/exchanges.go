package runner

import (
	"github.com/sirupsen/logrus"
)

// Exchanges retrieves trading information from exchanges and stores it
// in the database
type Exchanges struct {
	// ConnectionString is the database connection string
	ConnectionString string
	// SleepSeconds determines the sleep interval between fetching data
	SleepSeconds int
	// Logger for the runner
	Logger *logrus.Entry
}

// Run retrieves the exchange information at specified intevals
func (ex *Exchanges) Run() error {
	/*	var exchanges []connectors.Exchange

		// TODO: Fetch exchanges from the database

		exchanges = append(exchanges, &connectors.Crex24{
			Endpoint: config.ExchangeCrex24Endpoint,
		})
		exchanges = append(exchanges, &connectors.TradeOgre{
			Endpoint: config.ExchangeTradeogreEndpoint,
		})

		atomic.StoreUint32(&isRunning, 1)
		for atomic.LoadUint32(&isRunning) == 1 {

			db, err := gorm.Open("mysql", connectionString)
			if err != nil {
				logger.Fatalf("Unable to connect to database: %s", err)
			}
			defer db.Close()

			for _, exchange := range exchanges {
				logger.WithFields(log.Fields{
					"exchange": exchange.GetName(),
					"routine":  "exchange_sync",
				}).Info("Fetching trade information")

				ticker, err := exchange.GetTicker()
				if err != nil {
					logger.WithFields(log.Fields{
						"exchange": exchange.GetName(),
						"err":      err,
						"routine":  "exchange_sync",
					}).Warning("Unable to fetch trade information")
					goto skip
				}

				price := models.Price{
					Exchange:     exchange.GetName(),
					High:         ticker.High,
					Low:          ticker.Low,
					Last:         ticker.Last,
					Volume:       ticker.VolumeBTC,
					DateCaptured: time.Now().UTC(),
				}
				query := db.Save(&price)
				if err != nil {
					logger.WithFields(log.Fields{
						"exchange": exchange.GetName(),
						"err":      query.Error,
						"routine":  "exchange_sync",
					}).Warning("Unable to save trade information")
					goto skip
				}
			}

		skip:
			db.Close()
			logger.WithFields(log.Fields{
				"routine": "exchange_sync",
				"seconds": config.ExchangeSleepSeconds,
			}).Info("Sleeping")
			time.Sleep(time.Second * time.Duration(config.ExchangeSleepSeconds))
		}

	*/

	return nil
}

// Stop fetching data
func (ex *Exchanges) Stop() {
	ex.Logger.Info("Stopping exchange runner")
}
