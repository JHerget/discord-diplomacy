package sqsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	messageTypeDiscordMessage = "discord_message"
	maxDiscordContentLength   = 2000
)

type MessagePoster interface {
	PostMessage(channelID string, content string) error
}

type Consumer struct {
	client   *sqs.Client
	logger   *slog.Logger
	poster   MessagePoster
	queueURL string
}

type NotificationMessage struct {
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

func New(client *sqs.Client, logger *slog.Logger, poster MessagePoster, queueURL string) *Consumer {
	return &Consumer{
		client:   client,
		logger:   logger,
		poster:   poster,
		queueURL: queueURL,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	c.logger.Info("starting SQS consumer")

	for {
		if err := ctx.Err(); err != nil {
			c.logger.Info("stopping SQS consumer")
			return nil
		}

		output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
			VisibilityTimeout:   60,
		})
		if err != nil {
			if ctx.Err() != nil {
				c.logger.Info("stopping SQS consumer")
				return nil
			}
			c.logger.Error("receive SQS messages", "error", err)
			wait(ctx, 2*time.Second)
			continue
		}

		for _, message := range output.Messages {
			c.handleMessage(ctx, message)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, message sqstypes.Message) {
	if err := c.processMessage(message); err != nil {
		c.logger.Error("process SQS message", "message_id", aws.ToString(message.MessageId), "error", err)
		return
	}

	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		c.logger.Error("delete SQS message", "message_id", aws.ToString(message.MessageId), "error", err)
		return
	}

	c.logger.Info("processed SQS message", "message_id", aws.ToString(message.MessageId))
}

func (c *Consumer) processMessage(message sqstypes.Message) error {
	var notification NotificationMessage
	if err := json.Unmarshal([]byte(aws.ToString(message.Body)), &notification); err != nil {
		return fmt.Errorf("decode message body: %w", err)
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	if err := c.poster.PostMessage(notification.ChannelID, notification.Content); err != nil {
		return err
	}

	return nil
}

func (m NotificationMessage) Validate() error {
	if m.Type != messageTypeDiscordMessage {
		return fmt.Errorf("unsupported message type %q", m.Type)
	}
	if strings.TrimSpace(m.ChannelID) == "" {
		return fmt.Errorf("channel_id is required")
	}
	if strings.TrimSpace(m.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if utf8.RuneCountInString(m.Content) > maxDiscordContentLength {
		return fmt.Errorf("content exceeds %d characters", maxDiscordContentLength)
	}

	return nil
}

func wait(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
