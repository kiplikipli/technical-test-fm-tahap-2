package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kiplikipli/technical-test-fm-tahap-2/database"
	"github.com/kiplikipli/technical-test-fm-tahap-2/entity"
	"gorm.io/gorm"
)

type Transaction entity.Transaction

type NewTransactionRequest struct {
	UserID              uuid.UUID      `json:"user_id" validate:"required"`
	Type                sql.NullString `json:"type" validate:"required"`
	Amount              int64          `json:"amount" validate:"required"`
	Remarks             string         `json:"remarks"`
	Category            string         `json:"category"`
	CorrespondingUserID uuid.UUID      `json:"corresponding_user_id"`
}

func CreateDebitTransaction(targetUserId uuid.UUID, request NewTransactionRequest) (*Transaction, error) {
	db := database.DB
	transaction := &Transaction{}

	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.First(&user, &User{ID: targetUserId}).Error
		if err != nil {
			return err
		}

		if user.Balance < request.Amount {
			return errors.New("balance is not enough")
		}

		transaction = &Transaction{
			ID:            uuid.New(),
			UserID:        targetUserId,
			Type:          "DEBIT",
			Amount:        request.Amount,
			Remarks:       request.Remarks,
			Status:        "PENDING",
			BalanceBefore: user.Balance,
			BalanceAfter:  user.Balance - request.Amount,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		user.Balance = transaction.BalanceAfter
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func CreateCreditTransaction(targetUserId uuid.UUID, request NewTransactionRequest) (*Transaction, error) {
	db := database.DB
	transaction := &Transaction{}

	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.First(&user, &User{ID: targetUserId}).Error
		if err != nil {
			return err
		}

		transaction = &Transaction{
			ID:            uuid.New(),
			UserID:        targetUserId,
			Type:          "CREDIT",
			Amount:        request.Amount,
			Remarks:       request.Remarks,
			Status:        "PENDING",
			BalanceBefore: user.Balance,
			BalanceAfter:  user.Balance + request.Amount,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		user.Balance = transaction.BalanceAfter
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func CreateDebitTransactionWithDbTransaction(targetUserId uuid.UUID, request NewTransactionRequest, tx *gorm.DB) (*Transaction, error) {
	transaction := &Transaction{}
	var user User
	err := tx.First(&user, &User{ID: targetUserId}).Error
	if err != nil {
		return nil, err
	}

	if user.Balance < request.Amount {
		return nil, errors.New("balance is not enough")
	}

	transaction = &Transaction{
		ID:                  uuid.New(),
		UserID:              targetUserId,
		Type:                "DEBIT",
		Amount:              request.Amount,
		Remarks:             request.Remarks,
		Status:              "PENDING",
		BalanceBefore:       user.Balance,
		BalanceAfter:        user.Balance - request.Amount,
		CreatedAt:           time.Now(),
		CorrespondingUserID: &request.CorrespondingUserID,
	}

	if err := tx.Create(transaction).Error; err != nil {
		return nil, err
	}

	user.Balance = transaction.BalanceAfter
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func CreateCreditTransactionWithDbTransaction(targetUserId uuid.UUID, request NewTransactionRequest, tx *gorm.DB) (*Transaction, error) {
	transaction := &Transaction{}
	var user User
	err := tx.First(&user, &User{ID: targetUserId}).Error
	if err != nil {
		return nil, err
	}

	transaction = &Transaction{
		ID:                  uuid.New(),
		UserID:              targetUserId,
		Type:                "CREDIT",
		Amount:              request.Amount,
		Remarks:             request.Remarks,
		Status:              "PENDING",
		BalanceBefore:       user.Balance,
		BalanceAfter:        user.Balance + request.Amount,
		CreatedAt:           time.Now(),
		CorrespondingUserID: &request.CorrespondingUserID,
	}

	if err := tx.Create(transaction).Error; err != nil {
		return nil, err
	}

	user.Balance = transaction.BalanceAfter
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func CreateTransferTransaction(targetUserId uuid.UUID, requests []NewTransactionRequest) (*Transaction, error) {
	transactions, err := CreateMultipleTransactions(requests)
	if err != nil {
		return nil, err
	}

	selfTransaction := &Transaction{}
	for i := 0; i < len(transactions); i++ {
		if transactions[i].UserID == targetUserId {
			selfTransaction = transactions[i]
		}
	}

	return selfTransaction, nil
}

func CreateMultipleTransactions(requests []NewTransactionRequest) ([]*Transaction, error) {
	db := database.DB
	transactions := []*Transaction{}

	err := db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(requests); i++ {
			request := requests[i]
			if !request.Type.Valid {
				return errors.New("type is required")
			}

			if request.Type.String == "DEBIT" {
				transaction, err := CreateDebitTransactionWithDbTransaction(request.UserID, request, tx)
				if err != nil {
					return err
				}
				transactions = append(transactions, transaction)
			} else if request.Type.String == "CREDIT" {
				transaction, err := CreateCreditTransactionWithDbTransaction(request.UserID, request, tx)
				if err != nil {
					return err
				}
				transactions = append(transactions, transaction)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
