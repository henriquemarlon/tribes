package repository

import (
	"context"
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

func (r *CrowdfundingRepositorySqlite) CreateCrowdfunding(ctx context.Context, input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	if input.TotalObligation == nil {
		input.TotalObligation = uint256.NewInt(0)
	}
	err := r.Db.WithContext(ctx).Raw(`
		INSERT INTO crowdfundings (token, amount, creator, debt_issued, max_interest_rate, total_obligation, state, fundraising_duration, closes_at, maturity_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`, input.Token.String(), input.Amount.Hex(), input.Creator.String(), input.DebtIssued.Hex(), input.MaxInterestRate.Hex(),
		input.TotalObligation.Hex(), input.State, input.FundraisingDuration, input.ClosesAt,
		input.MaturityAt, input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create crowdfunding: %w", err)
	}
	return input, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingById(ctx context.Context, id uint) (*entity.Crowdfunding, error) {
	var result map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, token, amount, creator, debt_issued, max_interest_rate, total_obligation, state, fundraising_duration, closes_at, maturity_at, created_at, updated_at
		FROM crowdfundings WHERE id = ? LIMIT 1
	`, id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfunding by id: %w", err)
	}

	crowdfunding := &entity.Crowdfunding{
		Id:                  uint(result["id"].(int64)),
		Token:               common.HexToAddress(result["token"].(string)),
		Amount:              uint256.MustFromHex(result["amount"].(string)),
		Creator:             common.HexToAddress(result["creator"].(string)),
		DebtIssued:          uint256.MustFromHex(result["debt_issued"].(string)),
		MaxInterestRate:     uint256.MustFromHex(result["max_interest_rate"].(string)),
		TotalObligation:     uint256.MustFromHex(result["total_obligation"].(string)),
		State:               entity.CrowdfundingState(result["state"].(string)),
		FundraisingDuration: result["fundraising_duration"].(int64),
		ClosesAt:            result["closes_at"].(int64),
		MaturityAt:          result["maturity_at"].(int64),
		CreatedAt:           result["created_at"].(int64),
		UpdatedAt:           result["updated_at"].(int64),
	}

	var orders []map[string]interface{}
	err = r.Db.WithContext(ctx).Raw(`
		SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at
		FROM orders WHERE crowdfunding_id = ?
	`, crowdfunding.Id).Scan(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
	}

	for _, order := range orders {
		crowdfunding.Orders = append(crowdfunding.Orders, &entity.Order{
			Id:             uint(order["id"].(int64)),
			CrowdfundingId: uint(order["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(order["investor"].(string)),
			Amount:         uint256.MustFromHex(order["amount"].(string)),
			InterestRate:   uint256.MustFromHex(order["interest_rate"].(string)),
			State:          entity.OrderState(order["state"].(string)),
			CreatedAt:      order["created_at"].(int64),
			UpdatedAt:      order["updated_at"].(int64),
		})
	}

	return crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings(ctx context.Context) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, token, amount, creator, debt_issued, max_interest_rate, total_obligation, state, fundraising_duration, closes_at, maturity_at, created_at, updated_at
		FROM crowdfundings
	`).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all crowdfundings: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, data := range results {
		crowdfundings = append(crowdfundings, &entity.Crowdfunding{
			Id:                  uint(data["id"].(int64)),
			Token:               common.HexToAddress(data["token"].(string)),
			Amount:              uint256.MustFromHex(data["amount"].(string)),
			Creator:             common.HexToAddress(data["creator"].(string)),
			DebtIssued:          uint256.MustFromHex(data["debt_issued"].(string)),
			MaxInterestRate:     uint256.MustFromHex(data["max_interest_rate"].(string)),
			TotalObligation:     uint256.MustFromHex(data["total_obligation"].(string)),
			State:               entity.CrowdfundingState(data["state"].(string)),
			FundraisingDuration: data["fundraising_duration"].(int64),
			ClosesAt:            data["closes_at"].(int64),
			MaturityAt:          data["maturity_at"].(int64),
			CreatedAt:           data["created_at"].(int64),
			UpdatedAt:           data["updated_at"].(int64),
		})
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByInvestor(ctx context.Context, investor common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT DISTINCT c.id, c.token, c.amount, c.creator, c.debt_issued, c.max_interest_rate, 
		                c.total_obligation, c.state, c.fundraising_duration, c.closes_at, c.maturity_at, 
		                c.created_at, c.updated_at
		FROM crowdfundings c
		INNER JOIN orders o ON c.id = o.crowdfunding_id
		WHERE o.investor = ?
	`, investor.String()).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, data := range results {
		crowdfunding := &entity.Crowdfunding{
			Id:                  uint(data["id"].(int64)),
			Token:               common.HexToAddress(data["token"].(string)),
			Amount:              uint256.MustFromHex(data["amount"].(string)),
			Creator:             common.HexToAddress(data["creator"].(string)),
			DebtIssued:          uint256.MustFromHex(data["debt_issued"].(string)),
			MaxInterestRate:     uint256.MustFromHex(data["max_interest_rate"].(string)),
			TotalObligation:     uint256.MustFromHex(data["total_obligation"].(string)),
			State:               entity.CrowdfundingState(data["state"].(string)),
			FundraisingDuration: data["fundraising_duration"].(int64),
			ClosesAt:            data["closes_at"].(int64),
			MaturityAt:          data["maturity_at"].(int64),
			CreatedAt:           data["created_at"].(int64),
			UpdatedAt:           data["updated_at"].(int64),
		}

		var orders []map[string]interface{}
		err = r.Db.WithContext(ctx).Raw(`
			SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at
			FROM orders WHERE crowdfunding_id = ?
		`, crowdfunding.Id).Scan(&orders).Error
		if err != nil {
			return nil, fmt.Errorf("failed to load orders for crowdfunding with id %d: %w", crowdfunding.Id, err)
		}

		for _, order := range orders {
			crowdfunding.Orders = append(crowdfunding.Orders, &entity.Order{
				Id:             uint(order["id"].(int64)),
				CrowdfundingId: uint(order["crowdfunding_id"].(int64)),
				Investor:       common.HexToAddress(order["investor"].(string)),
				Amount:         uint256.MustFromHex(order["amount"].(string)),
				InterestRate:   uint256.MustFromHex(order["interest_rate"].(string)),
				State:          entity.OrderState(order["state"].(string)),
				CreatedAt:      order["created_at"].(int64),
				UpdatedAt:      order["updated_at"].(int64),
			})
		}
		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByCreator(ctx context.Context, creator common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, token, amount, creator, debt_issued, max_interest_rate, total_obligation, state, fundraising_duration, closes_at, maturity_at, created_at, updated_at
		FROM crowdfundings WHERE creator = ?
	`, creator.String()).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfundings = append(crowdfundings, &entity.Crowdfunding{
			Id:                  uint(result["id"].(int64)),
			Token:               common.HexToAddress(result["token"].(string)),
			Amount:              uint256.MustFromHex(result["amount"].(string)),
			Creator:             common.HexToAddress(result["creator"].(string)),
			DebtIssued:          uint256.MustFromHex(result["debt_issued"].(string)),
			MaxInterestRate:     uint256.MustFromHex(result["max_interest_rate"].(string)),
			TotalObligation:     uint256.MustFromHex(result["total_obligation"].(string)),
			State:               entity.CrowdfundingState(result["state"].(string)),
			FundraisingDuration: result["fundraising_duration"].(int64),
			ClosesAt:            result["closes_at"].(int64),
			MaturityAt:          result["maturity_at"].(int64),
			CreatedAt:           result["created_at"].(int64),
			UpdatedAt:           result["updated_at"].(int64),
		})
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) UpdateCrowdfunding(ctx context.Context, input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	crowdfunding, err := r.FindCrowdfundingById(ctx, input.Id)
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
	if input.FundraisingDuration != 0 {
		crowdfunding.FundraisingDuration = input.FundraisingDuration
	}
	if input.ClosesAt != 0 {
		crowdfunding.ClosesAt = input.ClosesAt
	}
	if input.MaturityAt != 0 {
		crowdfunding.MaturityAt = input.MaturityAt
	}
	crowdfunding.UpdatedAt = input.UpdatedAt

	res := r.Db.WithContext(ctx).Model(&entity.Crowdfunding{}).Where("id = ?", crowdfunding.Id).Updates(map[string]interface{}{
		"creator":              crowdfunding.Creator.String(),
		"debt_issued":          crowdfunding.DebtIssued.Hex(),
		"max_interest_rate":    crowdfunding.MaxInterestRate.Hex(),
		"total_obligation":     crowdfunding.TotalObligation.Hex(),
		"state":                crowdfunding.State,
		"fundraising_duration": crowdfunding.FundraisingDuration,
		"closes_at":            crowdfunding.ClosesAt,
		"maturity_at":          crowdfunding.MaturityAt,
		"updated_at":           crowdfunding.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update crowdfunding: %w", res.Error)
	}

	return r.FindCrowdfundingById(ctx, crowdfunding.Id)
}

func (r *CrowdfundingRepositorySqlite) DeleteCrowdfunding(ctx context.Context, id uint) error {
	res := r.Db.Delete(&entity.Crowdfunding{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete crowdfunding: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrCrowdfundingNotFound
	}
	return nil
}
