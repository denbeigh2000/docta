package checkers

import (
	"testing"

	"github.com/denbeigh2000/docta"

	"fmt"
	"net/http"
	"net/http/httptest"
)

type simpleHTTPHandler struct {
	f http.HandlerFunc
}

func (h simpleHTTPHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.f(rw, r)
}

func buildHandler(code int, responseText string) http.Handler {
	return simpleHTTPHandler{func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(code)
		rw.Write([]byte(responseText))
	}}
}

func TestHTTPCheckerSimple(t *testing.T) {
	handler := buildHandler(200, "OK")
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := HTTPChecker{"yellow", "red", RequestData{
		Endpoint: server.URL,
		Method:   "GET",
		Data:     nil,
	}}

	state := checker.Check()
	if state.State != docta.Green {
		t.Errorf("Expected state %v, got %v", docta.Green, state.State)
	}

	if state.Info != docta.DefaultInfo {
		t.Errorf("Expected info %v, got %v", docta.DefaultInfo, state.Info)
	}
}

func TestHTTPCheckerNon200(t *testing.T) {
	handler := buildHandler(500, "Not OK")
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := HTTPChecker{"yellow", "red", RequestData{
		Endpoint: server.URL,
		Method:   "GET",
		Data:     nil,
	}}

	state := checker.Check()
	expectedState := docta.Red
	if state.State != expectedState {
		t.Errorf("Expected state %v, got %v", expectedState, state.State)
	}

	expectedText := "Non-200 status code: 500"
	if state.Info != expectedText {
		t.Errorf("Expected info %v, got %v", expectedText, state.Info)
	}
}

func TestHTTPCheckerYellowText(t *testing.T) {
	forbiddenText := "YELLAHTEXT"
	handler := buildHandler(
		200, fmt.Sprintf("Some other data goes before %v and also after", forbiddenText),
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := HTTPChecker{forbiddenText, "red", RequestData{
		Endpoint: server.URL,
		Method:   "GET",
		Data:     nil,
	}}

	state := checker.Check()

	expectedState := docta.Yellow
	if state.State != expectedState {
		t.Errorf("Expected state %v, got %v", expectedState, state.State)
	}

	expectedText := fmt.Sprintf("Response body contains forbidden string %v", forbiddenText)
	if state.Info != expectedText {
		t.Errorf("Expected info %v, got %v", expectedText, state.Info)
	}
}
