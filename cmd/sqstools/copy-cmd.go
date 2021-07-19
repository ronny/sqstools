package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/ronny/sqstools"
	"github.com/ronny/sqstools/destinations"
)

func CreateCopyCommand() *ffcli.Command {
	fs := flag.NewFlagSet("sqstools copy", flag.ExitOnError)
	var (
		srcURL                   = fs.String("src", "", "the URL of the source SQS queue")
		destURL                  = fs.String("dest", "", "the URL of the destination, supported: file://, s3://, https:// (SQS only), and stdout:// (for debugging)")
		maxRcv                   = fs.Int("maxRcv", 10, "max messages to receive every time we call ReceiveMessage on the src SQS queue")
		visibilityTimeoutSeconds = fs.Int("vis", 10, "visibility timeout in seconds for every ReceiveMessage call")
		waitTimeSeconds          = fs.Int("wait", 0, "wait time in seconds for every ReceiveMessage call")
		region                   = fs.String("region", "us-east-1", "the AWS region")
		delete                   = fs.Bool("delete", false, "if true, deletes messages after they are successfully copied to the destination")
		destFormat               = fs.String("destFormat", "json", "format to use to serialise messages when not sending to an SQS queue, only `json` and `prettyjson` are supported")
		maxIters                 = fs.Int("maxIters", 0, "maximum number of ReceiveMessage iterations (0 = no maximum, repeat until empty)")
	)

	return &ffcli.Command{
		Name:       "copy",
		ShortUsage: "sqstools copy -src <src-queue-url> -dest <dest-url> [other flags]",
		ShortHelp:  "Copies messages from an SQS queue to one of the supported destinations",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if srcURL == nil || *srcURL == "" {
				fs.Usage()
				return errors.New("missing src queue URL")
			}

			loadAwsConfCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			awsCfg, err := awsConfig.LoadDefaultConfig(loadAwsConfCtx, awsConfig.WithRegion(*region))
			if err != nil {
				return fmt.Errorf("aws.LoadDefaultConfig: %w", err)
			}
			cancel()

			src, err := sqstools.NewSQSQueueSource(sqstools.NewSQSQueueSourceInput{
				QueueURL:                 *srcURL,
				MaxReceiveMessages:       int32(*maxRcv),
				VisibilityTimeoutSeconds: int32(*visibilityTimeoutSeconds),
				WaitTimeSeconds:          int32(*waitTimeSeconds),
				AWSConfig:                awsCfg,
			})
			if err != nil {
				return fmt.Errorf("NewSQSQueueSource: %w", err)
			}

			var dest sqstools.Destination
			{
				jsonMarshaler := destinations.NewJSONMarshaler(destinations.WithPrettyJSON(*destFormat == "prettyjson"))
				if strings.HasPrefix(*destURL, "stdout://") {
					dest = &destinations.StdoutDestination{
						Marshaler: jsonMarshaler,
					}
				} else if strings.HasPrefix(*destURL, "file://") {
					dest = &destinations.FileDestination{
						PathPrefix: strings.Replace(*destURL, "file://", "", 1),
						PathSuffix: ".json",
						Perm:       0644, // TODO
						Marshaler:  jsonMarshaler,
					}
				} else if strings.HasPrefix(*destURL, "s3://") {
					return errors.New("s3:// not yet implemented")
				} else if strings.HasPrefix(*destURL, "https://sqs.") {
					return errors.New("SQS destination not yet implemented")
				} else {
					return fmt.Errorf("unsupported destination: %s", *destURL)
				}
			}

			copy := sqstools.NewCopy(src, dest, *delete)

			log.Printf("src=%s, dest=%s\n", *srcURL, *destURL)

			err = copy.CopyRepeatedly(ctx, *maxIters)
			if err != nil {
				return fmt.Errorf("CopyRepeatedly: %w", err)
			}
			return nil
		},
	}
}
