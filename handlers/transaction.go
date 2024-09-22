package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kiplikipli/technical-test-fm-tahap-2/services"
)

type (
	CreateTopUpRequest struct {
		Amount  int64  `json:"amount" validate:"required"`
		Remarks string `json:"remarks"`
	}

	CreateTopUpResponse struct {
		TopUpID       string `json:"top_up_id"`
		AmountTopUp   int64  `json:"amount_top_up"`
		BalanceBefore int64  `json:"balance_before"`
		BalanceAfter  int64  `json:"balance_after"`
		CreatedDate   string `json:"created_date"`
	}

	CreatePaymentRequest struct {
		Amount  int64  `json:"amount" validate:"required"`
		Remarks string `json:"remarks"`
	}

	CreatePaymentResponse struct {
		PaymentID     string `json:"payment_id"`
		Amount        int64  `json:"amount"`
		Remarks       string `json:"remarks"`
		BalanceBefore int64  `json:"balance_before"`
		BalanceAfter  int64  `json:"balance_after"`
		CreatedDate   string `json:"created_date"`
	}

	CreateTransferRequest struct {
		Amount     int64  `json:"amount" validate:"required"`
		TargetUser string `json:"target_user" validate:"required"`
		Remarks    string `json:"remarks"`
	}

	CreateTransferResponse struct {
		TransferID    string `json:"transfer_id"`
		Amount        int64  `json:"amount"`
		Remarks       string `json:"remarks"`
		BalanceBefore int64  `json:"balance_before"`
		BalanceAfter  int64  `json:"balance_after"`
		CreatedDate   string `json:"created_date"`
	}
)

func CreateTopUp(c *fiber.Ctx) error {
	userUuid, err := extractUserUuidFromContext(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid UUID",
		})
	}

	json := new(CreateTopUpRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	newTransaction := services.NewTransactionRequest{
		UserID:   userUuid,
		Amount:   json.Amount,
		Remarks:  json.Remarks,
		Category: "TopUp",
	}
	transaction, err := services.CreateCreditTransaction(userUuid, newTransaction)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": &CreateTopUpResponse{
			TopUpID:       transaction.ID.String(),
			AmountTopUp:   transaction.Amount,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedDate:   transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func CreatePayment(c *fiber.Ctx) error {
	userUuid, err := extractUserUuidFromContext(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid UUID",
		})
	}

	json := new(CreateTopUpRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	newTransaction := services.NewTransactionRequest{
		UserID:   userUuid,
		Amount:   json.Amount,
		Remarks:  json.Remarks,
		Category: "Payment",
	}
	transaction, err := services.CreateDebitTransaction(userUuid, newTransaction)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": &CreatePaymentResponse{
			PaymentID:     transaction.ID.String(),
			Amount:        transaction.Amount,
			Remarks:       transaction.Remarks,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedDate:   transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func CreateTransfer(c *fiber.Ctx) error {
	userUuid, err := extractUserUuidFromContext(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid UUID",
		})
	}

	json := new(CreateTransferRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	targetUserUuid, err := uuid.Parse(json.TargetUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Target User UUID",
		})
	}

	newTransactions := []services.NewTransactionRequest{
		{
			UserID:   userUuid,
			Amount:   json.Amount,
			Remarks:  json.Remarks,
			Category: "Transfer",
			Type: sql.NullString{
				String: "DEBIT",
				Valid:  true,
			},
			CorrespondingUserID: targetUserUuid,
		},
		{
			UserID:   targetUserUuid,
			Amount:   json.Amount,
			Remarks:  json.Remarks,
			Category: "Transfer",
			Type: sql.NullString{
				String: "CREDIT",
				Valid:  true,
			},
			CorrespondingUserID: userUuid,
		},
	}
	transaction, err := services.CreateTransferTransaction(userUuid, newTransactions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": &CreateTransferResponse{
			TransferID:    transaction.ID.String(),
			Amount:        transaction.Amount,
			Remarks:       transaction.Remarks,
			BalanceBefore: transaction.BalanceBefore,
			BalanceAfter:  transaction.BalanceAfter,
			CreatedDate:   transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func extractUserUuidFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	stringUuid := c.Locals("userInfo").(jwt.MapClaims)["user_id"].(string)
	return uuid.Parse(stringUuid)
}
