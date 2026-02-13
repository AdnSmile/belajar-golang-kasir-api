package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(nameFilter string) ([]models.Product, error) {

	query := "SELECT p.id, p.name, p.price, p.stock, p.category_id FROM products p"

	args := []interface{}{}
	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%") // buat tambahan argument
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt16
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			p.CategoryID = &categoryID.Int16
		} else {
			p.CategoryID = nil
		}

		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.ProductDetail, error) {

	query := "SELECT s.id, s.name, s.price, s.stock, c.name FROM products s LEFT JOIN categories c ON s.category_id = c.id WHERE s.id = $1"

	var p models.ProductDetail
	var categoryName sql.NullString

	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	if categoryName.Valid {
		p.CategoryName = &categoryName.String
	} else {
		p.CategoryName = nil
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {

	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {

	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
