package qrcode

import (
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func GenerateQRBytes(address string) ([]byte, error) {
	// New QR Code
	qrc, err := qrcode.NewWith(
		address,
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	)
	if err != nil {
		return nil, err
	}

	// Transparent png
	options := []standard.ImageOption{
		standard.WithBgTransparent(),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	}

	// Write QR Code to bytes
	qrWriter := NewWriteCloser()
	w := standard.NewWithWriter(qrWriter, options...)

	// save file
	if err = qrc.Save(w); err != nil {
		return nil, err
	}

	return qrWriter.buf, err
}
