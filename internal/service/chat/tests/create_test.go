package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/AndreiMartynenko/chat-server/internal/model"
	"github.com/AndreiMartynenko/chat-server/internal/repository"
	repositoryMocks "github.com/AndreiMartynenko/chat-server/internal/repository/mocks"
	chatService "github.com/AndreiMartynenko/chat-server/internal/service/chat"
	"github.com/AndreiMartynenko/common/pkg/db"
	dbMocks "github.com/AndreiMartynenko/common/pkg/db/mocks"
	"github.com/AndreiMartynenko/common/pkg/db/transaction"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type messagesRepositoryMockFunc func(mc *minimock.Controller) repository.MessagesRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type transactorMockFunc func(mc *minimock.Controller) db.Transactor

	type args struct {
		ctx context.Context
		req *model.Chat
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		usernames = []string{"name1", "name2", "name3"}

		repositoryErr = fmt.Errorf("failed to create chat")

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		req = &model.Chat{
			Usernames: usernames,
		}

		reqLog = &model.Log{
			Text: fmt.Sprintf("Created chat with id: %v", id),
		}
	)

	tests := []struct {
		name                   string
		args                   args
		want                   string
		err                    error
		chatRepositoryMock     chatRepositoryMockFunc
		messagesRepositoryMock messagesRepositoryMockFunc
		logRepositoryMock      logRepositoryMockFunc
		transactorMock         transactorMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(minimock.AnyContext, req).Return(id, nil)
				return mock
			},
			messagesRepositoryMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repositoryMocks.NewMessagesRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, reqLog).Return(nil)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.CommitMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "chat repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(minimock.AnyContext, req).Return("", repositoryErr)
				return mock
			},
			messagesRepositoryMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repositoryMocks.NewMessagesRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "log repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(minimock.AnyContext, req).Return(id, nil)
				return mock
			},
			messagesRepositoryMock: func(mc *minimock.Controller) repository.MessagesRepository {
				mock := repositoryMocks.NewMessagesRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, reqLog).Return(repositoryErr)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepositoryMock := tt.chatRepositoryMock(mc)
			messagesRepositoryMock := tt.messagesRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			txManagerMock := transaction.NewTransactionManager(tt.transactorMock(mc))
			srv := chatService.NewService(chatRepositoryMock, messagesRepositoryMock, logRepositoryMock, txManagerMock)

			res, err := srv.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
