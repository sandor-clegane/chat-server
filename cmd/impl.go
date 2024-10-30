package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	desc "github.com/sandor-clegane/chat-server/internal/generated/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	chatTableName = "chat"

	chatColumnUsernames = "usernames"
	chatColumnID        = "id"
)

type server struct {
	db *pgxpool.Pool

	desc.UnimplementedChatV1Server
}

// Create ...
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	insertBuilder := sq.Insert(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatColumnUsernames).
		Values(
			pq.Array(req.GetUsernames()),
		).
		Suffix(fmt.Sprintf("RETURNING %s", chatColumnID))

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("error", err))
		return nil, err
	}

	var chatID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		slog.Error("failed to insert chat", slog.Any("error", err))
		return nil, err
	}

	slog.Info("chat inserted", slog.Any("chat_id", chatID))

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

// Delete ...
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	deleteBuilder := sq.Delete(chatTableName).
		Where(sq.Eq{chatColumnID: req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build delete query", slog.Any("error", err))
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to delete chat", slog.Any("error", err))
		return nil, err
	}

	slog.Info("chat deleted", slog.Any("chat_id", req.GetId()))

	return &emptypb.Empty{}, nil
}

// SendMessage ...
func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	// TODO: will be implemented later
	log.Printf("User: [%s] sent message [%s] at time [%s]", req.GetFrom(), req.Text, req.GetTimestamp().AsTime().String())

	return &emptypb.Empty{}, nil
}
