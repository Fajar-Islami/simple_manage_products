package dtos

type (
	FilterUsers struct {
		Limit    int    `query:"limit" validate:"omitempty,gt=0"`
		Page     int    `query:"page" validate:"omitempty,gt=0"`
		Fullname string `query:"fullname"`
	}

	UpdateUser struct {
		Fullname string `json:"fullname,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	ResGetUserData struct {
		Fullname string `json:"fullname"`
		Username string `json:"username"`
	}

	ResGetUser struct {
		Data []ResGetUserData `json:"data"`
		Pagination
	}
)
