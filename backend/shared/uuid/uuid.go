package uuid

import (
	"fmt"
	"github.com/google/uuid"
)

// Generate generates a new UUID. Returns def if an error occurs.
func Generate(def string) string {
	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to generate uuid: %w", err))
		return def
	}

	return id.String()
}
