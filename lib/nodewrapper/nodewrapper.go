package nodewrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type (
	Handler func(*Context, json.RawMessage) (interface{}, error)

	Context struct {
		AwsRequestID             string `json:"awsRequestId"`
		FunctionName             string `json:"functionName"`
		FunctionVersion          string `json:"functionVersion"`
		Invokeid                 string `json:"invokeid"`
		IsDefaultFunctionVersion bool   `json:"isDefaultFunctionVersion"`
		LogGroupName             string `json:"logGroupName"`
		LogStreamName            string `json:"logStreamName"`
		MemoryLimitInMB          string `json:"memoryLimitInMB"`
	}

	Payload struct {
		// custom event fields
		Event json.RawMessage `json:"event"`

		// default context object
		Context *Context `json:"context"`
	}

	Response struct {
		// Request id is an incremental integer
		// representing the request that has been
		// received by this go proc during it's
		// lifetime
		Success   *bool `json:"success,omitempty"`
		RequestId int   `json:"-"` // can retrun request_id for debugging"`
		// Any errors that occur during processing
		// or are returned by handlers are returned
		Error *string `json:"errorMessage,omitempty"`
		Token *string `json:"token,omitempty"`
		// Response output data returned if no errors are present
		Data *interface{} `json:"response,omitempty"`
	}
)

var requestId int // process req id

func NewErrorResponse(err error) *Response {
	e := err.Error()
	errorEscaped := strings.Replace(strings.Replace(e, "\n", "", -1), "\t", " ", -1)
	return &Response{
		RequestId: requestId,
		Error:     &errorEscaped,
	}
}

func NewResponse(data interface{}) *Response {
	success := true

	if token, ok := data.(map[string]interface{})["token"]; ok {

		delete(data.(map[string]interface{}), "token")

		return &Response{
			Success:   &success,
			Token:     token.(*string),
			RequestId: requestId,
			Data:      &data,
		}
	}

	return &Response{
		Success:   &success,
		RequestId: requestId,
		Data:      &data,
	}
}

func Run(handler Handler, rawOutput bool) {
	RunStream(handler, os.Stdin, os.Stdout, rawOutput)
}

func RunStream(handler Handler, Stdin io.Reader, Stdout io.Writer, rawOutput bool) {

	stdin := json.NewDecoder(Stdin)
	stdout := json.NewEncoder(Stdout)
	decodeError := false

	for ; ; requestId++ {
		if err := func() (err error) {
			defer func() {
				if e := recover(); e != nil {
					err = fmt.Errorf("Server Error: 'Panic' %v", e)
				}
			}()
			var payload Payload
			if err := stdin.Decode(&payload); err != nil {
				decodeError = true
				return fmt.Errorf("Bad Request: %v", err)
			}
			data, err := handler(payload.Context, payload.Event)
			if err != nil {
				return err
			}

			if rawOutput {
				// for this particular purpose return as bytestream
				return stdout.Encode(data)
			}
			return stdout.Encode(NewResponse(data))
		}(); err != nil {
			if encErr := stdout.Encode(NewErrorResponse(err)); encErr != nil {
				// bad times
				requestId++
				break
				log.Println("Failed to encode err response!", encErr.Error())
			} else {
				// if invalid JSON, break to advance stdin and restart the loop
				if decodeError {
					requestId++
					break
				}
			}
		}
	}

	RunStream(handler, Stdin, Stdout, rawOutput)
}
