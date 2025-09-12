package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"testing"
)

func TestRaceConditionOrders(t *testing.T) {
	url := "http://localhost:8080/api/v1/orders"
	payload := []byte(`{
		"user_id": "customer123",
		"items": [
			{
				"product_id": "1",
				"quantity": 2,
				"price": 10.99
			}
		]
	}`)

	var wg sync.WaitGroup
	clients := 20 // jumlah concurrent request

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
			if err != nil {
				t.Errorf("goroutine %d error create request: %v", id, err)
				return
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("goroutine %d error request: %v", id, err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				t.Logf("goroutine %d SUCCESS, response: %s", id, string(body))
			} else {
				t.Errorf("goroutine %d FAILED, status: %d, response: %s", id, resp.StatusCode, string(body))
			}
		}(i)
	}

	wg.Wait()
}
