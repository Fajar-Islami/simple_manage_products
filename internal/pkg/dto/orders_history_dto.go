package daos

type (
	FilterOrderHistory struct {
		Limit       int    `query:"limit"`
		Page        int    `query:"page"`
		Description string `query:"description"`
	}

	DataOrderHistory struct {
		Description string `json:"description" validate:"required"`
		OrderItemID string `json:"order_item_id" validate:"required"`
	}
)
