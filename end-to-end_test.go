package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	time.Sleep(time.Second)

	resp, err := http.Get("http://localhost:8080/")
	assert.NoError(t, err, "Error making request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be OK")
}
