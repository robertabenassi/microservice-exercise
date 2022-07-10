package http_handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"testing"

	"microservice-exercise/internal/data_model"
	"microservice-exercise/internal/faker"
	"microservice-exercise/internal/grpc"
	"microservice-exercise/internal/stream"
)

type getPortTest struct {
	key          string
	expectedPort data_model.Port
}

func TestHandleLoadPorts(t *testing.T) {
	client := NewRESTClient(
		&grpc.Client{
			Client: faker.NewFakePortServiceClient(), // mock of the
		},
		stream.NewStream(),
	)
	req := httptest.NewRequest("POST", "http://example.com/updatePorts", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(client.HandleUpdatePorts)
	handler.ServeHTTP(w, req)
	if w.Code == http.StatusInternalServerError {
		t.Fatal("Got %i while we expected 200", w.Code)
	}
}
