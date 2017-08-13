package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"os"

	"aws-lambda-barcode-generator/config"
	"aws-lambda-barcode-generator/lib/nodewrapper"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/pdf417"
	"github.com/boombuler/barcode/qr"
)

type LambdaInput struct {
	Width   *int    `json:"width,omitempty"`
	Height  *int    `json:"height,omitempty"`
	Message *string `json:"message,omitempty"`
	Type    *string `json:"type,omitempty"`
}

func main() {
	nodewrapper.Run(func(context *nodewrapper.Context, eventJSON json.RawMessage) (interface{}, error) {

		var input LambdaInput
		var c barcode.Barcode
		var err error

		if err = json.Unmarshal(eventJSON, &input); err != nil {
			// Add error prefixes to all final error strings to allow for AWS API Gateway regex
			return nil, fmt.Errorf("400 Bad Request: Invalid Request - cannot marshal JSON input. Object received: %s", string(eventJSON))
		}

		err = checkValidRequest(input)

		if err != nil {
			return nil, err
		}

		switch barcodeType := *input.Type; barcodeType {
		case "qr":
			c, err = qr.Encode(*input.Message, qr.M, qr.Unicode)
		case "pdf417":
			c, err = pdf417.Encode(*input.Message, 3)
		}

		if err != nil {
			return nil, fmt.Errorf("500 Server Error: Could not generate barcode. Details: %s", err)
		}

		// Scale the barcode to 200x200 pixels
		c, err = barcode.Scale(c, *input.Width, *input.Height)

		if err != nil {
			return nil, fmt.Errorf("500 Server Error: Could not scale the barcode to input width and height. Details: %s", err)
		}

		// if we have develop set to true, then also export as PNG (so can view on local machine)
		if config.DEVELOP {
			// create the output file
			file, _ := os.Create("code.png")
			defer file.Close()

			// encode the barcode as png
			png.Encode(file, c)
		}

		// write image to buffer
		buf := new(bytes.Buffer)
		err = png.Encode(buf, c)

		if err != nil {
			return nil, fmt.Errorf("500 Server Error: Could not write image to buffer. Details: %s", err)
		}

		// return the image as a base64 encoded string
		return base64.StdEncoding.EncodeToString(buf.Bytes()), err
	}, true)
}

func checkValidRequest(input LambdaInput) error {
	// check if input conforms with our spec
	if input.Height == nil {
		return fmt.Errorf("400 Bad Request: 'height' field is missing.")
	}

	if input.Width == nil {
		return fmt.Errorf("400 Bad Request: 'width' field is missing.")
	}

	if input.Message == nil {
		return fmt.Errorf("400 Bad Request: 'message' field is missing.")
	}

	if input.Type == nil {
		return fmt.Errorf("400 Bad Request: 'type' field is missing.")
	}

	h := *input.Height
	w := *input.Width
	m := *input.Message
	t := *input.Type

	if h < 150 {
		return fmt.Errorf("400 Bad Request: 'height' needs to be a minimum of 150.")
	}

	if h > 3000 {
		return fmt.Errorf("400 Bad Request: 'height' can be a maximum of 3000.")
	}

	if w < 150 {
		return fmt.Errorf("400 Bad Request: 'width' needs to be a minimum of 150.")
	}

	if w > 3000 {
		return fmt.Errorf("400 Bad Request: 'width' can be a maximum of 3000.")
	}

	if len(m) > 600 {
		return fmt.Errorf("400 Bad Request: Your message is too large, it can be a maximum of 600 bytes.")
	}

	if t != "qr" && t != "pdf417" {
		return fmt.Errorf("400 Bad Request: Invalid barcode type. Type needs to be: 'qr', 'pdf417'.")
	}

	return nil
}
