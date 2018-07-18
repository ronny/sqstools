# sqstools

Tools for performing common SQS tasks.

## Installation

Get the latest binary for your platform from the [releases page], rename it and run it with `-h` to see usage. This is
the easiest option if you don't have the go toolkit installed.

[releases page]: https://github.com/ronny/sqstools/releases

Alternatively you can get the source and build it yourself:

```
go get -u github.com/ronny/sqstools/.../
$GOBIN/sqs-move-messages -h
```

## `sqs-move-messages`

Moves messages from one SQS queue to another destination (currently only supports another SQS queue). This can be used to
requeue messages from a dead letter queue back to the main queue, for example.

Example:

Move 20 messages (10 at a time) from a DLQ to the main queue:

```
env AWS_PROFILE=foo $GOBIN/sqs-move-messages \
  -src https://sqs.ap-southeast-2.amazonaws.com/1234567890/MyDeadLetterQueue \
  -dest https://sqs.ap-southeast-2.amazonaws.com/1234567890/MyMainQueue \
  -srcMaxMsgsPerRcv 10 \
  -srcIters 2
```

Full usage:

```
Usage of sqs-move-messages:
  -dest string
        the URL of the destination SQS queue
  -logLevel string
        log level (panic, fatal, error, warning, info, debug) (default "info")
  -src string
        the URL of the source SQS queue
  -srcIters int
        how many iterations of ReceiveMessage (default 1)
  -srcMaxMsgsPerRcv int
        the maximum number of messages to receive at a time from the source queue (max 10) (default 1)
  -srcVisTimeout int
        visibility timeout (in seconds) of each message when receiving from the source queue (default 30)
  -srcWaitTime int
        how many seconds the SQS server should wait for each ReceiveMessage calls before returning with messages (max 20) (default 1)
```
