package dtos

type (
	FilterOrderHistory struct {
		Limit       int    `query:"limit"`
		Page        int    `query:"page"`
		Description string `query:"description"`
	}

	DataOrderHistory struct {
		DtosModel
		Description string `json:"description" validate:"required"`
		OrderItemID string `json:"order_item_id" validate:"required"`
	}
)
