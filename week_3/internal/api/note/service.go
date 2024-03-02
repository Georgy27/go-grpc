package note

import (
	"github.com/Georgy27/go-grpc/week_3/internal/service"
	desc "github.com/Georgy27/go-grpc/week_3/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
