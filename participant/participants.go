package participant


type Participant struct {
	Id           int    `json:"id"`
	Name         string `json:"name" binding:"required"`
}

