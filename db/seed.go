package db

import (
	"app/domain"
	"context"
	"fmt"
)

func (store *InMemoryContactStore) SeedMe(ctx context.Context) {
	contacts := make([]domain.Contact, 0)

	for i := 1; i <= 100; i++ {
		contact := domain.Contact{
			FirstName: fmt.Sprintf("User%d", i),
			LastName:  fmt.Sprintf("Lastname%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Phone:     fmt.Sprintf("+123456789%d", i),
		}
		contacts = append(contacts, contact)
	}

	for _, contact := range contacts {
		store.AddContact(ctx, contact)
	}
}
