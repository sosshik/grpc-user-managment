package domain

import (
	"time"

	proto "git.foxminded.ua/foxstudent106264/task-4.1/protos/gen/go/user_service"
	"github.com/google/uuid"
)

type Role int

type State int

const (
	Deleted State = iota - 1
	Banned
	Active
)

type DomainInterface interface {
	CreateUser(user UserProfileDTO) error
	GetUserByID(userID uuid.UUID) (UserProfileDTO, error)
	GetUserByEmail(email string) (UserProfileDTO, error)
	GetUsers() ([]*proto.UserInfo, error)
	UpdateUser(user UserProfileDTO) error
	DeleteUser(oid uuid.UUID) error
}

type UserProfileDTO struct {
	OID       uuid.UUID `json:"oid"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     State     `json:"state"`
}

type GetProfileDTO struct {
	OID       uuid.UUID `json:"oid"`
	Nickname  string    `json:"nickname"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     State     `json:"state"`
}

type Pagination[T any] struct {
	TotalItems  int `json:"total_items"`
	CurrentPage int `json:"current_page"`
	Users       []T `json:"users"`
}
