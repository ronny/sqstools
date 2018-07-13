package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ronny/sqstools/internal/mocks"
	"github.com/ronny/sqstools/internal/sqstools"
)

func TestMoveMessages(t *testing.T) {
	src := &mocks.MessageReceiveDeleter{}
	dest := &mocks.BatchMessageSender{}
	cfg := &config{
		sourceQueue:                    sqstools.NewSQSQueue("-source-"),
		destinationQueue:               sqstools.NewSQSQueue("-dest-"),
		sourceMaxReceiveMessages:       int64(1),
		sourceVisibilityTimeoutSeconds: int64(2),
		sourceWaitTimeSeconds:          int64(3),
		sourceReceiveIterations:        1,
	}

	src.On("ReceiveMessage", &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String("-source-"),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(2),
		WaitTimeSeconds:     aws.Int64(3),
	}).Return(&sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			&sqs.Message{
				Body:          aws.String("-Body-"),
				ReceiptHandle: aws.String("-ReceiptHandle-"),
			},
		},
	}, nil)

	dest.On("SendMessageBatch", &sqs.SendMessageBatchInput{
		QueueUrl: aws.String("-dest-"),
		Entries: []*sqs.SendMessageBatchRequestEntry{
			&sqs.SendMessageBatchRequestEntry{
				Id:          aws.String("0"),
				MessageBody: aws.String("-Body-"),
			},
		},
	}).Return(&sqs.SendMessageBatchOutput{
		Successful: []*sqs.SendMessageBatchResultEntry{
			&sqs.SendMessageBatchResultEntry{},
		},
		Failed: []*sqs.BatchResultErrorEntry{},
	}, nil)

	src.On("DeleteMessageBatch", &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String("-source-"),
		Entries: []*sqs.DeleteMessageBatchRequestEntry{
			&sqs.DeleteMessageBatchRequestEntry{
				Id:            aws.String("0"),
				ReceiptHandle: aws.String("-ReceiptHandle-"),
			},
		},
	}).Return(&sqs.DeleteMessageBatchOutput{
		Successful: []*sqs.DeleteMessageBatchResultEntry{
			&sqs.DeleteMessageBatchResultEntry{},
		},
		Failed: []*sqs.BatchResultErrorEntry{},
	}, nil)

	moveMessages(src, dest, cfg)

	src.AssertExpectations(t)
	dest.AssertExpectations(t)
}
