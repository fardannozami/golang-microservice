package main

import (
	"bytes"
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
				"quantity": 11,
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

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				t.Errorf("goroutine %d got unexpected status: %d", id, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
}
