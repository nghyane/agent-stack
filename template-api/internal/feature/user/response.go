package user

import (
	"time"

	"github.com/nghiahoang/template-api/internal/feature/user/data"
)

// UserResponse is the public representation of a user. It never includes the
// password hash. Map DB rows into this explicit struct; never return DB models
// directly.
//
// Contract types are feature-qualified (UserResponse, not DTO) because the
// OpenAPI schema namespace is flat: a generic name like "DTO" would collide
// with other features' response types.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// responseFromCreate maps a CreateUser row to the public response.
func responseFromCreate(row data.CreateUserRow) UserResponse {
	return UserResponse{
		ID:        row.ID.String(),
		Email:     row.Email,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// responseFromUser maps a full user row (login lookup) to the public response.
func responseFromUser(row data.User) UserResponse {
	return UserResponse{
		ID:        row.ID.String(),
		Email:     row.Email,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// responseFromByID maps a GetUserByID row to the public response.
func responseFromByID(row data.GetUserByIDRow) UserResponse {
	return UserResponse{
		ID:        row.ID.String(),
		Email:     row.Email,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
