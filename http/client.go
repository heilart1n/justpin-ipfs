package http

import (
	"github.com/ybbus/httpretry"
	"net/http"
	"time"
)

func NewClient(client *http.Client) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	return httpretry.NewCustomClient(
		client,
		// retry 5 times
		httpretry.WithMaxRetryCount(5),
		// retry on status == 429, if status >= 500, if err != nil, or if response was nil (status == 0)
		httpretry.WithRetryPolicy(func(statusCode int, err error) bool {
			return err != nil || statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError || statusCode == 0
		}),
		// every retry should wait one more 10 second
		// backoff will be: 5, 10, 20, 40, 80, ...
		httpretry.WithBackoffPolicy(httpretry.ExponentialBackoff(5*time.Second, time.Minute, time.Second)),
	)
}
