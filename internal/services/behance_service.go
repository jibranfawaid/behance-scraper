package services

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"net/http"
	e "optimus/internal/errors"
	"optimus/internal/models/responses"
	"optimus/internal/repositories/scraper"
)

type BehanceService interface {
	Search(ctx context.Context, search string) (*responses.SearchResponse, error)
}

type behanceService struct {
	behanceScraper scraper.BehanceScraper
}

func NewBehanceService(behanceScraper scraper.BehanceScraper) *behanceService {
	return &behanceService{
		behanceScraper: behanceScraper,
	}
}

func (s *behanceService) Search(ctx context.Context, search string) (searchResponse *responses.SearchResponse, err error) {
	ctx, sp := otel.Tracer("").Start(ctx, "Search")
	defer sp.End()

	result, err := s.behanceScraper.Search(ctx, search)
	if err != nil {
		log.WithContext(ctx).Error("Unable to retrieve search result: " + err.Error())
		return nil, errors.New(e.GeneralError)
	}

	searchResponse = &responses.SearchResponse{
		BaseResponse: responses.BaseResponse{
			Status:  http.StatusOK,
			Message: "Successfully fetched search result",
		},
		Data: result,
	}

	return searchResponse, err
}
