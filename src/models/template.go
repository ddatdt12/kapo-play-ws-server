package models

type Template struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	IsPublic    bool   `json:"isPublic"`
	CreatorID   uint   `json:"creatorId"`
	// Creator     User `json:"creator"`
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
}
