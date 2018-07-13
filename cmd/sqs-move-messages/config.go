package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ronny/sqstools/internal/sqstools"
	"github.com/sirupsen/logrus"
)

const (
	defaultSourceMaxReceiveMessages       int64 = 1
	defaultSourceVisibilityTimeoutSeconds int64 = 30
	defaultSourceWaitTimeSeconds          int64 = 1
	defaultSourceReceiveIterations        int   = 1

	defaultLogLevel string = "info"
)

type config struct {
	sourceQueue                    sqstools.SQSQueue
	destinationQueue               sqstools.SQSQueue
	sourceMaxReceiveMessages       int64
	sourceVisibilityTimeoutSeconds int64
	sourceWaitTimeSeconds          int64
	sourceReceiveIterations        int

	logLevel logrus.Level
}

func newConfig() (*config, error) {
	cfg := config{}

	var sourceQueueURL string
	var destinationQueueURL string
	var logLevel string

	flag.StringVar(&sourceQueueURL, "src", "", "the URL of the source SQS queue")
	flag.StringVar(&destinationQueueURL, "dest", "", "the URL of the destination SQS queue")

	flag.Int64Var(&cfg.sourceMaxReceiveMessages, "srcMaxMsgsPerRcv", defaultSourceMaxReceiveMessages, "the maximum number of messages to receive at a time from the source queue (max 10)")
	flag.Int64Var(&cfg.sourceVisibilityTimeoutSeconds, "srcVisTimeout", defaultSourceVisibilityTimeoutSeconds, "visibility timeout (in seconds) of each message when receiving from the source queue")
	flag.Int64Var(&cfg.sourceWaitTimeSeconds, "srcWaitTime", defaultSourceWaitTimeSeconds, "how many seconds the SQS server should wait for each ReceiveMessage calls before returning with messages (max 20)")
	flag.IntVar(&cfg.sourceReceiveIterations, "srcIters", defaultSourceReceiveIterations, "how many iterations of ReceiveMessage")

	flag.StringVar(&logLevel, "logLevel", defaultLogLevel, fmt.Sprintf("log level (%s)", logLevels()))

	flag.Parse()

	if sourceQueueURL == "" || destinationQueueURL == "" {
		usageAndExit()
	}

	cfg.sourceQueue = sqstools.NewSQSQueue(sourceQueueURL)
	cfg.destinationQueue = sqstools.NewSQSQueue(destinationQueueURL)

	var err error
	cfg.logLevel, err = logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func logLevels() string {
	var s []string

	for _, l := range logrus.AllLevels {
		s = append(s, l.String())
	}

	return strings.Join(s, ", ")
}

func usageAndExit() {
	flag.Usage()
	os.Exit(1)
}
