package note

import (
	"github.com/Georgy27/go-grpc/week_3/internal/client/db"
	//"github.com/Georgy27/go-grpc/week_3/internal/client/db"
	"github.com/Georgy27/go-grpc/week_3/internal/repository"
	"github.com/Georgy27/go-grpc/week_3/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}

func NewMockService(deps ...interface{}) service.NoteService {
	srv := serv{}
	for _, v := range deps {

		switch s := v.(type) {
		case repository.NoteRepository:
			srv.noteRepository = s

		case db.TxManager:
			srv.txManager = s

		}
	}

	return &srv
}
