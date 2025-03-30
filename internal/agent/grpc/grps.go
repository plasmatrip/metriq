package grpc

import (
	"context"

	"github.com/plasmatrip/metriq/internal/agent/config"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
	pb "github.com/plasmatrip/metriq/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Client pb.MetricsClient
	Conn   *grpc.ClientConn
	repo   storage.Repository
	config config.Config
	lg     logger.Logger
}

func NewGRPCClient(address string, repo storage.Repository, cfg config.Config, lg logger.Logger) (*GRPCClient, error) {
	conn, err := grpc.NewClient(":"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewMetricsClient(conn)
	return &GRPCClient{
		Client: client,
		Conn:   conn,
		repo:   repo,
		config: cfg,
		lg:     lg,
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.Conn.Close()
}

func (c *GRPCClient) SendMetrics() error {
	ctx := context.Background()
	metrics, err := c.repo.Metrics(ctx)
	if len(metrics) == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	// convert metrics
	sMetrics := make([]*pb.Metric, 0, len(metrics))
	for mName, metric := range metrics {
		var pbMetric *pb.Metric
		mMetric := metric.Convert(mName)

		switch mMetric.MType {
		case types.Gauge:
			pbMetric = &pb.Metric{
				Id:    mMetric.ID,
				Value: *mMetric.Value,
				Type:  mMetric.MType,
			}
		case types.Counter:
			pbMetric = &pb.Metric{
				Id:    mMetric.ID,
				Delta: *mMetric.Delta,
				Type:  mMetric.MType,
			}
		}

		sMetrics = append(sMetrics, pbMetric)
	}

	c.Client.UpdatesMetrics(ctx, &pb.UpdateMetricsRequest{
		Metrics: sMetrics,
	})

	return nil
}
