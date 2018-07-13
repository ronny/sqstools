package main

import (
	"strconv"

	"github.com/ronny/sqstools/internal/sqstools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Errorf("Invalid arg: %s", err)
		usageAndExit()
	}

	log.Infof("Source: %s", cfg.sourceQueue.URL())
	log.Infof(" Dest.: %s", cfg.destinationQueue.URL())
	log.Infof("Moving %d message(s) at a time for %d iterations", cfg.sourceMaxReceiveMessages, cfg.sourceReceiveIterations)

	srcSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.sourceQueue.Region()),
	}))
	srcSvc := sqs.New(srcSess)

	destSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(cfg.destinationQueue.Region()),
	}))
	destSvc := sqs.New(destSess)

	moveMessages(srcSvc, destSvc, cfg)
}

func moveMessages(src sqstools.MessageReceiveDeleter, dest sqstools.BatchMessageSender, cfg *config) {
	for i := 1; i <= cfg.sourceReceiveIterations; i++ {
		iterLogger := log.WithFields(log.Fields{"i": i, "msgsPerRcv": cfg.sourceMaxReceiveMessages, "iters": cfg.sourceReceiveIterations})

		result, err := src.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(cfg.sourceQueue.URL()),
			MaxNumberOfMessages: &cfg.sourceMaxReceiveMessages,
			VisibilityTimeout:   &cfg.sourceVisibilityTimeoutSeconds,
			WaitTimeSeconds:     &cfg.sourceWaitTimeSeconds,
		})

		if err != nil {
			iterLogger.Fatalf("Receiving messages from source: %v", err)
		}

		if len(result.Messages) == 0 {
			iterLogger.Infof("No messages in source queue")
			return
		}

		var sendMessageBatchRequestEntries []*sqs.SendMessageBatchRequestEntry
		var deleteMessageBatchRequestEntries []*sqs.DeleteMessageBatchRequestEntry

		for index, message := range result.Messages {
			sendMessageBatchRequestEntries = append(
				sendMessageBatchRequestEntries,
				&sqs.SendMessageBatchRequestEntry{
					Id:          aws.String(strconv.Itoa(index)),
					MessageBody: message.Body,
				},
			)

			deleteMessageBatchRequestEntries = append(
				deleteMessageBatchRequestEntries,
				&sqs.DeleteMessageBatchRequestEntry{
					Id:            aws.String(strconv.Itoa(index)),
					ReceiptHandle: message.ReceiptHandle,
				},
			)
		}

		sendOutput, err := dest.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  sendMessageBatchRequestEntries,
			QueueUrl: aws.String(cfg.destinationQueue.URL()),
		})

		if err != nil {
			iterLogger.Fatalf("Sending messages to destination: %v", err)
		}

		iterLogger.Infof("Sent to destination: %d successful, %d failed", len(sendOutput.Successful), len(sendOutput.Failed))
		iterLogger.Debugf("Send output: %v", sendOutput)

		resultDelete, err := src.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			QueueUrl: aws.String(cfg.sourceQueue.URL()),
			Entries:  deleteMessageBatchRequestEntries,
		})

		if err != nil {
			iterLogger.Fatalf("Deleting message: %v", err)
		}

		iterLogger.Infof("Deleted %d successfully, %d failed", len(resultDelete.Successful), len(resultDelete.Failed))
		iterLogger.Debugf("Delete batch output: %v", resultDelete)
	}
}
