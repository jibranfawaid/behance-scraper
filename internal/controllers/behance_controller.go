package controllers

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"optimus/internal/middleware"
	"optimus/internal/models/requests"
	"optimus/internal/services"
)

var behanceTracer = otel.Tracer("behanceController")

type behanceHandler struct {
	group *echo.Group

	behanceService services.BehanceService
}

func NewBehanceHandler(group *echo.Group, behanceService services.BehanceService) *behanceHandler {
	return &behanceHandler{
		group:          group,
		behanceService: behanceService,
	}
}

func (h *behanceHandler) MapRoutes() {
	h.group.GET("/search", h.Search(), middleware.Recovery())
}

func (h *behanceHandler) Search() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := behanceTracer.Start(c.Request().Context(), "Search")
		defer span.End(trace.WithStackTrace(true))

		log.WithContext(ctx).Info("Search API started")
		*c.Request() = *(c.Request().WithContext(ctx))

		var searchRequest requests.SearchRequest
		if err := c.Bind(&searchRequest); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
		}

		searchResponse, err := h.behanceService.Search(ctx, searchRequest.Search)
		if err != nil {
			return err
		}

		log.WithContext(ctx).Info("Search API has finished")
		return c.JSON(http.StatusOK, searchResponse)
	}
}
