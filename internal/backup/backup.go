package backup

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
)

type Backup struct {
	file *os.File
	s    storage.Repository
	ncdr *json.Encoder
	dcdr *json.Decoder
	l    *logger.Logger
}

func NewBackup(fn string, s storage.Repository, l *logger.Logger) (*Backup, error) {
	dir := filepath.Dir(fn)

	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Backup{
		file: file,
		s:    s,
		ncdr: json.NewEncoder(file),
		dcdr: json.NewDecoder(file),
		l:    l,
	}, nil
}

func (bkp Backup) Start(storeInterval int, restore bool) {
	if restore {
		if err := bkp.load(); err != nil {
			bkp.l.Sugar.Fatalw("error saving to backup: ", err)
		}
	}

	if storeInterval == 0 {
		go func() {
			c := make(chan bool)
			defer close(c)

			bkp.s.Backup(c)
			for {
				if <-c {
					bkp.save()
				}
			}
		}()
	} else {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		go func() {
			for range ticker.C {
				err := bkp.save()
				if err != nil {
					bkp.l.Sugar.Infow("error saving to backup: ", err)
				}
			}
		}()
	}

}

func (bkp Backup) save() error {
	err := bkp.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = bkp.file.Seek(0, 0)
	if err != nil {
		return err
	}

	for mName, metric := range bkp.s.Metrics() {
		err := bkp.ncdr.Encode(metric.Convert(mName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (bkp Backup) load() error {
	var jMetric models.Metrics
	var value any
	for {
		if err := bkp.dcdr.Decode(&jMetric); err == io.EOF {
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
