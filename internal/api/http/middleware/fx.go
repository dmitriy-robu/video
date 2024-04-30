package middleware

import (
	"go.uber.org/fx"
)

type Middleware struct {
	ClientAuthMiddleware *ClientAuthMiddleware
	AdminAuthMiddleware  *AdminAuthMiddleware
}

func NewMiddlewares(
	clientAuth *ClientAuthMiddleware,
	adminAuth *AdminAuthMiddleware,
) *Middleware {
	return &Middleware{
		ClientAuthMiddleware: clientAuth,
		AdminAuthMiddleware:  adminAuth,
	}
}

func NewMiddleware() fx.Option {
	return fx.Module(
		"middleware",
		fx.Provide(
			NewClientAuthMiddleware,
			NewAdminAuthMiddleware,
			NewMiddlewares,
		),
	)
}
