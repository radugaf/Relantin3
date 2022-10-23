package main

import (
	"net/http"

	"github.com/radugaf/RelantinV3/logger/models"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (config *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = config.readJSON(w, r, &requestPayload)

	// insert data
	event := models.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := config.Models.LogEntry.Insert(event)
	if err != nil {
		config.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	config.writeJSON(w, http.StatusAccepted, resp)
}
