package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	todo "github.com/Mobo140/projects/to_do_list"
	"github.com/Mobo140/projects/to_do_list/pkg/service"
	mock_service "github.com/Mobo140/projects/to_do_list/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string    //Имя теста
		inputBody            string    //Тело запроса
		inputUser            todo.User //Структура пользователя передаваемая в метод сервиса.
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK", //Успешная аутентификация
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{ //Невалидное тело запроса
			name:                 "Empty Fields",
			inputBody:            `{"username": "test", "password": "qwerty"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{ //Поведение хендлера при ошибке сервиса
			name:      "Service Failure",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	//Init testTable
	type mockBehavior func(s *mock_service.MockAuthorization, input signInInput)

	testTable := []struct {
		name                 string      //Имя теста
		inputBody            string      //Тело запроса
		inputData            signInInput //Данные для аутентификации
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputData: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("1", nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Empty Fields", //Невалидное тело запроса
			inputBody:            `{"name": "Test", "username": "test"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, input signInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message": "invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputData: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("1", errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message": "service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputData)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test Server
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}

}
