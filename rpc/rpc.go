package rpc

import "github.com/koala/middleware"

func BuildClientMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	if len(mids) == 0 {
		return handle
	}

	m := middleware.Chain(mids[0], mids...)
	return m(handle)
}