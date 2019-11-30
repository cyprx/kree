package kree

import (
	"database/sql"
)

type EndpointRepository interface {
	GetAll() ([]*endpoint, error)
	Create(e *endpoint) error
	Update(e *endpoint) error
}

type endpoint struct {
	domain string
	path   string
}

func (e *endpoint) Url() string {
	return e.domain + e.path
}

type endpointRepository struct {
	db *sql.DB
}

func NewEndpointRepository() EndpointRepository {
	db := GetDB()
	return &endpointRepository{db}
}

func (r *endpointRepository) GetAll() ([]*endpoint, error) {
	return nil, nil
}

func (r *endpointRepository) Create(e *endpoint) error {
	return nil
}

func (r *endpointRepository) Update(e *endpoint) error {
	return nil
}
