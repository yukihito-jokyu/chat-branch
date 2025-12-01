package queue

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
)

// 指定したトピックにメッセージを送信する
func NewPublisher(db *sql.DB, logger *slog.Logger) (message.Publisher, error) {
	publisher, err := watermillSQL.NewPublisher(
		db,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultMySQLSchema{},
		},
		watermill.NewSlogLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	return publisher, nil
}

// 指定したトピックからメッセージを受信する
func NewSubscriber(db *sql.DB, logger *slog.Logger) (message.Subscriber, error) {
	subscriber, err := watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:    watermillSQL.DefaultMySQLSchema{},
			OffsetsAdapter:   watermillSQL.DefaultMySQLOffsetsAdapter{},
			PollInterval:     1 * time.Second,
			InitializeSchema: true,
		},
		watermill.NewSlogLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}

// 指定したトピックにメッセージを送信する
func PublishTask(publisher message.Publisher, topic string, payload []byte) error {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	return publisher.Publish(topic, msg)
}
