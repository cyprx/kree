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

type EndpointRepository interface {
	GetAll() ([]*endpoint, error)
	Create(e *endpoint) error
	Update(e *endpoint) error
}

type endpoint struct {
	id     uint
	domain string
	path   string
	status EndpointStatus
}

func (e *endpoint) Url() string {
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

func (r *endpointRepository) GetAll() ([]*endpoint, error) {
	rows, err := r.db.Query("SELECT * FROM metric_endpoints WHERE status = ?", string(EndpointStatus_ACTIVE))
	if err != nil {
		return nil, err
	}
	var endpoints []*endpoint

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

		e := &endpoint{id, domain, path, EndpointStatus(status)}
		endpoints = append(endpoints, e)
	}
	return endpoints, nil
}

func (r *endpointRepository) Create(e *endpoint) error {
	return nil
}

func (r *endpointRepository) Update(e *endpoint) error {
	return nil
}
