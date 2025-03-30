package grpc

import (
	"context"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/plasmatrip/metriq/internal/types"
	pb "github.com/plasmatrip/metriq/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServer
	repo   storage.Repository
	config config.Config
	lg     logger.Logger
}

func NewGRPCServer(repo storage.Repository, config config.Config, lg logger.Logger) *GRPCServer {
	return &GRPCServer{repo: repo, config: config, lg: lg}
}

func (s *GRPCServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse

	//проверяем тип метрики
	if err := types.CheckMetricType(in.GetType()); err != nil {
		s.lg.Sugar.Infow("error in request handler", "error: ", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//проверяем имя метрики
	if len(in.GetId()) == 0 {
		s.lg.Sugar.Infow("error in request handler", "error: ", "the name of the metric is empty")
		return nil, status.Error(codes.InvalidArgument, "the name of the metric is empty")
	}

	metric, err := s.repo.Metric(ctx, in.GetId())
	if err != nil {
		s.lg.Sugar.Infoln("Metric not found")
		return nil, status.Error(codes.NotFound, "Metric not found")
	}

	respMetric := metric.Convert(in.GetId())

	response.Metric.Id = respMetric.ID
	response.Metric.Value = *respMetric.Value
	response.Metric.Delta = *respMetric.Delta
	response.Metric.Type = respMetric.MType

	return &response, nil
}

func (s *GRPCServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	var response pb.GetMetricsResponse

	metrics, err := s.repo.Metrics(ctx)
	if err != nil {
		s.lg.Sugar.Infow("error in request handler", "error: ", err)
		return nil, status.Error(codes.NotFound, "Metrics not found")
	}

	for _, metric := range metrics {
		respMetric := metric.Convert(metric.MetricType)

		var grpcMetric pb.Metric
		grpcMetric.Id = respMetric.ID
		grpcMetric.Value = *respMetric.Value
		grpcMetric.Delta = *respMetric.Delta
		grpcMetric.Type = respMetric.MType

		response.Metrics = append(response.Metrics, &grpcMetric)
	}

	return &response, nil
}

func (s *GRPCServer) UpdatesMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var metrics []models.Metrics
	var response pb.UpdateMetricsResponse

	for _, metric := range in.GetMetrics() {
		//проверяем тип метрики
		if err := types.CheckMetricType(metric.Type); err != nil {
			s.lg.Sugar.Infow("error in request handler", "error: ", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		//проверяем имя метрики
		if len(metric.Id) == 0 {
			s.lg.Sugar.Infow("error in request handler", "error: ", "the name of the metric is empty")
			return nil, status.Error(codes.InvalidArgument, "the name of the metric is empty")
		}

		metrics = append(metrics, models.Metrics{
			ID:    metric.Id,
			MType: metric.Type,
			Value: &metric.Value,
			Delta: &metric.Delta,
		})
	}

	if err := s.repo.SetMetrics(ctx, metrics); err != nil {
		s.lg.Sugar.Infow("error in request handler", "error: ", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.lg.Sugar.Infoln("Metrics updated, updated metrics count", len(metrics))

	return &response, nil
}

func (s *GRPCServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var metric *pb.Metric
	var response pb.UpdateMetricResponse

	metric = in.GetMetric()

	//проверяем тип метрики
	if err := types.CheckMetricType(metric.Type); err != nil {
		s.lg.Sugar.Infow("error in request handler", "error: ", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//проверяем имя метрики
	if len(metric.Id) == 0 {
		s.lg.Sugar.Infow("error in request handler", "error: ", "the name of the metric is empty")
		return nil, status.Error(codes.InvalidArgument, "the name of the metric is empty")
	}

	var value any
	switch metric.Type {
	case types.Counter:
		value = metric.Delta
	case types.Gauge:
		value = metric.Value
	}

	if err := s.repo.SetMetric(ctx, metric.Id, types.Metric{MetricType: metric.Type, Value: value}); err != nil {
		s.lg.Sugar.Infow("error in request handler", "error: ", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.lg.Sugar.Infoln("Metric updated: ", metric.Id)

	return &response, nil
}
