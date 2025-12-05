package user

type UpdateUserProfileRequest struct {
	Name   string `binding:"omitempty,max=50" json:"name"`
	Avatar string `binding:"omitempty,url"    json:"avatar"`
	UserID uint   `binding:"required"         json:"user_id"`
}

type UpdateUserProfileResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
