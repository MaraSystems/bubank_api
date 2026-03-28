package utils

type IDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type PageRequest struct {
	Limit  *int32 `form:"limit" binding:"omitempty,min=1,max=20"`
	Offset *int32 `form:"offset" binding:"omitempty,min=0"`
}
