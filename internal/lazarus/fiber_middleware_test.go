package lazarus

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
)

func TestMiddlewareRejectsWhenCapacityIsFull(t *testing.T) {
	limiter, err := NewLimiter(1)
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	app.Use(Middleware(limiter))

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})

	app.Get("/slow", func(c fiber.Ctx) error {
		close(requestStarted)
		<-releaseRequest
		return c.SendStatus(http.StatusOK)
	})

	firstRequestDone := make(chan error, 1)

	go func() {
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		resp, err := app.Test(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				firstRequestDone <- fiber.NewError(resp.StatusCode)
				return
			}
		}
		firstRequestDone <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(3 * time.Second):
		close(releaseRequest)
		t.Fatal("first request did not start")
	}

	secondReq := httptest.NewRequest(http.MethodGet, "/slow", nil)
	secondResp, err := app.Test(secondReq)
	if err != nil {
		close(releaseRequest)
		t.Fatalf("second request failed: %v", err)
	}
	defer secondResp.Body.Close()

	if secondResp.StatusCode != http.StatusServiceUnavailable {
		close(releaseRequest)
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			secondResp.StatusCode,
		)
	}

	if secondResp.Header.Get("Retry-After") != "1" {
		close(releaseRequest)
		t.Fatal("expected Retry-After header to be 1")
	}

	close(releaseRequest)

	select {
	case err := <-firstRequestDone:
		if err != nil {
			t.Fatalf("first request failed: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("first request did not finish")
	}

	stats := limiter.Stats()
	if stats.Active != 0 {
		t.Fatalf("expected 0 active requests, got %d", stats.Active)
	}
}
