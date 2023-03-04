package daos

type (
	FilterOrderItems struct {
		Limit         int    `query:"limit"`
		Page          int    `query:"page"`
		Name          string `query:"name"`
		PriceMoreThan int    `query:"price_more_than" validate:"gt=0;gtefield=PriceLessThan"`
		PriceLessThan int    `query:"price_less_than" validate:"gt=0"`
	}

	DataOrderItems struct {
		Description string `json:"description" validate:"required"`
		OrderItemID string `json:"order_item_id" validate:"required"`
	}
)
