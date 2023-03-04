package daos

type (
	FilterUsers struct {
		Limit    int    `query:"limit"`
		Page     int    `query:"page"`
		Fullname string `query:"fullname"`
	}

	UpdateUser struct {
		Nama     string `json:"fullname"`
		Username string `json:"username"`
	}
)
