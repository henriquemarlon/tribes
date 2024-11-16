package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
	err := r.Db.Create(input).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return input, nil
}

func (r *OrderRepositorySqlite) FindOrdersByState(crowdfundingId uint, state string) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := r.Db.Where("crowdfunding_id = ? AND state = ?", crowdfundingId, state).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByInvestor(investor common.Address) ([]*entity.Order, error) {
	return nil, nil
}

func (r *OrderRepositorySqlite) FindOrderById(id uint) (*entity.Order, error) {
	var order entity.Order
	err := r.Db.First(&order, "order_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}
	return &order, nil
}

func (r *OrderRepositorySqlite) FindOrdersByCrowdfundingId(id uint) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := r.Db.Where("crowdfunding_id = ?", id).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by crowdfunding ID: %w", err)
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindAllOrders() ([]*entity.Order, error) {
	var orders []*entity.Order
	err := r.Db.Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
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

	orderJSON["amount"] = input.Amount
	orderJSON["interest_rate"] = input.InterestRate
	orderJSON["state"] = input.State
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

	res := r.Db.Save(orderJSON)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update order: %w", res.Error)
	}
	return &order, nil
}

func (r *OrderRepositorySqlite) DeleteOrder(id uint) error {
	res := r.Db.Delete(&entity.Order{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}
