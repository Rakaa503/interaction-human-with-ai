package lazarus

import "github.com/gofiber/fiber/v3"

// Middleware mempertahankan kompatibilitas dengan pemanggilan lama.
func Middleware(limiter *Limiter) func(fiber.Ctx) error {
	return MiddlewareWithTelemetry(limiter, nil)
}

// MiddlewareWithTelemetry mencatat request nyata tanpa membuat antrean baru.
func MiddlewareWithTelemetry(limiter *Limiter, telemetry *Telemetry) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if limiter == nil {
			if telemetry != nil {
				telemetry.RecordRejected()
			}
			c.Set("Retry-After", "1")
			c.Status(fiber.StatusServiceUnavailable)
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "lazarus_unavailable",
			})
		}

		release, ok := limiter.TryAcquire()
		if !ok {
			if telemetry != nil {
				telemetry.RecordRejected()
			}
			c.Set("Retry-After", "1")
			c.Status(fiber.StatusServiceUnavailable)
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "server_overloaded",
				"message": "Server is temporarily at capacity. Retry shortly.",
			})
		}

		if telemetry != nil {
			finishTelemetry := telemetry.StartRequest()
			defer finishTelemetry()
		}

		defer release()
		return c.Next()
	}
}
