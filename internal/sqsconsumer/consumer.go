package sqsconsumer

import (
	"context"
	"discord-diplomacy/internal/types"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/bwmarrin/discordgo"
)

type MessagePoster interface {
	PostMessage(channelID string, message *discordgo.MessageSend) error
}

type Consumer struct {
	client   *sqs.Client
	logger   *slog.Logger
	poster   MessagePoster
	queueURL string
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
	var notification types.NotificationMessage
	if err := json.Unmarshal([]byte(aws.ToString(message.Body)), &notification); err != nil {
		return fmt.Errorf("decode message body: %w", err)
	}

	if err := notification.Validate(); err != nil {
		return err
	}

	msg := &discordgo.MessageSend{
		Content: notification.Content,
	}
	if err := c.poster.PostMessage(notification.ChannelID, msg); err != nil {
		return err
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
