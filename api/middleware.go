// api/middleware.go

package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
)

func (api *API) KeyAccess() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// TODO
		// - get the key from the header, or bail out 401 Unauthorized
		// - look it up, or bail out 401 Unauthorized
		// -	do not keep failed lookups, could be vector, e.g. massive hdr
		// - now add to log because we have it
		// - check access, we can hit this url or not? not bail 403 Forbidden
		// - stash in c, how?
		// - then auth for e.g. agent use in new-agent handler
		slogfiber.AddCustomAttributes(c, slog.String("key", "WHATEVER"))
		return c.Next()
	}

}
