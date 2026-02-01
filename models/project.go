package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	TechStack     json.RawMessage `json:"tech_stack"`
	Slug          string          `json:"slug"`
	DisplayStatus bool            `json:"display_status"`
	Image         *string         `json:"image"`
	Description   *string         `json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetAll returns all projects with display_status = true
func (r *ProjectRepository) GetAll() ([]Project, error) {
	query := `
		SELECT id, name, tech_stack, slug, display_status, image, description, created_at, updated_at
		FROM projects
		WHERE display_status = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.TechStack,
			&p.Slug,
			&p.DisplayStatus,
			&p.Image,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

// GetBySlug returns a project by its slug
func (r *ProjectRepository) GetBySlug(slug string) (*Project, error) {
	query := `
		SELECT id, name, tech_stack, slug, display_status, image, description, created_at, updated_at
		FROM projects
		WHERE slug = $1 AND display_status = true
	`

	var p Project
	err := r.db.QueryRow(context.Background(), query, slug).Scan(
		&p.ID,
		&p.Name,
		&p.TechStack,
		&p.Slug,
		&p.DisplayStatus,
		&p.Image,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}
