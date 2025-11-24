package entities

type CreateReviewRequest struct {
	CustomerRequestID string `json:"customer_request_id" validate:"required"`
	Rating            int    `json:"rating" validate:"min=1,max=5"`
	Comment           string `json:"comment" validate:"max=200"`
}

type SkipReviewRequest struct {
	CustomerRequestID string `json:"customer_request_id" validate:"required"`
}

type ReviewResponse struct {
	ReviewID          string `json:"review_id"`
	CustomerRequestID string `json:"customer_request_id"`
	ShopID            string `json:"shop_id"`
	UserID            string `json:"user_id"`
	Rating            int    `json:"rating"`
	Comment           string `json:"comment"`
	IsSkipped         bool   `json:"is_skipped"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type GetReviewsResponse struct {
	Reviews    []ReviewResponse `json:"reviews"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
