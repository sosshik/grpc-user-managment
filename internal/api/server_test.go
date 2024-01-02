// server_api_test.go

package api

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/sosshik/grpc-user-managment/internal/mocks"
	proto "github.com/sosshik/grpc-user-managment/protos/gen/go/user_service"
	"github.com/stretchr/testify/mock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestServerAPI_CreateUser(t *testing.T) {
	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}
	req1 := &proto.CreateUserRequest{
		User: &proto.UserInfo{
			Nickname:  "test",
			Email:     "test@example.com",
			FirstName: "test",
			LastName:  "test",
			Oid:       &proto.UUID{Value: ""},
		},
		Password: "Test123.",
	}
	req2 := &proto.CreateUserRequest{
		User: &proto.UserInfo{
			Nickname:  "test",
			Email:     "test@example.com",
			FirstName: "test",
			LastName:  "test",
			Oid:       &proto.UUID{Value: ""},
		},
		Password: "Test12",
	}
	req3 := &proto.CreateUserRequest{
		User: &proto.UserInfo{
			Nickname:  "test",
			Email:     "test@example.com",
			FirstName: "test",
			LastName:  "test",
			Oid:       &proto.UUID{Value: ""},
		},
		Password: "Test123.",
	}

	type args struct {
		ctx context.Context
		req *proto.CreateUserRequest
	}
	tests := []struct {
		name      string
		s         *ServerAPI
		args      args
		needsMock bool
		mockErr   error
		want      *proto.CreateUserResponse
		wantErr   bool
	}{
		{name: "Positive case", s: &s, args: args{ctx: context.Background(), req: req1}, needsMock: true, mockErr: nil, want: &proto.CreateUserResponse{}, wantErr: false},
		{name: "Wrong Pass", s: &s, args: args{ctx: context.Background(), req: req2}, want: &proto.CreateUserResponse{}, wantErr: true},
		{name: "DB Error", s: &s, args: args{ctx: context.Background(), req: req3}, needsMock: true, mockErr: errors.New("error"), want: &proto.CreateUserResponse{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.needsMock {
				mockDB.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockErr).Once()
			}

			_, err := tt.s.CreateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServerAPI_GetUserByEmail(t *testing.T) {
	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}
	req := &proto.GetUserByEmailRequest{
		Email: "test@example.com",
	}
	oid, _ := uuid.Parse("e93b6308-fbc2-40a7-90fc-84627f1580dd")

	type args struct {
		ctx context.Context
		req *proto.GetUserByEmailRequest
	}
	tests := []struct {
		name     string
		s        *ServerAPI
		args     args
		mockResp *proto.UserInfo
		mockErr  error
		want     *proto.GetUserByEmailResponse
		wantErr  bool
	}{
		{
			name: "No Error",
			s:    &s,
			args: args{ctx: context.Background(), req: req},
			mockResp: &proto.UserInfo{
				Oid:       &proto.UUID{Value: oid.String()},
				Nickname:  "test",
				Email:     "test@example.com",
				FirstName: "test",
				LastName:  "test",
			},
			mockErr: nil,
			want: &proto.GetUserByEmailResponse{
				User: &proto.UserInfo{
					Oid:       &proto.UUID{Value: oid.String()},
					Nickname:  "test",
					Email:     "test@example.com",
					FirstName: "test",
					LastName:  "test",
				},
			},
			wantErr: false,
		},
		{
			name:     "Error",
			s:        &s,
			args:     args{ctx: context.Background(), req: req},
			mockResp: &proto.UserInfo{},
			mockErr:  errors.New("error"),
			want:     &proto.GetUserByEmailResponse{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB.On("GetUserByEmail", mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got, err := tt.s.GetUserByEmail(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerAPI.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAPI_GetUserByID(t *testing.T) {
	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}
	oid, _ := uuid.Parse("e93b6308-fbc2-40a7-90fc-84627f1580dd")

	req := &proto.GetUserByIDRequest{
		Oid: &proto.UUID{Value: oid.String()},
	}

	type args struct {
		ctx context.Context
		req *proto.GetUserByIDRequest
	}
	tests := []struct {
		name     string
		s        *ServerAPI
		args     args
		mockResp *proto.UserInfo
		mockErr  error
		want     *proto.GetUserByIDResponse
		wantErr  bool
	}{
		{
			name: "No Error",
			s:    &s,
			args: args{ctx: context.Background(), req: req},
			mockResp: &proto.UserInfo{
				Oid:       &proto.UUID{Value: oid.String()},
				Nickname:  "test",
				Email:     "test@example.com",
				FirstName: "test",
				LastName:  "test",
			},
			mockErr: nil,
			want: &proto.GetUserByIDResponse{
				User: &proto.UserInfo{
					Oid:       &proto.UUID{Value: oid.String()},
					Nickname:  "test",
					Email:     "test@example.com",
					FirstName: "test",
					LastName:  "test",
				},
			},
			wantErr: false,
		},
		{
			name:     "Error",
			s:        &s,
			args:     args{ctx: context.Background(), req: req},
			mockResp: &proto.UserInfo{},
			mockErr:  errors.New("error"),
			want:     &proto.GetUserByIDResponse{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB.On("GetUserByID", mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got, err := tt.s.GetUserByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerAPI.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAPI_GetUsers(t *testing.T) {

	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}

	type args struct {
		ctx context.Context
		req *emptypb.Empty
	}
	tests := []struct {
		name     string
		s        *ServerAPI
		args     args
		mockResp []*proto.UserInfo
		mockErr  error
		want     *proto.GetUsersResponse
		wantErr  bool
	}{
		{
			name: "Positive case",
			s:    &s,
			args: args{ctx: context.Background(), req: &emptypb.Empty{}},
			mockResp: []*proto.UserInfo{
				{
					Oid:       &proto.UUID{Value: "sdgfadfgdfgas"},
					Nickname:  "test",
					Email:     "test@example.com",
					FirstName: "test",
					LastName:  "test",
				},
			},
			mockErr: nil,
			want: &proto.GetUsersResponse{
				Users: []*proto.UserInfo{
					{
						Oid:       &proto.UUID{Value: "sdgfadfgdfgas"},
						Nickname:  "test",
						Email:     "test@example.com",
						FirstName: "test",
						LastName:  "test",
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Negative case",
			s:        &s,
			args:     args{ctx: context.Background(), req: &emptypb.Empty{}},
			mockResp: []*proto.UserInfo{},
			mockErr:  errors.New("error"),
			want:     &proto.GetUsersResponse{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB.On("GetUsers").Return(tt.mockResp, tt.mockErr).Once()

			got, err := tt.s.GetUsers(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerAPI.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAPI_UpdateUser(t *testing.T) {

	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}

	oid := "e93b6308-fbc2-40a7-90fc-84627f1580dd"

	req1 := &proto.UpdateUserRequest{
		User: &proto.UserInfo{
			Oid:       &proto.UUID{Value: oid},
			Nickname:  "test",
			Email:     "test@example.com",
			FirstName: "test",
			LastName:  "test",
		},
	}
	req2 := &proto.UpdateUserRequest{
		User: &proto.UserInfo{
			Oid:       &proto.UUID{Value: oid},
			Nickname:  "test",
			Email:     "test@example.com",
			FirstName: "test",
			LastName:  "test",
		},
	}

	type args struct {
		ctx context.Context
		req *proto.UpdateUserRequest
	}
	tests := []struct {
		name      string
		s         *ServerAPI
		args      args
		needsMock bool
		mockErr   error
		want      *proto.UpdateUserResponse
		wantErr   bool
	}{

		{
			name:      "Positive case",
			s:         &s,
			args:      args{ctx: context.Background(), req: req1},
			needsMock: true,
			mockErr:   nil,
			want: &proto.UpdateUserResponse{
				IsOk: true,
			},
			wantErr: false,
		},
		{
			name:      "DB error case",
			s:         &s,
			args:      args{ctx: context.Background(), req: req2},
			needsMock: true,
			mockErr:   errors.New("error"),
			want: &proto.UpdateUserResponse{
				IsOk: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.needsMock {
				mockDB.On("UpdateUser", mock.Anything).Return(tt.mockErr).Once()
			}

			got, err := tt.s.UpdateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerAPI.UpdateUser() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestServerAPI_DeleteUser(t *testing.T) {

	mockDB := mocks.NewDomainInterface(t)
	s := ServerAPI{DB: mockDB}

	oid := "e93b6308-fbc2-40a7-90fc-84627f1580dd"

	type args struct {
		ctx context.Context
		req *proto.DeleteUserRequest
	}
	tests := []struct {
		name      string
		s         *ServerAPI
		args      args
		needsMock bool
		mockErr   error
		want      *proto.DeleteUserResponse
		wantErr   bool
	}{

		{
			name: "Positive case",
			s:    &s,
			args: args{ctx: context.Background(), req: &proto.DeleteUserRequest{
				Oid: &proto.UUID{Value: oid},
			}},
			needsMock: true,
			mockErr:   nil,
			want: &proto.DeleteUserResponse{
				IsOk: true,
			},
			wantErr: false,
		},
		{
			name: "Wrong OID case",
			s:    &s,
			args: args{ctx: context.Background(), req: &proto.DeleteUserRequest{
				Oid: &proto.UUID{Value: "oid"},
			}},
			needsMock: false,
			want: &proto.DeleteUserResponse{
				IsOk: false,
			},
			wantErr: true,
		},
		{
			name: "DB error case",
			s:    &s,
			args: args{ctx: context.Background(), req: &proto.DeleteUserRequest{
				Oid: &proto.UUID{Value: oid},
			}},
			needsMock: true,
			mockErr:   errors.New("error"),
			want: &proto.DeleteUserResponse{
				IsOk: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.needsMock {
				mockDB.On("DeleteUser", mock.Anything).Return(tt.mockErr).Once()
			}
			got, err := tt.s.DeleteUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerAPI.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServerAPI.DeleteUser() = %v, want %v", got, tt.want)
			}

		})
	}
}
