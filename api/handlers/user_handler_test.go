package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"cow_sso/api/dto/request"
	"cow_sso/api/dto/response"
	"cow_sso/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserHandler struct {
	userService *mocks.IUserService
}

type userMocks struct {
	userHandler func(f *mockUserHandler)
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		mocks       userMocks
		name        string
		expNickName int
	}{
		{
			name: "error get users",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetAll", mock.Anything).Return([]response.UserResponse{}, errors.New("error x"))
				},
			},
			expNickName: http.StatusInternalServerError,
		},
		{
			name: "full flow",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetAll", mock.Anything).Return([]response.UserResponse{}, nil)
				},
			},
			expNickName: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.GET(url, func(ctx *gin.Context) {
				handler.GetAll(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expNickName, res.Code)
		})
	}
}

func Test_GetByNickName(t *testing.T) {
	tests := []struct {
		mocks       userMocks
		name        string
		nickName    string
		expNickName int
	}{
		{
			name:     "code not sending",
			nickName: "",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expNickName: http.StatusBadRequest,
		},
		{
			name:     "error getting user by id",
			nickName: "ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetByNickName", mock.Anything, "ABC").Return(response.UserResponse{}, errors.New("x"))
				},
			},
			expNickName: http.StatusInternalServerError,
		},
		{
			name:     "full flow",
			nickName: "ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("GetByNickName", mock.Anything, "ABC").Return(response.UserResponse{
						NickName: "ABC",
					}, nil)
				},
			},
			expNickName: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.GET(url, func(ctx *gin.Context) {
				if tc.nickName != "" {
					ctx.AddParam("code", tc.nickName)
				}
				handler.GetByNickName(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expNickName, res.Code)
		})
	}
}

func Test_Create(t *testing.T) {
	tests := []struct {
		input   interface{}
		mocks   userMocks
		name    string
		expCode int
	}{
		{
			name:  "error on input",
			input: "ABC",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "error creating user",
			input: request.UserRequest{
				Name:     "a",
				NickName: "c",
				Email:    "d",
				LastName: "e",
			},
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Create", mock.Anything, request.UserRequest{
						Name:     "a",
						NickName: "c",
						Email:    "d",
						LastName: "e",
					}).Return(errors.New("error creating user"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name: "full flow",
			input: request.UserRequest{
				Name:     "a",
				NickName: "c",
				Email:    "d",
				LastName: "e",
			},
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Create", mock.Anything, request.UserRequest{
						Name:     "a",
						NickName: "c",
						Email:    "d",
						LastName: "e",
					}).Return(nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			url := "/users/"
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			engine.POST(url, func(ctx *gin.Context) {
				handler.Create(ctx)
			})
			res := httptest.NewRecorder()
			b, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewBuffer(b)))
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		mocks   userMocks
		name    string
		userID  string
		expCode int
	}{
		{
			name: "nick name isnt present",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {},
			},
			expCode: http.StatusBadRequest,
		},
		{
			name:   "nick name not found",
			userID: "abc",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Delete", mock.Anything, "abc").Return("", errors.New("nick name not found"))
				},
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "full flow",
			userID: "abc",
			mocks: userMocks{
				userHandler: func(f *mockUserHandler) {
					f.userService.Mock.On("Delete", mock.Anything, "abc").Return("test", nil)
				},
			},
			expCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockUserHandler{
				&mocks.IUserService{},
			}
			tc.mocks.userHandler(ms)
			handler := NewUserHandler(ms.userService)
			_, engine := gin.CreateTestContext(httptest.NewRecorder())
			url := "/users"
			engine.DELETE(url, func(ctx *gin.Context) {
				if tc.userID != "" {
					ctx.AddParam("code", tc.userID)
				}
				handler.Delete(ctx)
			})
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			engine.ServeHTTP(res, req)
			assert.Equal(t, tc.expCode, res.Code)
		})
	}
}
