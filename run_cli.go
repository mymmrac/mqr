package main

import (
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

func runCLI(app *cli.Context, data string) error {
	code, err := qrCodeFromData(data, 0, qrcode.RecoveryLevel(app.Int("recovery-level")))
	if err != nil {
		return err
	}

	var output string
	if app.Bool("big") {
		output = code.ToString(app.Bool("inverted"))
	} else {
		output = code.ToSmallString(app.Bool("inverted"))
	}

	fmt.Println("\n " + strings.ReplaceAll(output, "\n", "\n "))
	return nil
}
