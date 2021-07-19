package destinations

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type FileDestination struct {
	PathPrefix string
	PathSuffix string
	Perm       fs.FileMode
	Marshaler  Marshaler
}

func (d *FileDestination) HandleMessages(ctx context.Context, messages []types.Message) error {
	if d.Marshaler == nil {
		return errors.New("FileDestination missing Marshaler")
	}

	if d.PathPrefix == "" {
		return errors.New("FileDestination missing PathPrefix")
	}

	if d.PathSuffix == "" {
		return errors.New("FileDestination missing PathSuffix")
	}

	if d.Perm == 0 {
		d.Perm = 0644
	}

	for _, message := range messages {
		bytes, err := d.Marshaler.Marshal(message)
		if err != nil {
			return fmt.Errorf("Marshal: %w", err)
		}

		name := path.Join(d.PathPrefix, *message.MessageId+d.PathSuffix)
		err = os.WriteFile(name, bytes, d.Perm)
		if err != nil {
			return fmt.Errorf("os.WriteFile: %w", err)
		}
	}

	return nil
}
