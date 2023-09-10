package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

func run(app *cli.Context) error {
	args := app.Args()
	if args.Len() > 1 {
		return fmt.Errorf("expected only one argument")
	}

	var data string

	stat, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("stdin: %w", err)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var stdin []byte
		stdin, err = io.ReadAll(os.Stdin)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("read stdin: %w", err)
		}
		if len(stdin) == 0 {
			return fmt.Errorf("no data provided")
		}
		data = string(stdin)
	} else {
		data = args.First()
	}

	if data == "" || app.Bool("tui") {
		return runTUI(app, data)
	} else {
		return runCLI(app, data)
	}
}
