package domain

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	LoginRequest
	Name          string `json:"name"`
	Bio           string `json:"bio"`
	ImageData     string
	ImageFileName string
}

type SwipeRequest struct {
	SwipedID  uint   `json:"swiped_id"`
	SwipeType string `json:"type"`
}
