package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionId int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionId)
	if err != nil {
		return nil, err
	}

	var (
		valueStrings []string
		valueArgs    []interface{}
	)

	for i, d := range details {
		details[i].TransactionID = transactionId

		base := i * 4
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d)",
				base+1, base+2, base+3, base+4),
		)

		valueArgs = append(valueArgs,
			transactionId,
			d.ProductID,
			d.Quantity,
			d.Subtotal,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO transaction_details
		(transaction_id, product_id, quantity, subtotal)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionId,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
