// Testing requires the node package lambda-test
// Download from github link:

package main

import (
	"aws-lambda-barcode-generator/lib/osutil"
	"encoding/json"
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

type lambdaInputTest struct {
	Width   interface{} `json:"width,omitempty"`
	Height  interface{} `json:"height,omitempty"`
	Message *string     `json:"message,omitempty"`
	Type    *string     `json:"type,omitempty"`
}

type lambdaOutput struct {
	Response *interface{} `json:"response,omitempty"`
	Error    *string      `json:"errorMessage,omitempty"`
}

// Tests will be called on the compiled binary so ensure we are running the latest version
func TestBuild(t *testing.T) {
	assert := assert.New(t)

	_, err := osutil.Run("", "go build")
	assert.Nil(err, fmt.Sprintf("Build failed: %v", err))
}

func TestMissingHeight(t *testing.T) {
	width := 200
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'height' field is missing.")
}

func TestMissingWidth(t *testing.T) {
	height := 200
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'width' field is missing.")
}

func TestMissingMessage(t *testing.T) {
	height := 200
	width := 200
	ctype := "qr"

	l := lambdaInputTest{
		Height: &height,
		Width:  &width,
		Type:   &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'message' field is missing.")
}

func TestMissingType(t *testing.T) {
	height := 200
	width := 200
	message := "Test"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'type' field is missing.")
}

func TestInvalidType(t *testing.T) {
	height := 200
	width := 200
	message := "Test"
	ctype := "test"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Type:    &ctype,
		Message: &message,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: Invalid barcode type. Type needs to be: 'qr', 'pdf417'.")
}

func TestInvalidHeight(t *testing.T) {
	height := "testnotallowed"
	width := 200
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Contains(*resp.Error, "400 Bad Request: Invalid Request - cannot marshal JSON input.")
}

func TestInvalidWidth(t *testing.T) {
	height := 200
	width := "testnotallowed"
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Contains(*resp.Error, "400 Bad Request: Invalid Request - cannot marshal JSON input.")
}

func TestMinimumHeight(t *testing.T) {
	height := 149
	width := 200
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'height' needs to be a minimum of 150.")
}

func TestMaximumHeight(t *testing.T) {
	height := 3001
	width := 200
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'height' can be a maximum of 3000.")
}

func TestMinimumWidth(t *testing.T) {
	height := 200
	width := 149
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'width' needs to be a minimum of 150.")
}

func TestMaximumWidth(t *testing.T) {
	height := 200
	width := 3001
	message := "Test"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: 'width' can be a maximum of 3000.")
}

func TestMaximalMessageLength(t *testing.T) {
	height := 200
	width := 200
	message := ""
	ctype := "qr"

	for i := 0; i <= 100; i++ {
		message += "abcdefgh"
	}

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: Your message is too large, it can be a maximum of 600 bytes.")
}

func TestSuccesfulBarcode(t *testing.T) {
	height := 200
	width := 600
	message := "測試這個"
	ctype := "qr"

	l := lambdaInputTest{
		Height:  &height,
		Width:   &width,
		Message: &message,
		Type:    &ctype,
	}

	runLambda(l, 10)

	/*resp, err := runLambda(l, 10)

	assert := assert.New(t)
	assert.Nil(err, fmt.Sprintf("lambda-test returned an unexpected error: %v", err))
	assert.Equal(*resp.Error, "400 Bad Request: Your message is too large, it can be a maximum of 600 bytes.")*/

}

// Helper method, runs lambda-test and passes in the JSON
func runLambda(request lambdaInputTest, timeout int) (*lambdaOutput, error) {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := osutil.Run("", fmt.Sprintf("lambda-test -e %v -t %v", string(requestJSON), timeout))
	responseObject := &lambdaOutput{}

	if err != nil {
		return nil, err
	} else {
		json.Unmarshal(response, responseObject)
	}

	return responseObject, nil
}
