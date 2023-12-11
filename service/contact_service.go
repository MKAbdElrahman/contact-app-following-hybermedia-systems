package service

import (
	"app/domain"
	"context"
)

type ContactStore interface {
	AddContact(ctx context.Context, contact domain.Contact) error
	GetContacts(ctx context.Context, query string) ([]domain.Contact, error)
	GetContactByID(ctx context.Context, id int) (*domain.Contact, error)
}

type ContactService struct {
	ContactStore ContactStore
}

func NewContactService(store ContactStore) *ContactService {
	return &ContactService{
		ContactStore: store,
	}
}
