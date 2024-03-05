package tests

import (
	"context"
	"github.com/Georgy27/go-grpc/week_3/internal/client/db"
	"github.com/brianvoe/gofakeit"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	txMocks "github.com/Georgy27/go-grpc/week_3/internal/client/db/mocks"
	"github.com/Georgy27/go-grpc/week_3/internal/model"
	"github.com/Georgy27/go-grpc/week_3/internal/repository"
	repoMocks "github.com/Georgy27/go-grpc/week_3/internal/repository/mocks"
	"github.com/Georgy27/go-grpc/week_3/internal/service/note"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type noteRepositoryMockFunc func(mc *minimock.Controller) repository.NoteRepository
	type transactionManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.NoteInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.BeerName()
		content = gofakeit.BeerBlg()

		//repoErr = fmt.Errorf("repo error")

		req = &model.NoteInfo{
			Title:   title,
			Content: content,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name                   string
		args                   args
		want                   int64
		err                    error
		noteRepositoryMock     noteRepositoryMockFunc
		transactionManagerMock transactionManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			transactionManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, func(ctx context.Context) error {
					return nil
				}).Return(nil)
				return mock
			},
			noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
				mock := repoMocks.NewNoteRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		//{
		//	name: "service error case",
		//	args: args{
		//		ctx: ctx,
		//		req: req,
		//	},
		//	want: 0,
		//	err:  repoErr,
		//	noteRepositoryMock: func(mc *minimock.Controller) repository.NoteRepository {
		//		mock := repoMocks.NewNoteRepositoryMock(mc)
		//		mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
		//		return mock
		//	},
		//},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteRepoMock := tt.noteRepositoryMock(mc)
			noteTxMock := tt.transactionManagerMock(mc)

			service := note.NewService(noteRepoMock, noteTxMock)
			newID, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
