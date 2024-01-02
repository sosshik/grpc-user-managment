package domain

import (
	"github.com/google/uuid"
	proto "github.com/sosshik/grpc-user-managment/protos/gen/go/user_service"
)

type Role int

type State int

const (
	Deleted State = iota - 1
	Banned
	Active
)

type DomainInterface interface {
	CreateUser(user *proto.UserInfo, pass string, state State) error
	GetUserByID(oid uuid.UUID) (*proto.UserInfo, error)
	GetUserByEmail(email string) (*proto.UserInfo, error)
	GetUsers() ([]*proto.UserInfo, error)
	UpdateUser(user *proto.UserInfo) error
	DeleteUser(oid uuid.UUID) error
}
