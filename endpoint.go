package kree

import (
	"database/sql"
	"log"
)

type EndpointStatus string

const (
	EndpointStatus_ACTIVE   EndpointStatus = "ACTIVE"
	EndpointStatus_INACTIVE EndpointStatus = "INACTIVE"
)

// EndpointRepository is interface for CRUD operations related to Endpoint
type EndpointRepository interface {
	GetAll() ([]*Endpoint, error)
	Create(e *Endpoint) error
	Update(e *Endpoint) error
}

type Endpoint struct {
	id     uint
	domain string
	path   string
	status EndpointStatus
}

func (e *Endpoint) Url() string {
	return e.domain + e.path
}

type endpointRepository struct {
	db *sql.DB
}

func NewEndpointRepository() EndpointRepository {
	db := GetDB()
	if db == nil {
		log.Fatal("no database found")
	}
	return &endpointRepository{db}
}

func (r *endpointRepository) GetAll() ([]*Endpoint, error) {
	rows, err := r.db.Query("SELECT * FROM metric_endpoints WHERE status = ?", string(EndpointStatus_ACTIVE))
	if err != nil {
		return nil, err
	}
	var endpoints []*Endpoint

	for rows.Next() {
		var (
			id     uint
			domain string
			path   string
			status string
		)

		err = rows.Scan(&id, &domain, &path, &status)
		if err != nil {
			return nil, err
		}

		e := &Endpoint{id, domain, path, EndpointStatus(status)}
		endpoints = append(endpoints, e)
	}
	return endpoints, nil
}

func (r *endpointRepository) Create(e *Endpoint) error {
	return nil
}

func (r *endpointRepository) Update(e *Endpoint) error {
	return nil
}
