package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   bool        `json:"error"`
}

func (app *Config) readJSON(writer http.ResponseWriter, request *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB

	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	// Why the second call to decode?
	// Because { "foo": "bar" } is correct.
	// But { "foo": "bar" } { "baz": "qux" } is not.
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}

func (app *Config) writeJSON(writer http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// https://gobyexample.com/variadic-functions
	// Go does not support optional parameters.
	// So the use of a variadic lets you "emulate" that, since you're allowed to have a zero variadic parameter.
	if len(headers) > 0 {
		for k, v := range headers[0] {
			writer.Header().Set(k, v[0])
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusInternalServerError

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(writer, statusCode, payload)
}
