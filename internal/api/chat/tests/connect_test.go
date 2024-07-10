package tests

/*import (
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	chatAPI "github.com/AndreiMartynenko/chat-server/internal/api/chat"
	"github.com/AndreiMartynenko/chat-server/internal/model"
	"github.com/AndreiMartynenko/chat-server/internal/service"
	serviceMocks "github.com/AndreiMartynenko/chat-server/internal/service/mocks"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		req    *desc.ConnectRequest
		stream desc.ChatV1_ConnectServer
	}

	var (
		mc = minimock.NewController(t)

		chatID   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		username = "username"

		serviceErr = fmt.Errorf("service error")

		stream        desc.ChatV1_ConnectServer
		serviceStream model.Stream

		req = &desc.ConnectRequest{
			ChatId:   chatID,
			Username: username,
		}
	)

	tests := []struct {
		name            string
		args            args
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				req:    req,
				stream: stream,
			},
			err: nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.ConnectMock.Expect(chatID, username, serviceStream).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				req:    req,
				stream: stream,
			},
			err: serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.ConnectMock.Expect(chatID, username, serviceStream).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chatAPI.NewImplementation(chatServiceMock)

			err := api.Connect(tt.args.req, tt.args.stream)
			require.Equal(t, tt.err, err)
		})
	}
}*/
