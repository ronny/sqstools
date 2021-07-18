package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	rootFlagSet := flag.NewFlagSet("sqstools", flag.ExitOnError)

	root := &ffcli.Command{
		ShortUsage:  "sqstools <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{CreateCopyCommand()},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	err := root.ParseAndRun(context.Background(), os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
}
