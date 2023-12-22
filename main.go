package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

type LoadTestConfig struct {
	URL         string
	NumRequests int
	Method      string
	Data        []byte
}

type LoadTestResult struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	Duration           time.Duration
}

func makeRequest(url string, method string, data []byte, requestNum int, wg *sync.WaitGroup, mu *sync.Mutex, result *LoadTestResult) {
	defer wg.Done()

	startTime := time.Now()
	resp, err := http.NewRequest(method, url, nil)
	if err != nil {
		mu.Lock()
		result.FailedRequests++
		mu.Unlock()
		fmt.Printf(color.RedString("Request %d failed: %v\n"), requestNum, err)
		return
	}

	if data != nil {
		resp.Body = http.NoBody
		resp.ContentLength = 0
		resp.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(data)), nil
		}
	}

	client := http.Client{}
	response, err := client.Do(resp)
	elapsed := time.Since(startTime)

	if err != nil {
		mu.Lock()
		result.FailedRequests++
		mu.Unlock()
		fmt.Printf(color.RedString("Request %d failed: %v\n"), requestNum, err)
		return
	}

	defer response.Body.Close()

	mu.Lock()
	result.SuccessfulRequests++
	result.TotalRequests++
	mu.Unlock()

	fmt.Printf(color.GreenString("Request %d completed in %v\n"), requestNum, elapsed)
	fmt.Printf(color.GreenString("Response Status: %v\n"), response.Status)
	fmt.Printf(color.GreenString("Response Body: %v\n"), response.Body)
	fmt.Printf(color.GreenString("Response Headers: %v\n"), response.Header)
	fmt.Println(color.GreenString("--------------------------------------------------"))
	successMsg := fmt.Sprintf("Request %d completed in %v\n", requestNum, elapsed)
	logAndPrint(color.GreenString, successMsg+fmt.Sprintf("Response Status: %v\n", response.Status)+fmt.Sprintf("Response Body: %v\n", response.Body)+fmt.Sprintf("Response Headers: %v\n", response.Header)+"--------------------------------------------------")

}

func logAndPrint(colorFunc func(string, ...interface{}) string, msg string) {
	fmt.Print(colorFunc(msg))
	log.Println(msg)
	file, err := os.OpenFile("request_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(msg + "\n"); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

func getDuration() time.Duration {
	var durationMillis int
	fmt.Print(color.YellowString("Please enter the duration in milliseconds (0 for infinite): "))
	_, err := fmt.Scanln(&durationMillis)
	if err != nil {
		fmt.Println(color.RedString("Invalid duration"))
		return 0
	}

	return time.Duration(durationMillis) * time.Millisecond
}

func runLoadTestWithRate(config LoadTestConfig, duration time.Duration) LoadTestResult {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var result LoadTestResult

	startTime := time.Now()
	ticker := time.NewTicker(duration)

	defer ticker.Stop()

	for i := 1; i <= config.NumRequests; i++ {
		wg.Add(1)
		go makeRequest(config.URL, config.Method, config.Data, i, &wg, &mutex, &result)

		if duration > 0 {
			<-ticker.C
		}
	}

	wg.Wait()
	result.Duration = time.Since(startTime)

	return result
}

func main() {
	fmt.Println("    __                       __   _____   __                            ")
	fmt.Println("   / /   ____   ____ _  ____/ /  / ___/  / /_  ____    _____   ____ ___ ")
	fmt.Println("  / /   / __ \\ / __ `/ / __  /   \\__ \\  / __/ / __ \\  / ___/  / __ `__ \\")
	fmt.Println(" / /___/ /_/ // /_/ / / /_/ /   ___/ / / /_  / /_/ / / /     / / / / / /")
	fmt.Println("/_____/\\____/ \\__,_/  \\__,_/   /____/  \\__/  \\____/ /_/     /_/ /_/ /_/ ")
	fmt.Println(" ")
	fmt.Println(" ")

	var url string
	fmt.Print(color.YellowString("Please enter the URL to load test: "))
	fmt.Scanln(&url)

	var numRequests int
	fmt.Print(color.YellowString("Please enter the number of requests to send: "))
	fmt.Scanln(&numRequests)

	var method string
	fmt.Print(color.YellowString("Please enter the HTTP method to use(POST/GET/PUT/DELETE): "))
	fmt.Scanln(&method)

	var data []byte
	fmt.Print(color.YellowString("Do you want to send data? (y/n): "))
	var sendData string
	fmt.Scanln(&sendData)
	if sendData == "y" || sendData == "Y" {
		fmt.Print(color.YellowString("Please enter the data to send: "))
		fmt.Scanln(&data)
	}

	loadTestConfig := LoadTestConfig{
		URL:         url,
		NumRequests: numRequests,
		Method:      method,
		Data:        data,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println(color.RedString("Load Test Stopped"))
		os.Exit(0)
	}()

	duration := getDuration()
	if duration < 0 {
		fmt.Println(color.RedString("Invalid duration"))
		return
	}

	result := runLoadTestWithRate(loadTestConfig, duration)

	logFile, err := os.OpenFile("request_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	fmt.Printf(color.GreenString("Load Test Completed in %v\n"), result.Duration)
	fmt.Printf("Total Requests: %d\n", result.TotalRequests)
	fmt.Printf(color.GreenString("Successful Requests: %d\n"), result.SuccessfulRequests)
	fmt.Printf(color.RedString("Failed Requests: %d\n"), result.FailedRequests)
	fmt.Print("Press ENTER to exit...")
	fmt.Scanln()
	fmt.Println("Exiting.")
}
