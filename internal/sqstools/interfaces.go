package sqstools

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageReceiveDeleter interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessageBatch(*sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error)
}

type BatchMessageSender interface {
	SendMessageBatch(*sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error)
}
