package order_items_usecase

import (
	"context"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
)

type OrderItemsUseCase interface {
	GetAllOrderItems(ctx context.Context, params daos.FilterOrderItems) (res []daos.OrderItems, err *helper.ErrorStruct)
	GetOrderItemsByID(ctx context.Context, orderItemsid int) (res daos.OrderItems, err *helper.ErrorStruct)
	CreateOrderItems(ctx context.Context, data daos.OrderItems) (res uint, err *helper.ErrorStruct)
	UpdateOrderItemsByID(ctx context.Context, orderItemsid int, data daos.OrderItems) (res string, err *helper.ErrorStruct)
	DeleteOrderItemsByID(ctx context.Context, orderItemsid int) (res string, err *helper.ErrorStruct)
}
