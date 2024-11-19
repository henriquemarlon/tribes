package repository

import (
	"encoding/json"
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

func (r *OrderRepositorySqlite) CreateOrder(input *entity.Order) (*entity.Order, error) {
	err := r.Db.Model(&entity.Order{}).Create(map[string]interface{}{
		"crowdfunding_id": input.CrowdfundingId,
		"investor":        input.Investor.String(),
		"amount":          input.Amount.String(),
		"interest_rate":   input.InterestRate.String(),
		"state":           input.State,
		"created_at":      input.CreatedAt,
		"updated_at":      input.UpdatedAt,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return r.FindOrderById(input.Id)
}

func (r *OrderRepositorySqlite) FindOrderById(id uint) (*entity.Order, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE id = ? LIMIT 1", id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}

	order := &entity.Order{
		Id:             uint(result["id"].(int64)),
		CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
		Investor:       common.HexToAddress(result["investor"].(string)),
		Amount:         uint256.MustFromHex(result["amount"].(string)),
		InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
		State:          entity.OrderState(result["state"].(string)),
		CreatedAt:      result["created_at"].(int64),
		UpdatedAt:      result["updated_at"].(int64),
	}

	return order, nil
}

func (r *OrderRepositorySqlite) FindOrdersByCrowdfundingId(id uint) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", id).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by crowdfunding ID: %w", err)
	}

	var orders []*entity.Order
	for _, result := range results {
		order := &entity.Order{
			Id:             uint(result["id"].(int64)),
			CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(result["investor"].(string)),
			Amount:         uint256.MustFromHex(result["amount"].(string)),
			InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
			State:          entity.OrderState(result["state"].(string)),
			CreatedAt:      result["created_at"].(int64),
			UpdatedAt:      result["updated_at"].(int64),
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByState(crowdfundingId uint, state string) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ? AND state = ?", crowdfundingId, state).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}

	var orders []*entity.Order
	for _, result := range results {
		order := &entity.Order{
			Id:             uint(result["id"].(int64)),
			CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(result["investor"].(string)),
			Amount:         uint256.MustFromHex(result["amount"].(string)),
			InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
			State:          entity.OrderState(result["state"].(string)),
			CreatedAt:      result["created_at"].(int64),
			UpdatedAt:      result["updated_at"].(int64),
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByInvestor(investor common.Address) ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE investor = ?", investor.String()).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by investor: %w", err)
	}

	var orders []*entity.Order
	for _, result := range results {
		order := &entity.Order{
			Id:             uint(result["id"].(int64)),
			CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(result["investor"].(string)),
			Amount:         uint256.MustFromHex(result["amount"].(string)),
			InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
			State:          entity.OrderState(result["state"].(string)),
			CreatedAt:      result["created_at"].(int64),
			UpdatedAt:      result["updated_at"].(int64),
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositorySqlite) FindAllOrders() ([]*entity.Order, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}

	var orders []*entity.Order
	for _, result := range results {
		order := &entity.Order{
			Id:             uint(result["id"].(int64)),
			CrowdfundingId: uint(result["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(result["investor"].(string)),
			Amount:         uint256.MustFromHex(result["amount"].(string)),
			InterestRate:   uint256.MustFromHex(result["interest_rate"].(string)),
			State:          entity.OrderState(result["state"].(string)),
			CreatedAt:      result["created_at"].(int64),
			UpdatedAt:      result["updated_at"].(int64),
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepositorySqlite) UpdateOrder(input *entity.Order) (*entity.Order, error) {
	var orderJSON map[string]interface{}
	err := r.Db.Where("id = ?", input.Id).First(&orderJSON).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrOrderNotFound
		}
		return nil, err
	}

	if input.Amount.Sign() > 0 {
		orderJSON["amount"] = input.Amount.String()
	}
	if input.InterestRate.Sign() > 0 {
		orderJSON["interest_rate"] = input.InterestRate.String()
	}
	if input.State != "" {
		orderJSON["state"] = input.State
	}
	orderJSON["updated_at"] = input.UpdatedAt

	orderBytes, err := json.Marshal(orderJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order: %w", err)
	}

	var order entity.Order
	if err = json.Unmarshal(orderBytes, &order); err != nil {
		return nil, err
	}
	if err = order.Validate(); err != nil {
		return nil, entity.ErrInvalidOrder
	}

	res := r.Db.Save(&order)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update order: %w", res.Error)
	}
	return &order, nil
}

func (r *OrderRepositorySqlite) DeleteOrder(id uint) error {
	res := r.Db.Delete(&entity.Order{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete order: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}
