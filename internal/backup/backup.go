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
	c config.Config
	s storage.Repository
	l *logger.Logger
}

func NewBackup(c config.Config, s storage.Repository, l *logger.Logger) (*Backup, error) {
	dir := filepath.Dir(c.FileStoragePath)
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}

	// file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	return nil, err
	// }

	return &Backup{
		// file: file,
		c: c,
		s: s,
		l: l,
	}, nil
}

func (bkp Backup) Start() {
	if bkp.c.Restore {
		if err := bkp.load(); err != nil {
			bkp.l.Sugar.Fatalw("error saving to backup: ", err)
		}
	}

	if bkp.c.StoreInterval == 0 {
		go func() {
			c := make(chan bool)
			defer close(c)

			bkp.s.SetBackup(c)
			for {
				if <-c {
					// mu := sync.Mutex{}
					// mu.Lock()
					bkp.Save()
					// mu.Unlock()
					// fmt.Println("save")
				}
			}
		}()
	} else {
		ticker := time.NewTicker(time.Duration(bkp.c.StoreInterval) * time.Second)
		go func() {
			for range ticker.C {
				err := bkp.Save()
				if err != nil {
					bkp.l.Sugar.Infow("error saving to backup: ", err)
				}
			}
		}()
	}

}

func (bkp Backup) Save() error {
	// mu := sync.Mutex{}
	// mu.Lock()
	// defer mu.Unlock()

	file, err := os.OpenFile(bkp.c.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for mName, metric := range bkp.s.Metrics() {
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

	file, err := os.OpenFile(bkp.c.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
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

		if err := bkp.s.SetMetric(jMetric.ID, types.Metric{MetricType: jMetric.MType, Value: value}); err != nil {
			return err
		}
	}
	return nil
}
