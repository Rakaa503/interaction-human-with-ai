package lazarus

import "github.com/gofiber/fiber/v3"

// Middleware membatasi request yang sedang diproses.
// Request yang tidak mendapat slot langsung ditolak tanpa antrean.
func Middleware(limiter *Limiter) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if limiter == nil {
			c.Set("Retry-After", "1")
			c.Status(fiber.StatusServiceUnavailable)
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "lazarus_unavailable",
			})
		}

		release, ok := limiter.TryAcquire()
		if !ok {
			c.Set("Retry-After", "1")
			c.Status(fiber.StatusServiceUnavailable)
			return c.JSON(fiber.Map{
				"success": false,
				"error":   "server_overloaded",
				"message": "Server is temporarily at capacity. Retry shortly.",
			})
		}

		defer release()
		return c.Next()
	}
}
