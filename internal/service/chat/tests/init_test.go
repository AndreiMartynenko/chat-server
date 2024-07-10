package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/AndreiMartynenko/chat-server/internal/repository"
	repositoryMocks "github.com/AndreiMartynenko/chat-server/internal/repository/mocks"
	chatService "github.com/AndreiMartynenko/chat-server/internal/service/chat"
	"github.com/AndreiMartynenko/common/pkg/db"
	dbMocks "github.com/AndreiMartynenko/common/pkg/db/mocks"
	"github.com/AndreiMartynenko/common/pkg/db/transaction"
)

func TestInitChannels(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type messagesRepositoryMockFunc func(mc *minimock.Controller) repository.MessagesRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type transactorMockFunc func(mc *minimock.Controller) db.Transactor

	type args struct {
		ctx context.Context
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		repositoryErr = fmt.Errorf("failed to init existing chats")
		ids           = []string{"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"}
	)

	tests := []struct {
		name                   string
		args                   args
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
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.GetChatsMock.Expect(minimock.AnyContext).Return(ids, nil)
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
				return mock
			},
		},
		{
			name: "chat repository error case",
			args: args{
				ctx: ctx,
			},
			err: repositoryErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.GetChatsMock.Expect(minimock.AnyContext).Return(nil, repositoryErr)
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

			err := srv.InitChannels(tt.args.ctx)
			require.Equal(t, tt.err, err)
		})
	}
}
