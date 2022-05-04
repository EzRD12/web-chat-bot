package tests

import (
	"bot/core"
	"bot/test/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func initHttpMock() {
	core.Client = &mocks.MockClient{}
}

func initHttpClient() {
	core.Client = &http.Client{}
}

// TestStockQuote calls GetStockQuote with a code, checking
// for a valid return value.
func TestStockCode(t *testing.T) {
	initHttpClient()
	code := "aapl.us"
	msg, err := core.GetStockQuote(code)
	if msg == "Could not get stock quote." || err != nil {
		t.Fatalf(`GetStockQuote("aapl.us") = %q, %v, want "AAPL.US quote is 148.69 per share.", nil`, msg, err)
	}
}

// TestStockEmpty calls GetStockQuote with a code, checking
// for a valid return value.
func TestStockEmpty(t *testing.T) {
	initHttpClient()
	code := ""
	msg, err := core.GetStockQuote(code)
	if msg != "" || err == nil {
		t.Fatalf(`GetStockQuote("") = %q, %v, want "Could not get stock quote.", invalid code`, msg, err)
	}
}

// TestStockWith500Response calls GetStockQuote with an empty code, checking
// for an invalid return message.
func TestStockWith500Response(t *testing.T) {
	initHttpMock()
	code := "aapl.us"
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       nil,
		}, errors.New("Internal server error")
	}
	msg, err := core.GetStockQuote(code)
	assert.NotNil(t, err)
	assert.Empty(t, msg)
	assert.Equal(t, err.Error(), "error parsing CSV from URL")
}

// TestReadUrlWith500Response calls ReadCSVFromUrl with a url, checking
// for a invalid response with a 500 status code.
func TestReadUrlWith500Response(t *testing.T) {
	initHttpMock()
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       nil,
		}, errors.New("Internal server error")
	}
	msg, err := core.ReadCSVFromUrl("")
	assert.NotNil(t, err)
	assert.Empty(t, msg)
	assert.Contains(t, err.Error(), "Internal server error")
}

// TestReadUrlWithBadRequestResponse calls ReadCSVFromUrl with a url, checking
// for a invalid response with a 400 status code.
func TestReadUrlWithBadRequestResponse(t *testing.T) {
	initHttpMock()
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       nil,
		}, nil
	}
	msg, err := core.ReadCSVFromUrl("")
	assert.NotNil(t, err)
	assert.Empty(t, msg)
	assert.Contains(t, err.Error(), "invalid response status code")
}

// TestStockNonExisting calls GetStockQuote with a non existing code, checking
// for an expected message.
func TestStockNonExisting(t *testing.T) {
	initHttpClient()
	code := "HELLO WORLD"
	msg, err := core.GetStockQuote(code)
	if msg == "" || err != nil {
		t.Fatalf(`GetStockQuote("HELLO WORLD") = %q, %v, want "Could not get stock quote.", nil`, msg, err)
	}
}
