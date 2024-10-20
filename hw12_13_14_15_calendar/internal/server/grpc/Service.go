package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	port    string
	storage storage.Storage
	pb.UnimplementedEventServiceServer
	srv grpc.Server
}

func NewService(app *app.App) *Service {
	return &Service{
		port:    app.Config().GRPCConf.Port,
		storage: app.Storage(),
	}
}

func (s *Service) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(func(
		ctx context.Context,
		request interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (response interface{}, err error) {
		start := time.Now()
		response, err = handler(ctx, request)
		duration := time.Since(start)
		method := info.FullMethod

		if err != nil {
			slog.Error(fmt.Sprintf("Request failed: method=%s, request=%v, error=%v, duration=%s, status=%s",
				method, request, err, duration, status.Code(err)))
		} else {
			slog.Info("Request succeeded: method=%s, request=%v, response=%v, duration=%s",
				method, request, response, duration)
		}

		return response, err
	}))

	pb.RegisterEventServiceServer(srv, s)

	if err := srv.Serve(lis); err != nil {
		slog.Error("Failed to serve", "error", err)
	}

	return nil
}

func (s *Service) Stop() error {
	if s != nil {
		s.srv.GracefulStop()
	}

	return nil
}

func requestIDFromContext(ctx context.Context) string {
	requestID := ""

	if ctx != nil {
		meta, ok := metadata.FromIncomingContext(ctx)

		if ok {
			ids := meta.Get("request_id")
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}
	}

	return requestID
}

func getEvent(req *pb.EventRequest) (*entities.Event, error) {
	reqEvent := req.GetEvent()

	if reqEvent == nil {
		return nil, errors.New("reqEvent is nil")
	}

	event := &entities.Event{
		ID:          reqEvent.Id,
		Title:       reqEvent.Title,
		DateTime:    reqEvent.DateTime.AsTime(),
		Duration:    time.Duration(reqEvent.Duration),
		Description: reqEvent.Description,
		OwnerID:     reqEvent.UserId,
		NotifyTime:  reqEvent.NotifiedTime,
	}

	return event, nil
}

func (s *Service) Add(ctx context.Context, req *pb.EventRequest) (*pb.EmptyResponse, error) {
	requestID := requestIDFromContext(ctx)

	slog.Info("Add for request id " + requestID)

	event, err := getEvent(req)
	if err != nil {
		slog.Error("error extracting event", err)
		return nil, status.Error(codes.InvalidArgument, "request event is required")
	}

	err = s.storage.Add(event)
	if err != nil {
		slog.Error("Error adding event to storage: %v", err)
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Service) Update(ctx context.Context, req *pb.EventRequest) (*pb.EmptyResponse, error) {
	requestID := requestIDFromContext(ctx)

	slog.Info("Update for request id " + requestID)

	event, err := getEvent(req)
	if err != nil {
		slog.Error("error extracting event", err)
		return nil, status.Error(codes.InvalidArgument, "request event is required")
	}

	err = s.storage.Change(event)
	if err != nil {
		slog.Error("Error adding event to storage: %v", err)
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Service) Delete(ctx context.Context, req *pb.EventRequest) (*pb.EmptyResponse, error) {
	requestID := requestIDFromContext(ctx)

	slog.Info("Delete for request id " + requestID)

	event, err := getEvent(req)
	if err != nil {
		slog.Error("error extracting event", err)
		return nil, status.Error(codes.InvalidArgument, "request event is required")
	}

	err = s.storage.Delete(event)
	if err != nil {
		slog.Error("Error adding event to storage: %v", err)
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func (s *Service) List(ctx context.Context, le *pb.ListEvents) (*pb.EventsResponse, error) {
	requestID := requestIDFromContext(ctx)

	slog.Info("List for request id " + requestID)

	events, err := s.storage.List()

	switch le.Period {
	case pb.Period_DAY:
		slog.Info("Get day events")
		events = entities.GetTodayEvents(events)
	case pb.Period_WEEK:
		slog.Info("Get weekly events")
		events = entities.GetWeekEvents(events)
	case pb.Period_MONTH:
		slog.Info("Get monthly events")
		events = entities.GetMonthEvents(events)
	default:
		slog.Info("Get all events")
	}

	if err != nil {
		slog.Error("Error adding event to storage: %v", err)
		return nil, err
	}

	protoEvents := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		protoEvent := &pb.Event{
			Id:           event.ID,
			Title:        event.Title,
			DateTime:     timestamppb.New(event.DateTime),
			Duration:     int64(event.Duration.Seconds()),
			Description:  event.Description,
			UserId:       event.OwnerID,
			NotifiedTime: event.NotifyTime,
		}
		protoEvents = append(protoEvents, protoEvent)
	}

	return &pb.EventsResponse{Events: protoEvents}, nil
}
