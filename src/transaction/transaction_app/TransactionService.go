package transaction_app

import (
	"loyalty-campaigns/src/common/models"
	"loyalty-campaigns/src/common/utils"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_ports"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_requests"
	"loyalty-campaigns/src/transaction/transaction_domain/transaction_structs/transaction_responses"
	"sync"
	"time"
)

type ITransactionService interface {
	CreateTransaction(req transaction_requests.CreateTransactionRequest) (*transaction_responses.TransactionResponse, error)
	GetTransaction(id uint) (*transaction_responses.TransactionResponse, error)
	ListTransactionsByUser(userID uint) ([]transaction_responses.TransactionResponse, error)
	ListTransactionsByBranch(branchID uint) ([]transaction_responses.TransactionResponse, error)
	GetTransactionsByDateRange(startDate, endDate time.Time) ([]transaction_responses.TransactionResponse, error)
	GetTotalAmountByUserAndDateRange(userID uint, startDate, endDate time.Time) (float64, error)
}

type transactionService struct {
	transactionRepo transaction_ports.ITransactionRepository
	logger          utils.ILogger
}

var (
	transactionServiceInstance *transactionService
	transactionServiceOnce     sync.Once
)

func NewTransactionService(transactionRepo transaction_ports.ITransactionRepository) ITransactionService {
	transactionServiceOnce.Do(func() {
		transactionServiceInstance = &transactionService{
			transactionRepo: transactionRepo,
			logger:          utils.NewLogger(),
		}
	})
	return transactionServiceInstance
}

func (s *transactionService) CreateTransaction(req transaction_requests.CreateTransactionRequest) (*transaction_responses.TransactionResponse, error) {
	transaction := &models.Transaction{
		UserID:   req.UserID,
		BranchID: req.BranchID,
		Amount:   req.Amount,
		Date:     req.Date,
	}

	err := s.transactionRepo.Create(transaction)
	if err != nil {
		s.logger.Error("Error al crear transacción", err)
		return nil, err
	}

	return mapTransactionToResponse(transaction), nil
}

func (s *transactionService) GetTransaction(id uint) (*transaction_responses.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener transacción", err)
		return nil, err
	}

	return mapTransactionToResponse(transaction), nil
}

func (s *transactionService) ListTransactionsByUser(userID uint) ([]transaction_responses.TransactionResponse, error) {
	transactions, err := s.transactionRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("Error al listar transacciones del usuario", err)
		return nil, err
	}

	return mapTransactionsToResponses(transactions), nil
}

func (s *transactionService) ListTransactionsByBranch(branchID uint) ([]transaction_responses.TransactionResponse, error) {
	transactions, err := s.transactionRepo.GetByBranchID(branchID)
	if err != nil {
		s.logger.Error("Error al listar transacciones de la sucursal", err)
		return nil, err
	}

	return mapTransactionsToResponses(transactions), nil
}

func (s *transactionService) GetTransactionsByDateRange(startDate, endDate time.Time) ([]transaction_responses.TransactionResponse, error) {
	transactions, err := s.transactionRepo.GetByDateRange(startDate, endDate)
	if err != nil {
		s.logger.Error("Error al obtener transacciones por rango de fechas", err)
		return nil, err
	}

	return mapTransactionsToResponses(transactions), nil
}

func (s *transactionService) GetTotalAmountByUserAndDateRange(userID uint, startDate, endDate time.Time) (float64, error) {
	totalAmount, err := s.transactionRepo.GetTotalAmountByUserAndDateRange(userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Error al obtener monto total de transacciones del usuario por rango de fechas", err)
		return 0, err
	}

	return totalAmount, nil
}

func mapTransactionToResponse(transaction *models.Transaction) *transaction_responses.TransactionResponse {
	return &transaction_responses.TransactionResponse{
		ID:       transaction.ID,
		UserID:   transaction.UserID,
		BranchID: transaction.BranchID,
		Amount:   transaction.Amount,
		Date:     transaction.Date,
	}
}

func mapTransactionsToResponses(transactions []models.Transaction) []transaction_responses.TransactionResponse {
	responses := make([]transaction_responses.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *mapTransactionToResponse(&transaction)
	}
	return responses
}
