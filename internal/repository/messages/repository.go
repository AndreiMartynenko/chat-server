package messages

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/AndreiMartynenko/chat-server/internal/model"
	"github.com/AndreiMartynenko/chat-server/internal/repository"
	"github.com/AndreiMartynenko/chat-server/internal/repository/messages/converter"
	modelRepo "github.com/AndreiMartynenko/chat-server/internal/repository/messages/model"
	"github.com/AndreiMartynenko/common/pkg/db"
)

const (
	tableName = "messages"

	chatIDColumn    = "chat_id"
	fromColumn      = "from_user"
	textColumn      = "text"
	timestampColumn = "timestamp"
)

type repo struct {
	db db.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(db db.Client) repository.MessagesRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatID string, message *model.Message) error {
	id, err := uuid.Parse(chatID)
	if err != nil {
		return err
	}

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, fromColumn, textColumn, timestampColumn).
		Values(id, message.From, message.Text, message.Timestamp)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "messages_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetMessages(ctx context.Context, chatID string) ([]*model.Message, error) {
	id, err := uuid.Parse(chatID)
	if err != nil {
		return nil, err
	}

	builderSelect := sq.Select(fromColumn, textColumn, timestampColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "messages_repository.GetMessages",
		QueryRaw: query,
	}

	var messages []*modelRepo.Message
	err = r.db.DB().ScanAllContext(ctx, &messages, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToMessagesFromRepo(messages), nil
}

func (r *repo) DeleteChat(ctx context.Context, chatID string) error {
	id, err := uuid.Parse(chatID)
	if err != nil {
		return err
	}

	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "messages_repository.DeleteChat",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
