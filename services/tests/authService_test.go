package services

import (
	"testing"

	"github.com/Renan-Parise/codium/entities"
	"github.com/Renan-Parise/codium/errors"
	"github.com/Renan-Parise/codium/services"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[string]entities.User
}

func (m *mockUserRepository) FindByID(id int) (*entities.User, error) {
	panic("unimplemented")
}

func (m *mockUserRepository) DeactivateUser(userID int) error {
	panic("unimplemented")
}

func (m *mockUserRepository) DeleteInactiveUsers() error {
	panic("unimplemented")
}

func (m *mockUserRepository) FindByEmail(email string) (*entities.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.NewQueryError("user not found")
	}
	return &user, nil
}

func (m *mockUserRepository) Create(user entities.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.NewQueryError("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepository) Update(user entities.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return errors.NewQueryError("user not found")
	}
	m.users[user.Username] = user
	return nil
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	err = service.Register(user)
	assert.NotNil(t, err)
}

func TestLogin(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	token, err := service.Login("testuser", "password123")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	_, err = service.Login("testuser", "wrongpassword")
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	user.Password = "newpassword123"
	err = service.Update(user)
	assert.Nil(t, err)

	_, err = service.Login("testuser", "password123")
	assert.NotNil(t, err)

	token, err := service.Login("testuser", "newpassword123")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}