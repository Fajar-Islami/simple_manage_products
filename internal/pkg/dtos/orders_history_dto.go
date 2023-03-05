package dtos

type (
	FilterOrderHistory struct {
		Limit       int    `query:"limit" validate:"omitempty,gt=0"`
		Page        int    `query:"page" validate:"omitempty,gt=0"`
		Description string `query:"description"`
	}

	DataOrderHistory struct {
		DtosModel
		Description string `json:"description" validate:"required"`
		OrderItemID string `json:"order_item_id" validate:"required"`
	}
)
