package request

// ListMarketsRequest represents the request parameters for listing markets.
type ListMarketsRequest struct {
	Category string `uri:"category" binding:"required"`
	Page     int    `form:"page" binding:"min=0"`
	Limit    int    `form:"limit" binding:"min=0,max=100"`
	Status   string `form:"status" binding:"omitempty,oneof=open closed settled"`
}
