package requests

type SearchRequest struct {
	Search string `query:"search" validate:"required"`
}
