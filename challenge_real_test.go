package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	b := []byte(`
		{
			"first_name": "primero",
			"last_name": "uno",
			"email": "1"
		}`)
	req, err := http.NewRequest("POST", "/cards", bytes.NewReader(b))
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(createCard)

	handler.ServeHTTP(rr, req)

	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())
}
