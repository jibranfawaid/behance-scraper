package server

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/playwright-community/playwright-go"
	log "github.com/sirupsen/logrus"
	"optimus/internal/controllers"
	"optimus/internal/errors"
	"optimus/internal/models/responses"
	"optimus/internal/repositories/scraper"
	"optimus/internal/services"
	"os"
	"time"
)

const (
	TimeoutMessage = "request timeout"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func RunServer(ctx context.Context, pw *playwright.Playwright) {
	e := echo.New()
	defer e.Close()

	e.Validator = &CustomValidator{validator: validator.New()}
	e.IPExtractor = echo.ExtractIPDirect()

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:                    middleware.DefaultSkipper,
		ErrorMessage:               TimeoutMessage,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {},
		Timeout:                    10 * time.Minute,
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.WithContext(c.Request().Context()).Error("Error: ", err)

		e := errors.GetError(err.Error())
		response := responses.ErrorResponse{
			BaseResponse: responses.BaseResponse{
				Status:  e.Status,
				Message: e.Message,
			},
		}

		c.JSON(response.Status, response)
	}

	controllers.NewBehanceHandler(e.Group("behance/v1"), services.NewBehanceService(
		scraper.NewBehanceScraper(pw),
	)).MapRoutes()

	go func() {
		log.Fatal(e.Start(":" + os.Getenv("PORT")))
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error when shutting down: " + err.Error())
	}
}
