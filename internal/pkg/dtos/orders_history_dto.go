package dtos

import "encoding/json"

type (
	FilterOrderHistory struct {
		Limit       int    `query:"limit" validate:"omitempty,gt=0"`
		Page        int    `query:"page" validate:"omitempty,gt=0"`
		Description string `query:"description"`
	}

	ReqCreateDataOrderHistoryItem struct {
		Description string `json:"description" validate:"required"`
		OrderItemID int    `json:"order_item_id" validate:"required"`
	}

	ReqUpdateDataOrderHistoryItem struct {
		Description string `json:"description,omitempty"`
		OrderItemID int    `json:"order_item_id,omitempty"`
	}

	ResDataOrderHistoryItem struct {
		DtosModel
		Description string                `json:"description"`
		OrderItem   ResDataOrderItemsData `json:"order_item"`
	}

	ResDataOrderHistory struct {
		Data []ResDataOrderHistoryItem `json:"data"`
		Pagination
	}
)

func (rdoi *ResDataOrderHistoryItem) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &rdoi); err != nil {
		return err
	}

	return nil
}
