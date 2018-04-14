package runner

import (
	"errors"
	"sync/atomic"
	"time"

	"bitbucket.org/iliveit/other-projects/stellite-api/src/connectors"
	"bitbucket.org/iliveit/other-projects/stellite-api/src/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Pools retrieves all the pool information and stores it in the database
type Pools struct {
	// ConnectionString is the database connection string
	ConnectionString string
	// SleepSeconds determines the sleep interval between fetching data
	SleepSeconds int
	// Logger for the runner
	Logger *logrus.Entry

	// isRunning keeps the routine running
	isRunning uint32
}

// Run the fetching routine
func (p *Pools) Run() error {
	p.Logger.Debug("Fetching pools from database")

	db, err := gorm.Open("mysql", p.ConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	atomic.StoreUint32(&p.isRunning, 1)
	for atomic.LoadUint32(&p.isRunning) == 1 {
		var poolModels []models.Pool
		query := db.Where("is_enabled = 1").Find(&poolModels)
		if query.Error != nil {
			p.Logger.Error(query.Error)
			return query.Error
		}

		for _, poolModel := range poolModels {
			stats, err := p.GetPoolStats(poolModel)
			if err != nil {
				p.Logger.WithFields(logrus.Fields{
					"api_type": poolModel.APIType,
					"name":     poolModel.Name,
				}).Error(err)
				continue
			}

			p.Logger.WithFields(logrus.Fields{
				"hashrate":   stats.Hashrate,
				"height":     stats.Height,
				"miners":     stats.Miners,
				"last_block": stats.LastBlockTime,
			}).Info("Stats retrieved")

			poolModel.Miners = stats.Miners
			poolModel.Hashrate = stats.Hashrate
			poolModel.LastBlock = stats.LastBlockTime
			poolModel.LastUpdate = time.Now()

			query = db.Save(&poolModel)
			if query.Error != nil {
				p.Logger.WithFields(logrus.Fields{
					"api_type": poolModel.APIType,
					"name":     poolModel.Name,
				}).Error(err)
				continue
			}
		}
		time.Sleep(time.Millisecond * 500)

	}

	return err
}

// GetPoolStats gets the pool stats
func (p *Pools) GetPoolStats(poolModel models.Pool) (connectors.PoolStats, error) {
	switch poolModel.APIType {
	case "cryptonote-pool":
		p.Logger.WithFields(logrus.Fields{
			"api_type": poolModel.APIType,
			"name":     poolModel.Name,
		}).Info("Fetching info for pool")

		connector := connectors.CryptonotePool{
			Endpoint: poolModel.Endpoint,
		}
		return connector.GetStats()
	case "nodejs-pool":
		p.Logger.WithFields(logrus.Fields{
			"api_type": poolModel.APIType,
			"name":     poolModel.Name,
		}).Info("Fetching info for pool")

		connector := connectors.NodeJSPool{
			Endpoint: poolModel.Endpoint,
		}
		return connector.GetStats()
	}
	return connectors.PoolStats{}, errors.New("Invalid API Type")
}

// Stop fetching data
func (p *Pools) Stop() {
	p.Logger.Info("Stopping pool runner")
	atomic.StoreUint32(&p.isRunning, 0)
}
