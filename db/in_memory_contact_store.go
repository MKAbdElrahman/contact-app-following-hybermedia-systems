package db

import (
	"app/domain"
	"context"
)

type InMemoryContactStore struct {
	contacts []domain.Contact
}

func NewInMemoryContactStore() *InMemoryContactStore {
	return &InMemoryContactStore{
		contacts: make([]domain.Contact, 0),
	}
}

func (store *InMemoryContactStore) AddContact(ctx context.Context, contact domain.Contact) error {

	store.contacts = append(store.contacts, contact)
	return nil
}

func (store *InMemoryContactStore) GetContacts(ctx context.Context) ([]domain.Contact, error) {
	return store.contacts, nil
}
