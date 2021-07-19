package destinations

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type StdoutDestination struct {
	Marshaler Marshaler
}

func (d *StdoutDestination) HandleMessages(ctx context.Context, messages []types.Message) error {
	if d.Marshaler == nil {
		return errors.New("StdoutDestination missing Marshaler")
	}

	for _, message := range messages {
		bytes, err := d.Marshaler.Marshal(message)
		if err != nil {
			return fmt.Errorf("Marshal: %w", err)
		}
		fmt.Fprintln(os.Stdout, string(bytes))
	}
	return nil
}
