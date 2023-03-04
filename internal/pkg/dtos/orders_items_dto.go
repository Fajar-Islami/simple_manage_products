package dtos

import (
	"time"
)

type (
	FilterOrderItems struct {
		Limit         int    `query:"limit"`
		Page          int    `query:"page"`
		Name          string `query:"name"`
		PriceMoreThan int    `query:"price_more_than" validate:"omitempty,gt=0"`
		PriceLessThan int    `query:"price_less_than" validate:"omitempty,gt=0"`
	}

	ReqDataOrderItems struct {
		Name      string `json:"name" validate:"required"`
		Price     int    `json:"price" validate:"required,gt=0"`
		ExpiredAt string `json:"expired_at" validate:"required"`
	}

	ReqDataUpdateOrderItems struct {
		Name      string `json:"name,omitempty"`
		Price     int    `json:"price,omitempty" validate:"omitempty,gt=0"`
		ExpiredAt string `json:"expired_at,omitempty"`
	}

	ResDataOrderItems struct {
		Data []ResDataOrderItemsData `json:"data"`
		Pagination
	}

	ResDataOrderItemsData struct {
		DtosModel
		Name      string    `json:"name"`
		Price     int       `json:"price"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)
