package backup

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
)

type Backup struct {
	cfg  config.Config
	stor storage.Repository
	lg   *logger.Logger
}

func NewBackup(cfg config.Config, stor storage.Repository, lg *logger.Logger) (*Backup, error) {
	dir := filepath.Dir(cfg.FileStoragePath)
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}

	return &Backup{
		cfg:  cfg,
		stor: stor,
		lg:   lg,
	}, nil
}

func (bkp Backup) Start() {
	if bkp.cfg.Restore {
		if err := bkp.load(); err != nil {
			bkp.lg.Sugar.Fatalw("error loading from backup: ", err)
		}
	}

	if bkp.cfg.StoreInterval == 0 {
		go func() {
			c := make(chan struct{})
			defer close(c)

			bkp.stor.SetBackup(c)
			select {
			case <-c:
				bkp.Save()
			default:
			}
			// for {
			// 	if <-c {
			// 		bkp.Save()
			// 	}
			// }
		}()
	} else {
		go func() {
			ticker := time.NewTicker(time.Duration(bkp.cfg.StoreInterval) * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				err := bkp.Save()
				if err != nil {
					bkp.lg.Sugar.Infow("error saving to backup: ", err)
				}
			}
		}()
	}

}

func (bkp Backup) Save() error {
	file, err := os.OpenFile(bkp.cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for mName, metric := range bkp.stor.Metrics() {
		err := encoder.Encode(metric.Convert(mName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (bkp Backup) load() error {
	var jMetric models.Metrics
	var value any

	file, err := os.OpenFile(bkp.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		if err := decoder.Decode(&jMetric); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch jMetric.MType {
		case types.Counter:
			value = *jMetric.Delta
		case types.Gauge:
			value = *jMetric.Value
		}

		if err := bkp.stor.SetMetric(jMetric.ID, types.Metric{MetricType: jMetric.MType, Value: value}); err != nil {
			return err
		}
	}
	return nil
}
