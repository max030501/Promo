package prize

type Prize struct {
	Id           int    `json:"id"`
	Description  string `json:"description" binding:"required"`
}
