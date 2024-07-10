package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatAPI "github.com/AndreiMartynenko/chat-server/internal/api/chat"
	"github.com/AndreiMartynenko/chat-server/internal/model"
	"github.com/AndreiMartynenko/chat-server/internal/service"
	serviceMocks "github.com/AndreiMartynenko/chat-server/internal/service/mocks"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

func TestSend(t *testing.T) {
	t.Parallel()

	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

		from      = "from"
		text      = "text"
		timestamp = timestamppb.Now()

		message = &model.Message{
			From:      from,
			Text:      text,
			Timestamp: timestamp.AsTime(),
		}

		serviceErr = fmt.Errorf("service error")

		req = &desc.SendMessageRequest{
			Message: &desc.Message{
				From:      from,
				Text:      text,
				Timestamp: timestamp,
			},
			ChatId: chatID,
		}

		res = &empty.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, chatID, message).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, chatID, message).Return(serviceErr)
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

			res, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
