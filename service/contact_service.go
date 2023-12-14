package service

import (
	"app/domain"
	"context"
)

type ContactUpdateParams struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type ContactStore interface {
	AddContact(ctx context.Context, contact domain.Contact) error
	UpdateContact(ctx context.Context, id int, params ContactUpdateParams) error
	DeleteContact(ctx context.Context, id int) error
	GetContacts(ctx context.Context, query string) ([]domain.Contact, error)
	GetContactByID(ctx context.Context, id int) (*domain.Contact, error)
	GetContactByEmail(ctx context.Context, email string) (*domain.Contact, error)
}

type ContactService struct {
	ContactStore ContactStore
}

func NewContactService(store ContactStore) *ContactService {
	return &ContactService{
		ContactStore: store,
	}
}
