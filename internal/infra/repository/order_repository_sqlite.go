package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/datatype"
	"gorm.io/gorm"
)

type OrderRepositorySqlite struct {
	Db *gorm.DB
}

func NewOrderRepositorySqlite(db *gorm.DB) *OrderRepositorySqlite {
	return &OrderRepositorySqlite{Db: db}
}

func (r *OrderRepositorySqlite) CreateOrder(ctx context.Context, input *entity.Order) (*entity.Order, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return input, nil
}

func (r *OrderRepositorySqlite) FindOrderById(ctx context.Context, id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.Db.WithContext(ctx).First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}
	return &order, nil
}

func (r *OrderRepositorySqlite) FindOrdersByCrowdfundingId(ctx context.Context, id uint) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.WithContext(ctx).Where("crowdfunding_id = ?", id).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by crowdfunding ID: %w", err)
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByState(ctx context.Context, crowdfundingId uint, state string) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.WithContext(ctx).
		Where("crowdfunding_id = ? AND state = ?", crowdfundingId, state).
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindOrdersByInvestor(ctx context.Context, investor datatype.Address) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.WithContext(ctx).Where("investor = ?", investor).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by investor: %w", err)
	}
	return orders, nil
}

func (r *OrderRepositorySqlite) FindAllOrders(ctx context.Context) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.WithContext(ctx).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}
	return orders, nil
}code .

	}
	if && input.InterestRate.Sign() > 0 {
		order.InterestRate = input.InterestRate
	}
	if input.State != "" {
		order.State = input.State
	}
	order.UpdatedAt = input.UpdatedAt

	if err := r.Db.WithContext(ctx).Save(order).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return order, nil
}

func (r *OrderRepositorySqlite) DeleteOrder(ctx context.Context, id uint) error {
	res := r.Db.WithContext(ctx).Delete(&entity.Order{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete order: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}
