package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService(t *testing.T) {
	err := logger.InitLogger("INFO")
	require.NoError(t, err)
	storage := memorystorage.New()
	config := config.Config{
		Logger: config.LoggerConf{Level: "INFO"},
		Server: config.ServerConf{
			Host: "localhost",
			Port: "80",
		},
		DBType: config.DBType{
			Type: "memory",
		},
		GRPCConf: config.GRPCConf{
			Port: "8080",
		},
	}
	application := app.New(storage, config)
	service := NewService(application)

	lis, err := net.Listen("tcp", ":8080")
	require.NoError(t, err)

	srv := grpc.NewServer()

	pb.RegisterEventServiceServer(srv, service)

	go func() {
		if err := srv.Serve(lis); err != nil {
			return
		}
	}()

	require.NoError(t, err)

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() {
		conn.Close()
	}()

	client := pb.NewEventServiceClient(conn)

	md := metadata.New(nil)
	md.Append("request_id", "1")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	in := &pb.EventRequest{
		Event: &pb.Event{
			Id:           "",
			Title:        "Mur",
			DateTime:     timestamppb.Now(),
			Duration:     100,
			Description:  "Mur",
			UserId:       "1",
			NotifiedTime: 1000,
		},
	}

	_, err = client.Add(ctx, in)
	require.NoError(t, err)
	_, err = client.Add(ctx, in)
	require.NoError(t, err)

	list, err := storage.List()
	require.NoError(t, err)
	require.Len(t, list, 2)

	changedID := list[0].ID
	in.Event.Id = changedID
	in.Event.Title = "Meow"
	_, err = client.Update(ctx, in)
	require.NoError(t, err)

	events, err := storage.List()
	require.NoError(t, err)
	for _, event := range events {
		if event.ID == changedID {
			require.Equal(t, event.Title, in.Event.Title)
		} else if event.ID == "Mur" {
			require.Equal(t, event.Title, "Mur")
		}
	}

	client.Delete(ctx, in)
	events, err = storage.List()
	require.Len(t, events, 1)

	require.NoError(t, err)
}
