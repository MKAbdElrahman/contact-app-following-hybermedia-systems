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

func (store *InMemoryContactStore) SeedMe(ctx context.Context) {
	contacts := []domain.Contact{
		{FirstName: "Mohamed", LastName: "Ali", Email: "mohamed.ali@example.com", Phone: "+211111111111"},
		{FirstName: "Sayed", LastName: "Kamal", Email: "sayed.kamal@example.com", Phone: "+511111111111"},
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "+1234567890"},
		{FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Phone: "+9876543210"},
		{FirstName: "Alice", LastName: "Smith", Email: "alice.smith@example.com", Phone: "+1122334455"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob.johnson@example.com", Phone: "+9988776655"},
		{FirstName: "Emily", LastName: "Wilson", Email: "emily.wilson@example.com", Phone: "+6677889900"},
		{FirstName: "David", LastName: "Miller", Email: "david.miller@example.com", Phone: "+5544332211"},
		{FirstName: "Sophia", LastName: "Clark", Email: "sophia.clark@example.com", Phone: "+6677001122"},
		{FirstName: "Daniel", LastName: "Jones", Email: "daniel.jones@example.com", Phone: "+1122334455"},
		// Add more contacts as needed...
	}
	for _, contact := range contacts {
		store.AddContact(ctx, contact)
	}
}
