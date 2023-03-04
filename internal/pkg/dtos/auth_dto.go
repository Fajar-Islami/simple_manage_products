package dtos

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterRequest struct {
		Nama     string `json:"fullname" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	LoginResp struct {
		DtosModel
		Nama     string `json:"fullname"`
		Username string `json:"username"`
	}
)
