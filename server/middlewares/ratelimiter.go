package middlewares

import (
	"fmt"
	"net/http"

	"pfserver/config"
	"pfserver/core"
	"pfserver/utils"

	"github.com/wassimbj/gorl"
)

func RateLimit(f http.HandlerFunc, opts gorl.RLOpts) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		rdsClient, err := config.Redis().Client()
		if err != nil {
			core.Respond(res, core.ResOpts{
				Status: http.StatusInternalServerError,
				Msg:    "Something is wrong !!",
			})
			return
		}

		result, _ := gorl.RateLimiter(req.Context(), gorl.RLOpts{
			Attempts:    opts.Attempts,
			Prefix:      opts.Prefix,
			Duration:    opts.Duration,
			Id:          utils.GetUserIp(req),
			RedisClient: rdsClient,
		})

		// block the user
		if result.Block {
			core.Respond(res, core.ResOpts{
				Status: http.StatusTooManyRequests,
				Msg:    fmt.Sprintf("You have reached the limit, come back after %d seconds", result.TimeLeft/1000),
			})
			return
		}

		// allow the user
		f(res, req)

	}
}
