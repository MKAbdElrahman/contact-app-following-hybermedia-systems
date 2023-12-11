package db

import (
	"app/domain"
	"app/service"
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

func (store *InMemoryContactStore) GetContactByID(ctx context.Context, id int) (*domain.Contact, error) {
	for _, contact := range store.contacts {
		if contact.ID == id {
			return &contact, nil
		}
	}
	return nil, domain.ErrorContactNotFound
}

func (store *InMemoryContactStore) UpdateContact(ctx context.Context, id int, updateParams service.ContactUpdateParams) error {
	for i, contact := range store.contacts {
		if contact.ID == id {
			// Update the contact fields if non-empty in the updateParams
			if updateParams.FirstName != "" {
				store.contacts[i].FirstName = updateParams.FirstName
			}
			if updateParams.LastName != "" {
				store.contacts[i].LastName = updateParams.LastName
			}
			if updateParams.Phone != "" {
				store.contacts[i].Phone = updateParams.Phone
			}
			if updateParams.Email != "" {
				store.contacts[i].Email = updateParams.Email
			}
			return nil
		}
	}
	return domain.ErrorContactNotFound
}
