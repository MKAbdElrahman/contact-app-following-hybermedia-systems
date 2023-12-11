package db

import (
	"app/domain"
	"context"
	"strings"
)

type InMemoryContactStore struct {
	contacts []domain.Contact
}

func NewInMemoryContactStore() *InMemoryContactStore {
	return &InMemoryContactStore{
		contacts: make([]domain.Contact, 0),
	}
}

var counter = 0

func (store *InMemoryContactStore) AddContact(ctx context.Context, contact domain.Contact) error {
	contact.ID = counter
	store.contacts = append(store.contacts, contact)
	counter++
	return nil
}

func (store *InMemoryContactStore) GetContacts(ctx context.Context, query string) ([]domain.Contact, error) {
	var results []domain.Contact

	// Case-insensitive substring search
	query = strings.ToLower(query)

	for _, contact := range store.contacts {
		if strings.Contains(strings.ToLower(contact.FirstName), query) ||
			strings.Contains(strings.ToLower(contact.LastName), query) ||
			strings.Contains(strings.ToLower(contact.Phone), query) ||
			strings.Contains(strings.ToLower(contact.Email), query) {
			results = append(results, contact)
		}
	}

	return results, nil
}
