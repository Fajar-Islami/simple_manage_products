package dtos

import (
	"encoding/json"
	"time"
)

type (
	FilterOrderItems struct {
		Limit         int    `query:"limit" validate:"omitempty,gt=0"`
		Page          int    `query:"page" validate:"omitempty,gt=0"`
		Name          string `query:"name"`
		PriceMoreThan int    `query:"price_more_than" validate:"omitempty,gt=0"`
		PriceLessThan int    `query:"price_less_than" validate:"omitempty,gt=0"`
		WithExpired   bool   `query:"with_expired" `
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

func (rdoi *ResDataOrderItemsData) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &rdoi); err != nil {
		return err
	}

	return nil
}
