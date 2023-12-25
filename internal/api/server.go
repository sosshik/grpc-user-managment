package api

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"

	"git.foxminded.ua/foxstudent106264/task-4.1/internal/domain"
	proto "git.foxminded.ua/foxstudent106264/task-4.1/protos/gen/go/user_service"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ServerAPI struct {
	proto.UnimplementedUserServiceServer
	DB domain.DomainInterface
}

func (s *ServerAPI) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	inputUser := req.GetUser()

	if err := CheckPassword(req.GetPassword()); err != nil {
		return &proto.CreateUserResponse{}, fmt.Errorf("CreateUser: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Warnf("CreateUser - unable to generate hash for password: %s", err)
		return &proto.CreateUserResponse{}, fmt.Errorf("CreateUser - unable to generate hash for password: %w", err)
	}

	user := domain.UserProfileDTO{
		OID:       uuid.New(),
		Nickname:  inputUser.Nickname,
		Email:     inputUser.Email,
		FirstName: inputUser.FirstName,
		LastName:  inputUser.LastName,
		Password:  string(hash),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		State:     domain.Active,
	}

	crErr := s.DB.CreateUser(user)
	if crErr != nil {
		log.Warnf("CreateUser: %s", err)
		return &proto.CreateUserResponse{}, fmt.Errorf("CreateUser: %w", err)
	}

	log.Infof("Successfully created user %s", user.Nickname)
	return &proto.CreateUserResponse{
		Oid: &proto.UUID{Value: user.OID.String()},
	}, nil
}

func (s *ServerAPI) GetUserByEmail(ctx context.Context, req *proto.GetUserByEmailRequest) (*proto.GetUserByEmailResponse, error) {
	user, err := s.DB.GetUserByEmail(req.GetEmail())
	if err != nil {
		log.Warnf("GetUserByEmail: %s", err)
		return &proto.GetUserByEmailResponse{}, fmt.Errorf("GetUserByEmail: %w", err)
	}

	return &proto.GetUserByEmailResponse{
		User: &proto.UserInfo{
			Oid:       &proto.UUID{Value: user.OID.String()},
			Nickname:  user.Nickname,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *ServerAPI) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	oid, err := uuid.Parse(req.Oid.GetValue())
	if err != nil {
		log.Warnf("GetUserByID: unable to parse uuid:%s", err)
		return &proto.GetUserByIDResponse{}, fmt.Errorf("GetUserByID: %w", err)
	}
	user, err := s.DB.GetUserByID(oid)
	if err != nil {
		log.Warnf("GetUserByEmail: %s", err)
		return &proto.GetUserByIDResponse{}, fmt.Errorf("GetUserByID: %w", err)
	}

	return &proto.GetUserByIDResponse{
		User: &proto.UserInfo{
			Oid:       &proto.UUID{Value: user.OID.String()},
			Nickname:  user.Nickname,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}
func (s *ServerAPI) GetUsers(ctx context.Context, req *emptypb.Empty) (*proto.GetUsersResponse, error) {

	users, err := s.DB.GetUsers()
	if err != nil {
		log.Warnf("GetUsers:%s", err)
		return &proto.GetUsersResponse{}, fmt.Errorf("GetUsers: %w", err)
	}

	return &proto.GetUsersResponse{
		Users: users,
	}, nil
}
func (s *ServerAPI) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	grpcUser := req.GetUser()
	oid, err := uuid.Parse(grpcUser.Oid.GetValue())
	if err != nil {
		log.Warnf("DeleteUser: unable to parse uuid:%s", err)
		return &proto.UpdateUserResponse{IsOk: false}, fmt.Errorf("DeleteUser: %w", err)
	}
	user := domain.UserProfileDTO{
		OID:       oid,
		Nickname:  grpcUser.GetNickname(),
		Email:     grpcUser.GetEmail(),
		FirstName: grpcUser.GetFirstName(),
		LastName:  grpcUser.GetLastName(),
	}

	err = s.DB.UpdateUser(user)
	if err != nil {
		log.Warnf("UpdateUser:%s", err)
		return &proto.UpdateUserResponse{IsOk: false}, fmt.Errorf("UpdateUser: %w", err)
	}

	return &proto.UpdateUserResponse{IsOk: true}, nil
}
func (s *ServerAPI) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	oid, err := uuid.Parse(req.Oid.GetValue())
	if err != nil {
		log.Warnf("DeleteUser: unable to parse uuid:%s", err)
		return &proto.DeleteUserResponse{IsOk: false}, fmt.Errorf("DeleteUser: %w", err)
	}

	err = s.DB.DeleteUser(oid)
	if err != nil {
		log.Warnf("DeleteUser:%s", err)
		return &proto.DeleteUserResponse{IsOk: false}, fmt.Errorf("DeleteUser: %w", err)
	}

	return &proto.DeleteUserResponse{IsOk: true}, nil
}

func CheckPassword(psw string) error {

	if len(psw) < 8 {

		return errors.New("password is too short, should be at least 8 symbols")

	}

	var lower, upper, number, symbol bool

	for _, letter := range psw {

		if unicode.IsLower(letter) {
			lower = true
		}
		if unicode.IsUpper(letter) {
			upper = true
		}
		if unicode.IsNumber(letter) {
			number = true
		}
		if unicode.IsSymbol(letter) || unicode.IsPunct(letter) {
			symbol = true
		}
	}

	if lower && upper && number && symbol {
		return nil
	}
	return errors.New("wrong password format: password must contatin at least 1 upper case letter, 1 lower case letter, 1 number and 1 symbol")
}
