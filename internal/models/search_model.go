package models

type SearchModel struct {
	ImageUrl   string `json:"imageUrl"`
	ProjectUrl string `json:"projectUrl"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	TotalLikes int    `json:"totalLikes"`
	TotalViews int    `json:"totalViews"`
}
