package main

import (
	"fmt"

	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

func runCLI(app *cli.Context, data string) error {
	code, err := qrCodeFromData(data, 0, qrcode.RecoveryLevel(app.Int("recovery-level")))
	if err != nil {
		return err
	}

	var qrCodeData string
	if app.Bool("big") {
		qrCodeData = code.ToString(app.Bool("inverted"))
	} else {
		qrCodeData = code.ToSmallString(app.Bool("inverted"))
	}

	fmt.Println(qrCodeData)
	return nil
}
