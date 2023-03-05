package dtos

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterRequest struct {
		Fullname string `json:"fullname" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	LoginResp struct {
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
)
