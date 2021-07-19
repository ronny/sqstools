package sqstools

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Source interface {
	ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error)
	DeleteMessageBatch(ctx context.Context, messages []types.Message) (*sqs.DeleteMessageBatchOutput, error)
}

type Destination interface {
	HandleMessages(ctx context.Context, messages []types.Message) error
}

type Copy struct {
	src    Source
	dest   Destination
	delete bool
}

func NewCopy(src Source, dest Destination, delete bool) *Copy {
	return &Copy{
		src:    src,
		dest:   dest,
		delete: delete,
	}
}

func (c *Copy) CopyRepeatedly(ctx context.Context, maxIterations int) error {
	i := 0
	for {
		handled, err := c.CopySingleBatch(ctx)
		if err != nil {
			return fmt.Errorf("CopySingleBatch: %w", err)
		}
		i = i + 1
		if handled == 0 || (maxIterations > 1 && i >= maxIterations) {
			return nil
		}
	}
}

func (c *Copy) CopySingleBatch(ctx context.Context) (int, error) {
	handled := 0

	receiveOutput, err := c.src.ReceiveMessage(ctx)
	if err != nil {
		return handled, fmt.Errorf("ReceiveMessage: %w", err)
	}

	if len(receiveOutput.Messages) == 0 {
		return handled, nil
	}

	err = c.dest.HandleMessages(ctx, receiveOutput.Messages)
	if err != nil {
		return handled, fmt.Errorf("HandleMessages: %w", err)
	}
	handled = len(receiveOutput.Messages)

	if c.delete {
		deleteOutput, err := c.src.DeleteMessageBatch(ctx, receiveOutput.Messages)
		if err != nil {
			return handled, fmt.Errorf("DeleteMessageBatch: %w", err)
		}
		if len(deleteOutput.Failed) != 0 {
			return handled, fmt.Errorf("DeleteMessageBatch: failed to delete some/all messages: %v", deleteOutput.Failed)
		}
	}
	return handled, nil
}
