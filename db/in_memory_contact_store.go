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

func (store *InMemoryContactStore) GetContactByEmail(ctx context.Context, email string) (*domain.Contact, error) {
	for _, contact := range store.contacts {
		if contact.Email == email {
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

func (store *InMemoryContactStore) DeleteContact(ctx context.Context, id int) error {
	for i, contact := range store.contacts {
		if contact.ID == id {
			// Remove the contact from the slice by slicing it
			store.contacts = append(store.contacts[:i], store.contacts[i+1:]...)
			return nil
		}
	}
	return domain.ErrorContactNotFound
}

// GetContactsPaginated implementation for InMemoryContactStore
func (store *InMemoryContactStore) GetContactsPaginated(ctx context.Context, query string, page, pageSize int) ([]domain.Contact, error) {
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Check if startIndex is out of bounds
	if startIndex >= len(store.contacts) {
		return nil, nil
	}

	// Adjust endIndex to prevent going beyond the slice length
	if endIndex > len(store.contacts) {
		endIndex = len(store.contacts)
	}

	return store.contacts[startIndex:endIndex], nil
}

// GetTotalContacts implementation for InMemoryContactStore
func (store *InMemoryContactStore) GetTotalContacts(ctx context.Context, query string) (int, error) {
	// Implement your logic to get the total number of contacts based on the query
	// For simplicity, return the total number of contacts in the store
	return len(store.contacts), nil
}
