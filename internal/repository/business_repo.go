package repository

import (
	"database/sql"
	"YALP/internal/domain"
)

type BusinessRepository interface {
	Create(b *domain.Business) (int64, error)
	FindByID(id int64) (*domain.Business, error)
	ListAll() ([]domain.Business, error)
	Search(keyword string) ([]domain.Business, error)
	Update(b *domain.Business) error
}

type businessRepository struct {
	db *sql.DB
}

func NewBusinessRepository(db *sql.DB) BusinessRepository {
	return &businessRepository{db: db}
}

func (r *businessRepository) Create(b *domain.Business) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO businesses (name, category, description, address, contact_info, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id
	`, b.Name, b.Category, b.Description, b.Address, b.ContactInfo, b.OwnerID).Scan(&id)
	return id, err
}

func (r *businessRepository) FindByID(id int64) (*domain.Business, error) {
	var b domain.Business
	err := r.db.QueryRow(`SELECT id, name, category, description, address, contact_info, owner_id, created_at, updated_at FROM businesses WHERE id=$1`, id).
		Scan(&b.ID, &b.Name, &b.Category, &b.Description, &b.Address, &b.ContactInfo, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &b, err
}

func (r *businessRepository) ListAll() ([]domain.Business, error) {
	rows, err := r.db.Query(`SELECT id, name, category, description, address, contact_info, owner_id, created_at, updated_at FROM businesses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []domain.Business
	for rows.Next() {
		var b domain.Business
		if err := rows.Scan(&b.ID, &b.Name, &b.Category, &b.Description, &b.Address, &b.ContactInfo, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, b)
	}
	return businesses, nil
}

func (r *businessRepository) Search(keyword string) ([]domain.Business, error) {
	rows, err := r.db.Query(`
		SELECT id, name, category, description, address, contact_info, owner_id, created_at, updated_at
		FROM businesses
		WHERE name ILIKE '%' || $1 || '%' OR category ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
	`, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []domain.Business
	for rows.Next() {
		var b domain.Business
		if err := rows.Scan(&b.ID, &b.Name, &b.Category, &b.Description, &b.Address, &b.ContactInfo, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, b)
	}
	return businesses, nil
}

func (r *businessRepository) Update(b *domain.Business) error {
	_, err := r.db.Exec(`
		UPDATE businesses SET name=$1, category=$2, description=$3, address=$4, contact_info=$5, owner_id=$6, updated_at=NOW()
		WHERE id=$7
	`, b.Name, b.Category, b.Description, b.Address, b.ContactInfo, b.OwnerID, b.ID)
	return err
}
