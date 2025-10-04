package mocks

import "github.com/mahmoud-shabban/snippetbox/internal/models"

type UserModel struct{}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dup@mock.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "test@mock.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (u *UserModel) Exits(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
