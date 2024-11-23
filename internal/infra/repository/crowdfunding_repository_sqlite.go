package repository

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type CrowdfundingRepositorySqlite struct {
	Db *gorm.DB
}

func NewCrowdfundingRepositorySqlite(db *gorm.DB) *CrowdfundingRepositorySqlite {
	return &CrowdfundingRepositorySqlite{Db: db}
}

func (r *CrowdfundingRepositorySqlite) CreateCrowdfunding(input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	if input.TotalObligation == nil {
		input.TotalObligation = uint256.NewInt(0)
	}
	err := r.Db.Raw(`
		INSERT INTO crowdfundings 
		(creator, debt_issued, max_interest_rate, total_obligation, state, expires_at, maturity_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`, input.Creator.String(), input.DebtIssued.Hex(), input.MaxInterestRate.Hex(),
		input.TotalObligation.Hex(), input.State, input.ExpiresAt, input.MaturityAt, input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, err
	}
	return r.FindCrowdfundingById(input.Id)
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingById(id uint) (*entity.Crowdfunding, error) {
	var result map[string]interface{}
	err := r.Db.Raw(`
		SELECT id, creator, debt_issued, max_interest_rate, total_obligation, state, expires_at, maturity_at, created_at, updated_at 
		FROM crowdfundings WHERE id = ? LIMIT 1
	`, id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfunding by id: %w", err)
	}

	crowdfunding := r.mapToCrowdfundingEntity(result)

	var orders []map[string]interface{}
	err = r.Db.Raw(`
		SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at
		FROM orders WHERE crowdfunding_id = ?
	`, crowdfunding.Id).Scan(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
	}

	for _, order := range orders {
		crowdfunding.Orders = append(crowdfunding.Orders, r.mapToOrderEntity(order))
	}

	return crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings() ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(`
		SELECT id, creator, debt_issued, max_interest_rate, total_obligation, state, expires_at, maturity_at, created_at, updated_at
		FROM crowdfundings
	`).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all crowdfundings: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := r.mapToCrowdfundingEntity(result)
		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByCreator(creator common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(`
		SELECT id, creator, debt_issued, max_interest_rate, total_obligation, state, expires_at, maturity_at, created_at, updated_at
		FROM crowdfundings WHERE creator = ?
	`, creator.String()).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := r.mapToCrowdfundingEntity(result)
		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByInvestor(investor common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(`
		SELECT DISTINCT c.id, c.creator, c.debt_issued, c.max_interest_rate, 
		                c.total_obligation, c.state, c.expires_at, c.maturity_at, 
		                c.created_at, c.updated_at
		FROM crowdfundings c
		INNER JOIN orders o ON c.id = o.crowdfunding_id
		WHERE o.investor = ?
	`, investor.String()).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := r.mapToCrowdfundingEntity(result)
		var orders []map[string]interface{}
		err = r.Db.Raw(`
			SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at
			FROM orders
			WHERE crowdfunding_id = ?
		`, crowdfunding.Id).Scan(&orders).Error
		if err != nil {
			return nil, fmt.Errorf("failed to load orders for crowdfunding with id %d: %w", crowdfunding.Id, err)
		}

		for _, order := range orders {
			crowdfunding.Orders = append(crowdfunding.Orders, r.mapToOrderEntity(order))
		}

		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) UpdateCrowdfunding(input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	crowdfunding, err := r.FindCrowdfundingById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to find crowdfunding for update: %w", err)
	}

	if input.DebtIssued != nil && input.DebtIssued.Sign() > 0 {
		crowdfunding.DebtIssued = input.DebtIssued
	}
	if input.MaxInterestRate != nil && input.MaxInterestRate.Sign() > 0 {
		crowdfunding.MaxInterestRate = input.MaxInterestRate
	}
	if input.TotalObligation != nil && input.TotalObligation.Sign() > 0 {
		crowdfunding.TotalObligation = input.TotalObligation
	}
	if input.State != "" {
		crowdfunding.State = input.State
	}
	if input.ExpiresAt != 0 {
		crowdfunding.ExpiresAt = input.ExpiresAt
	}
	if input.MaturityAt != 0 {
		crowdfunding.MaturityAt = input.MaturityAt
	}
	crowdfunding.UpdatedAt = input.UpdatedAt

	res := r.Db.Model(&entity.Crowdfunding{}).Where("id = ?", crowdfunding.Id).Updates(map[string]interface{}{
		"debt_issued":       crowdfunding.DebtIssued.Hex(),
		"max_interest_rate": crowdfunding.MaxInterestRate.Hex(),
		"total_obligation":  crowdfunding.TotalObligation.Hex(),
		"state":             crowdfunding.State,
		"expires_at":        crowdfunding.ExpiresAt,
		"maturity_at":       crowdfunding.MaturityAt,
		"updated_at":        crowdfunding.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update crowdfunding: %w", res.Error)
	}

	return r.FindCrowdfundingById(crowdfunding.Id)
}

func (r *CrowdfundingRepositorySqlite) DeleteCrowdfunding(id uint) error {
	res := r.Db.Delete(&entity.Crowdfunding{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete crowdfunding: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrCrowdfundingNotFound
	}
	return nil
}

// Helper to map query result to Crowdfunding entity
func (r *CrowdfundingRepositorySqlite) mapToCrowdfundingEntity(result map[string]interface{}) *entity.Crowdfunding {
	return &entity.Crowdfunding{
		Id:              uint(result["id"].(int64)),
		Creator:         common.HexToAddress(result["creator"].(string)),
		DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
		MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
		TotalObligation: uint256.MustFromHex(result["total_obligation"].(string)),
		State:           entity.CrowdfundingState(result["state"].(string)),
		ExpiresAt:       result["expires_at"].(int64),
		MaturityAt:      result["maturity_at"].(int64),
		CreatedAt:       result["created_at"].(int64),
		UpdatedAt:       result["updated_at"].(int64),
	}
}

// Helper to map query result to Order entity
func (r *CrowdfundingRepositorySqlite) mapToOrderEntity(order map[string]interface{}) *entity.Order {
	return &entity.Order{
		Id:             uint(order["id"].(int64)),
		CrowdfundingId: uint(order["crowdfunding_id"].(int64)),
		Investor:       common.HexToAddress(order["investor"].(string)),
		Amount:         uint256.MustFromHex(order["amount"].(string)),
		InterestRate:   uint256.MustFromHex(order["interest_rate"].(string)),
		State:          entity.OrderState(order["state"].(string)),
		CreatedAt:      order["created_at"].(int64),
		UpdatedAt:      order["updated_at"].(int64),
	}
}
