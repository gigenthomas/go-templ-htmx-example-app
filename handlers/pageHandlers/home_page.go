package pageHandlers

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
	"github.com/alekLukanen/go-templ-htmx-example-app/ui/pages"
)

type HomePageHandler struct {
	ctx    context.Context
	logger *slog.Logger
	*services.ServiceMesh
}

func NewHomePageHandler(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh) *HomePageHandler {
	return &HomePageHandler{
		ctx:         ctx,
		logger:      logger,
		ServiceMesh: serviceMesh,
	}
}

func (obj *HomePageHandler) RegisterPublicRoutes(echoHandler *echo.Echo) {
	echoHandler.GET("/", obj.BasePage)
}

func (obj *HomePageHandler) BasePage(echoCtx echo.Context) error {
	isLoggedIn := false
	if cookie, err := echoCtx.Cookie("token"); err == nil {
		if _, err := obj.UserAuthenticationService.ValidateToken(echoCtx.Request().Context(), cookie.Value); err == nil {
			isLoggedIn = true
		}
	}
	//isLoggedIn = true
	component := pages.Homepage(isLoggedIn)

	handlers.Render(echoCtx, &component)
	return nil
}
