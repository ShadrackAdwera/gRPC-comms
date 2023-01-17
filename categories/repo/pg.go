package repo

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Category: Category{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	Category Category
}

// Category is the structure which holds one category from the database.
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdy"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAll returns a slice of all categories
func (u *Category) GetAll() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM categories;`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedBy,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// GetOne returns one category by id
func (u *Category) GetOne(id int) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM categories WHERE id = $1`

	var category Category
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedBy,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update updates one category in the database, using the information
// stored in the receiver u
func (u *Category) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE categories SET
		name = $1,
		description = $2,
		updated_at = $3
		where id = $4
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Name,
		u.Description,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one category from the database, by Category.ID
func (u *Category) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM categories WHERE id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new category into the database, and returns the ID of the newly inserted row
func (u *Category) Insert(category Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `INSERT INTO categories (name, description, created_by, created_at, updated_at)
		values ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRowContext(ctx, stmt,
		category.Name,
		category.Description,
		category.CreatedBy,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
