package main

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func qrCodeFromData(data string, version int, recoveryLevel qrcode.RecoveryLevel) (*qrcode.QRCode, error) {
	var err error
	var code *qrcode.QRCode

	if version == 0 {
		code, err = qrcode.New(data, recoveryLevel)
	} else {
		code, err = qrcode.NewWithForcedVersion(data, version, recoveryLevel)
	}
	if err != nil {
		return nil, fmt.Errorf("generate QR code: %w", err)
	}

	return code, err
}
