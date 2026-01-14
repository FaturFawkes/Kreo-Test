package request

// GetCategoryOverviewRequest represents the request parameters for getting category overview.
type GetCategoryOverviewRequest struct {
	Category string `uri:"category" binding:"required,min=2,max=50"`
}
