package repository

import (
	"context"
	"database/sql"
	"strings"
)

type CategoriesRepository struct {
	db *sql.DB
}

type Category struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"updated_at"`
}

type SelectedCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c *CategoriesRepository) Create(ctx context.Context, name string) error {
	query := `INSERT INTO categories(name) VALUE(?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := c.db.ExecContext(ctx, query, &name)
	if err != nil {
		duplicateKey := "Error 1062"
		switch {
		case strings.Contains(err.Error(), duplicateKey):
			return ErrDuplicateCategory
		default:
			return err
		}
	}

	return nil
}

func (c *CategoriesRepository) GetByID(ctx context.Context, id int) (*Category, error) {
	query := `SELECT  id, name, created_at,updated_at FROM categories WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	cat := &Category{}

	err := c.db.QueryRowContext(ctx, query, id).Scan(&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdateAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}

	return cat, nil
}

func (c *CategoriesRepository) GetCategories(ctx context.Context, qp PaginatedCategoriesQuery) ([]*Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY created_at ` + qp.Sort + ` LIMIT ? OFFSET ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, qp.Limit, qp.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		category := &Category{}
		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoriesRepository) Update(ctx context.Context, category *Category) error {
	query := `UPDATE categories SET name = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := c.db.ExecContext(ctx, query, &category.Name, &category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
