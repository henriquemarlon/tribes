package repository

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type OrderRepositorySqlite struct {
	Db *gorm.DB
}

func NewOrderRepositorySqlite(db *gorm.DB) *OrderRepositorySqlite {
	return &OrderRepositorySqlite{
		Db: db,
	}
}

func (r *OrderRepositorySqlite) CreateOrder(ctx context.Context, input *entity.Order) (*entity.Order, error) {
	err := r.Db.Raw(`
		INSERT INTO orders (crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`, input.CrowdfundingId, input.Investor.String(), input.Amount.Hex(), input.InterestRate.Hex(), input.State, input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return input, nil
}

func (r *OrderRepositorySqlite) FindOrderById(ctx context.Context, id uint) (*entity.Order, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE id = ? LIMIT 1", id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}

	return &entity.Order{
		Id:             uint(result["id"].(int64)),
		CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
		Investor:       common.HexToAddress(result["investor"].(string)),
		Amount:         uint256.MustFromHex(result["amount"].(string)),
		InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
		State:          entity.OrderState(result["state"].(string)),
		CreatedAt:      result["created_at"].(int64),
		UpdatedAt:      result["updated_at"].(int64),
	}, nil
}

func (r *OrderRepositorySqlite) FindOrdersByCrowdfundingId(ctx context.Context, id uint) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", id).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by crowdfunding ID: %w", err)
	}

	var orders []*entity.Order
	for _, data := range results {
		orders = append(orders, &entity.Order{
			Id:             uint(data["id"].(int64)),
			CrowdfundingId: uint(data["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(data["investor"].(string)),
			Amount:         uint256.MustFromHex(data["amount"].(string)),
			InterestRate:   uint256.MustFromHex(data["interest_rate"].(string)),
			State:          entity.OrderState(data["state"].(string)),
			CreatedAt:      data["created_at"].(int64),
			UpdatedAt:      data["updated_at"].(int64),
		})
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByState(ctx context.Context, crowdfundingId uint, state string) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ? AND state = ?", crowdfundingId, state).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}

	var orders []*entity.Order
	for _, data := range results {
		orders = append(orders, &entity.Order{
			Id:             uint(data["id"].(int64)),
			CrowdfundingId: uint(data["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(data["investor"].(string)),
			Amount:         uint256.MustFromHex(data["amount"].(string)),
			InterestRate:   uint256.MustFromHex(data["interest_rate"].(string)),
			State:          entity.OrderState(data["state"].(string)),
			CreatedAt:      data["created_at"].(int64),
			UpdatedAt:      data["updated_at"].(int64),
		})
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByInvestor(ctx context.Context, investor common.Address) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE investor = ?", investor.String()).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by investor: %w", err)
	}

	var orders []*entity.Order
	for _, data := range results {
		orders = append(orders, &entity.Order{
			Id:             uint(data["id"].(int64)),
			CrowdfundingId: uint(data["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(data["investor"].(string)),
			Amount:         uint256.MustFromHex(data["amount"].(string)),
			InterestRate:   uint256.MustFromHex(data["interest_rate"].(string)),
			State:          entity.OrderState(data["state"].(string)),
			CreatedAt:      data["created_at"].(int64),
			UpdatedAt:      data["updated_at"].(int64),
		})
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindAllOrders(ctx context.Context) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}

	var orders []*entity.Order
	for _, data := range results {
		orders = append(orders, &entity.Order{
			Id:             uint(data["id"].(int64)),
			CrowdfundingId: uint(data["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(data["investor"].(string)),
			Amount:         uint256.MustFromHex(data["amount"].(string)),
			InterestRate:   uint256.MustFromHex(data["interest_rate"].(string)),
			State:          entity.OrderState(data["state"].(string)),
			CreatedAt:      data["created_at"].(int64),
			UpdatedAt:      data["updated_at"].(int64),
		})
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) UpdateOrder(ctx context.Context, input *entity.Order) (*entity.Order, error) {
	order, err := r.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}

	if input.Amount != nil && input.Amount.Sign() > 0 {
		order.Amount = input.Amount
	}
	if input.InterestRate != nil && input.InterestRate.Sign() > 0 {
		order.InterestRate = input.InterestRate
	}
	if input.State != "" {
		order.State = input.State
	}
	order.UpdatedAt = input.UpdatedAt

	res := r.Db.Model(&entity.Order{}).Where("id = ?", input.Id).Updates(map[string]interface{}{
		"amount":        order.Amount.Hex(),
		"interest_rate": order.InterestRate.Hex(),
		"state":         order.State,
		"updated_at":    order.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update order: %w", res.Error)
	}
	return order, nil
}

func (r *OrderRepositorySqlite) DeleteOrder(ctx context.Context, id uint) error {
	res := r.Db.Delete(&entity.Order{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete order: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}
