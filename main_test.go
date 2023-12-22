package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set up test data
	url := server.URL
	method := "GET"
	data := []byte("testData")

	var wg sync.WaitGroup
	var mu sync.Mutex
	result := LoadTestResult{}

	// Run the makeRequest function in a goroutine
	wg.Add(1)
	go makeRequest(url, method, data, 1, &wg, &mu, &result)

	// Wait for the goroutine to finish
	wg.Wait()

	// Check the result
	if result.TotalRequests != 1 {
		t.Errorf("Expected total requests to be 1, got %d", result.TotalRequests)
	}
	if result.SuccessfulRequests != 1 {
		t.Errorf("Expected successful requests to be 1, got %d", result.SuccessfulRequests)
	}
	if result.FailedRequests != 0 {
		t.Errorf("Expected failed requests to be 0, got %d", result.FailedRequests)
	}
}

func TestRunLoadTestWithRate(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set up test data
	url := server.URL
	method := "GET"
	data := []byte("testData")
	config := LoadTestConfig{
		URL:         url,
		NumRequests: 5,
		Method:      method,
		Data:        data,
	}
	duration := time.Millisecond * 100

	// Run the load test
	result := runLoadTestWithRate(config, duration)

	// Check the result
	if result.TotalRequests != 5 {
		t.Errorf("Expected total requests to be 5, got %d", result.TotalRequests)
	}
	if result.SuccessfulRequests != 5 {
		t.Errorf("Expected successful requests to be 5, got %d", result.SuccessfulRequests)
	}
	if result.FailedRequests != 0 {
		t.Errorf("Expected failed requests to be 0, got %d", result.FailedRequests)
	}
}
