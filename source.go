package sqstools

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSQueueSource struct {
	queueURL                 string
	maxReceiveMessages       int32
	visibilityTimeoutSeconds int32
	waitTimeSeconds          int32
	client                   *sqs.Client
}

type NewSQSQueueSourceInput struct {
	QueueURL                 string
	MaxReceiveMessages       int32
	VisibilityTimeoutSeconds int32
	WaitTimeSeconds          int32
	AWSConfig                aws.Config
}

func NewSQSQueueSource(input NewSQSQueueSourceInput) (*SQSQueueSource, error) {
	if input.QueueURL == "" {
		return nil, errors.New("missing QueueURL")
	}
	return &SQSQueueSource{
		queueURL:                 input.QueueURL,
		maxReceiveMessages:       input.MaxReceiveMessages,
		visibilityTimeoutSeconds: input.VisibilityTimeoutSeconds,
		waitTimeSeconds:          input.WaitTimeSeconds,
		client:                   sqs.NewFromConfig(input.AWSConfig),
	}, nil
}

func (src *SQSQueueSource) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	return src.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(src.queueURL),
		MaxNumberOfMessages: src.maxReceiveMessages,
		VisibilityTimeout:   src.visibilityTimeoutSeconds,
		WaitTimeSeconds:     src.waitTimeSeconds,
	})
}

func (src *SQSQueueSource) DeleteMessageBatch(ctx context.Context, messages []types.Message) (*sqs.DeleteMessageBatchOutput, error) {
	entriesToDelete := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		id := strconv.Itoa(i)
		entriesToDelete[i] = types.DeleteMessageBatchRequestEntry{
			Id:            &id,
			ReceiptHandle: message.ReceiptHandle,
		}
	}

	return src.client.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(src.queueURL),
		Entries:  entriesToDelete,
	})
}
