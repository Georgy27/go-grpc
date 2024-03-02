package service

import (
	"context"

	"github.com/Georgy27/go-grpc/week_3/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
}
