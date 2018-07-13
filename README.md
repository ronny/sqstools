# sqstools

Tools for performing common SQS tasks.

## Installation

Get the latest binary for your platform from the [releases page], rename it and run it with `-h` to see usage. This is
the easiest option if you don't have the go toolkit installed.

Alternatively you can get the source and build it yourself:

```
go get -u github.com/ronny/sqstools/.../
$GOBIN/sqs-move-messages -h
```

## `sqs-move-messages`

Moves messages from one SQS queue to another destination (currently only supports another SQS queue). This can be used to
requeue messages from a dead letter queue back to the main queue, for example.

[releases page]: https://github.com/ronny/sqstools/releases
