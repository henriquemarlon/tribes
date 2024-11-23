package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

func (h *OrderInspectHandlers) FindOrderByIdHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findOrderById := order_usecase.NewFindOrderByIdUseCase(h.OrderRepository)
	res, err := findOrderById.Execute(ctx, &order_usecase.FindOrderByIdInputDTO{
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

func (h *OrderInspectHandlers) FindBisdByCrowdfundingIdHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findOrdersByCrowdfundingId := order_usecase.NewFindOrdersByCrowdfundingIdUseCase(h.OrderRepository)
	res, err := findOrdersByCrowdfundingId.Execute(ctx, &order_usecase.FindOrdersByCrowdfundingIdInputDTO{
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

func (h *OrderInspectHandlers) FindAllOrdersHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findAllOrders := order_usecase.NewFindAllOrdersUseCase(h.OrderRepository)
	res, err := findAllOrders.Execute(ctx)
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

func (h *OrderInspectHandlers) FindOrdersByInvestorHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	address := strings.ToLower(router.PathValue(ctx, "address"))
	findOrdersByInvestor := order_usecase.NewFindOrdersByInvestorUseCase(h.OrderRepository)
	res, err := findOrdersByInvestor.Execute(ctx, &order_usecase.FinsOrdersByInvestorInputDTO{
		Investor: common.HexToAddress(address),
	})
	if err != nil {
		return fmt.Errorf("failed to find orders by investor: %w", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}
