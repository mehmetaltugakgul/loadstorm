# LoadStorm

LoadStorm is a simple command-line tool for conducting load tests on web services. It allows you to send multiple HTTP requests to a specified URL and analyze the performance of your service under different conditions.

## Features

- Send a specified number of HTTP requests to a target URL.
- Choose the HTTP method for the requests.
- Include optional data payload with each request.
- View detailed information about each request, including response status, body, and headers.
- Log results to a file for further analysis.

## Getting Started

### Prerequisites

Make sure you have Go installed on your machine.

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/mehmetaltugakgul/loadstorm.git
   cd loadstorm

2. 
     ```bash
   go run main.go
   

Usage
Enter the URL to load test when prompted.
Specify the number of requests to send.
Choose the HTTP method to use (e.g., GET, POST).
Optionally provide data to include in the request body.
The load test will run, and the results will be displayed.
Stopping the Load Test:
Press Ctrl+C to stop the load test at any time.

Configuration
You can customize the load test by modifying the source code. Configuration options are available in the loadtest.go and duration.go files.

Results
Results of the load test, including successful and failed requests, will be displayed at the end of the test. Additionally, detailed logs are saved in the request_logs.txt file.

License
This project is licensed under the MIT License - see the LICENSE.md file for details.

Acknowledgments
This tool was created as part of a project to analyze web service performance.
Thanks to the Go community for providing helpful libraries.
Feel free to contribute or report issues. Happy load testing!
