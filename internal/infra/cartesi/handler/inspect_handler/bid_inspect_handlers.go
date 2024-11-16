package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/order_usecase"
	"github.com/tribeshq/tribes/pkg/router"
)

type OrderInspectHandlers struct {
	OrderRepository entity.OrderRepository
}

func NewOrderInspectHandlers(orderRepository entity.OrderRepository) *OrderInspectHandlers {
	return &OrderInspectHandlers{
		OrderRepository: orderRepository,
	}
}

func (h *OrderInspectHandlers) FindOrderByIdHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findOrderById := order_usecase.NewFindOrderByIdUseCase(h.OrderRepository)
	res, err := findOrderById.Execute(&order_usecase.FindOrderByIdInputDTO{
		Id: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}
	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	env.Report(order)
	return nil
}

func (h *OrderInspectHandlers) FindBisdByCrowdfundingIdHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findOrdersByCrowdfundingId := order_usecase.NewFindOrdersByCrowdfundingIdUseCase(h.OrderRepository)
	res, err := findOrdersByCrowdfundingId.Execute(&order_usecase.FindOrdersByCrowdfundingIdInputDTO{
		CrowdfundingId: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find orders by crowdfunding id: %v", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}

func (h *OrderInspectHandlers) FindAllOrdersHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllOrders := order_usecase.NewFindAllOrdersUseCase(h.OrderRepository)
	res, err := findAllOrders.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}
	allOrders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all orders: %w", err)
	}
	env.Report(allOrders)
	return nil
}
