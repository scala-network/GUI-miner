package runner

import "github.com/sirupsen/logrus"

// Blockchain runner retrieves information from the blockchain
type Blockchain struct {
	// DaemonEndpoint is the blockchain daemon's endpoint used to retrieve
	// the information from
	DaemonEndpoint string
	// ConnectionString is the database connection string
	ConnectionString string
	// SleepSeconds determines the sleep interval between fetching data
	SleepSeconds int
	// Logger for the runner
	Logger *logrus.Entry
}

// Run the blockchain information retriever
func (chain *Blockchain) Run() error {

	/*
		logger.WithFields(log.Fields{
			"daemon_endpoint": config.XtlDaemonEndpoint,
			"routine":         "chain",
		}).Info("Run")
		atomic.StoreUint32(&isRunning, 1)

		daemon := connectors.Daemon{
			Endpoint: config.XtlDaemonEndpoint,
		}
		for atomic.LoadUint32(&isRunning) == 1 {
			// Get last block height from database
			db, err := gorm.Open("mysql", connectionString)
			if err != nil {
				logger.Fatalf("Unable to connect to database: %s", err)
			}
			defer db.Close()

			// Since I' hacking this quickly and using gotos, declare the vars here
			var blockInfo connectors.GetBlockResponse

			var lastBlock models.Block
			query := db.Last(&lastBlock)
			if query.Error != nil {
				if query.Error == gorm.ErrRecordNotFound {
					// If none, start at 0
					logger.WithFields(log.Fields{
						"err": query.Error,
					}).Warning("No block information found - starting from height 0")
				} else {
					logger.WithFields(log.Fields{
						"err":     query.Error,
						"routine": "chain",
					}).Error("Unable to query database")
					goto skip
				}
			} else {
				lastBlock.Height++
			}

			// Keep reading the chain until we are done
			for atomic.LoadUint32(&isRunning) == 1 {
				logger.WithFields(log.Fields{
					"height":  lastBlock.Height,
					"routine": "chain",
				}).Info("Reading blockchain")

				// Fetch information for the given height
				blockInfo, err = daemon.GetBlockInfo(lastBlock.Height)
				if err != nil {
					logger.WithFields(log.Fields{
						"height":  lastBlock.Height,
						"err":     err,
						"routine": "chain",
					}).Warning("Unable to get block info from daemon")
					goto skip
				}
				if blockInfo.Error.Message != "" {
					logger.WithFields(log.Fields{
						"height":  lastBlock.Height,
						"err":     blockInfo.Error.Message,
						"routine": "chain",
					}).Warning("End of chain")
					goto skip
				}
				// Only save non orphans
				if blockInfo.Result.BlockHeader.OrphanStatus == false {
					lastBlock = models.Block{
						Height:     blockInfo.Result.BlockHeader.Height,
						Difficulty: blockInfo.Result.BlockHeader.Difficulty,
						Reward:     float64(blockInfo.Result.BlockHeader.Reward) / float64(100.00),
						Timestamp:  time.Unix(blockInfo.Result.BlockHeader.Timestamp, 0).UTC(),
						TxCount:    blockInfo.Result.BlockHeader.NumTxes,
					}
					query = db.Save(&lastBlock)
					if query.Error != nil {
						logger.WithFields(log.Fields{
							"err":     query.Error,
							"routine": "chain",
						}).Error("Unable to save block to database")
						goto skip
					}
				}
				lastBlock.Height++
			}

		skip:
			db.Close()
			logger.WithFields(log.Fields{
				"routine": "chain",
				"seconds": config.BlockchainSleepSeconds,
			}).Info("Sleeping")
			time.Sleep(time.Second * time.Duration(config.BlockchainSleepSeconds))
		}

	*/
	return nil
}

// Stop syncing
func (chain *Blockchain) Stop() {
	chain.Logger.Info("Stopping blockchain runner")
}
