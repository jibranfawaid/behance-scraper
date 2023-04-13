package responses

import "optimus/internal/models"

type SearchResponse struct {
	BaseResponse
	Data []*models.SearchModel `json:"data"`
}
