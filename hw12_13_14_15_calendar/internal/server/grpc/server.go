package grpcserver

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	pb "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/api/gen"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/app"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/config"
	logger "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	storage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	grpcServer *grpc.Server
	conf       config.ServerConfig
	app        *app.App
	logger     *logger.Logger
}

type CalendarServer struct {
	pb.UnimplementedCalendarServiceServer
	app    *app.App
	logger *logger.Logger
}

func NewCalendarServer(app *app.App, logger *logger.Logger) *CalendarServer {
	return &CalendarServer{
		app:    app,
		logger: logger,
	}
}

func (s *CalendarServer) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.Event, error) {
	if err := validateCreateEventRequest(req); err != nil {
		return nil, err
	}

	event := &storage.Event{
		Title:       req.Title,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
		Description: req.Description,
		UserID:      req.UserId,
		NotifyAt:    req.NotifyAt.AsTime(),
	}

	createdEvent, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create event: %v", err))
		return nil, status.Error(codes.Internal, "failed to create event")
	}

	return convertDomainToProtoEvent(createdEvent), nil
}

func (s *CalendarServer) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.Event, error) {
	if req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event cannot be nil")
	}

	event := &storage.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		StartTime:   req.Event.StartTime.AsTime(),
		EndTime:     req.Event.EndTime.AsTime(),
		Description: req.Event.Description,
		UserID:      req.Event.UserId,
		NotifyAt:    req.Event.NotifyAt.AsTime(),
	}

	updatedEvent, err := s.app.UpdateEvent(ctx, event.ID, event)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update event: %v", err))
		return nil, status.Error(codes.Internal, "failed to update event")
	}

	return convertDomainToProtoEvent(updatedEvent), nil
}

func (s *CalendarServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete event: %v", err))
		return nil, status.Error(codes.Internal, "failed to delete event")
	}

	return &emptypb.Empty{}, nil
}

func (s *CalendarServer) ListEventsForDay(ctx context.Context,
	req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	startTime := req.Date.AsTime()
	endTime := startTime.Add(24 * time.Hour)

	events, err := s.app.ListEvents(ctx, req.UserId, startTime, endTime)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		return nil, status.Error(codes.Internal, "failed to list events")
	}

	return &pb.ListEventsResponse{
		Events: convertDomainToProtoEvents(events),
	}, nil
}

func (s *CalendarServer) ListEventsForWeek(ctx context.Context,
	req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	startTime := req.Date.AsTime()
	endTime := startTime.Add(7 * 24 * time.Hour)

	events, err := s.app.ListEvents(ctx, req.UserId, startTime, endTime)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		return nil, status.Error(codes.Internal, "failed to list events")
	}

	return &pb.ListEventsResponse{
		Events: convertDomainToProtoEvents(events),
	}, nil
}

func (s *CalendarServer) ListEventsForMonth(ctx context.Context,
	req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	startTime := req.Date.AsTime()
	endTime := startTime.AddDate(0, 1, 0)

	events, err := s.app.ListEvents(ctx, req.UserId, startTime, endTime)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		return nil, status.Error(codes.Internal, "failed to list events")
	}

	return &pb.ListEventsResponse{
		Events: convertDomainToProtoEvents(events),
	}, nil
}

// Helper functions

func validateCreateEventRequest(req *pb.CreateEventRequest) error {
	if req.Title == "" {
		return status.Error(codes.InvalidArgument, "title cannot be empty")
	}
	if req.StartTime == nil {
		return status.Error(codes.InvalidArgument, "start time cannot be nil")
	}
	if req.EndTime == nil {
		return status.Error(codes.InvalidArgument, "end time cannot be nil")
	}
	if req.StartTime.AsTime().After(req.EndTime.AsTime()) {
		return status.Error(codes.InvalidArgument, "end time must be after start time")
	}
	if req.UserId <= 0 {
		return status.Error(codes.InvalidArgument, "invalid user id")
	}
	return nil
}

func convertDomainToProtoEvent(event *storage.Event) *pb.Event {
	return &pb.Event{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		UserId:      event.UserID,
		NotifyAt:    timestamppb.New(event.NotifyAt),
	}
}

func convertDomainToProtoEvents(events []*storage.Event) []*pb.Event {
	result := make([]*pb.Event, len(events))
	for i, event := range events {
		result[i] = convertDomainToProtoEvent(event)
	}
	return result
}

func NewGrpcServer(logger *logger.Logger, app *app.App,
	conf config.ServerConfig) app.Server {
	return &GrpcServer{
		conf:   conf,
		app:    app,
		logger: logger,
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.conf.Listener.Host, strconv.Itoa(s.conf.Listener.Port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(s.loggingInterceptor),
	)

	calendarServer := NewCalendarServer(s.app, s.logger)
	pb.RegisterCalendarServiceServer(s.grpcServer, calendarServer)

	reflection.Register(s.grpcServer)

	s.logger.Info(fmt.Sprintf("Starting gRPC server on %s", addr))

	go func() {
		<-ctx.Done()
		s.Stop(context.Background())
	}()

	if err := s.grpcServer.Serve(listener); err != nil {
		s.logger.Error(fmt.Sprintf("gRPC server failed: %v", err))
	}

	return nil
}

func (s *GrpcServer) Stop(ctx context.Context) error {
	s.logger.Info("Stopping gRPC server...")

	// Create a channel to signal completion
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or context deadline
	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Forced to stop gRPC server due to timeout")
		s.grpcServer.Stop()
		return ctx.Err()
	}
}

// Logging interceptor for all gRPC methods.
func (s *GrpcServer) loggingInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Log the incoming request
	s.logger.Info(fmt.Sprintf("gRPC request received: method: %s, time:%s",
		info.FullMethod, start.Format(time.RFC3339)),
	)

	// Execute the handler
	resp, err := handler(ctx, req)

	// Log the completion
	duration := time.Since(start)
	if err != nil {
		s.logger.Error(fmt.Sprintf("gRPC request failed: method: %s, time:%s, err %v",
			info.FullMethod, duration.String(), err))
	} else {
		s.logger.Info(fmt.Sprintf("gRPC request completed: method: %s, time:%s",
			info.FullMethod, duration.String()))
	}

	return resp, err
}
