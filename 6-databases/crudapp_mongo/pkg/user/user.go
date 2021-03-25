package user

import "errors"

type User struct {
	ID       uint32
	Login    string
	password string
}

type UserRepo struct {
	data map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*User{
			"rvasily": &User{
				ID:       1,
				Login:    "rvasily",
				password: "love",
			},
		},
	}
}

var (
	ErrNoUser  = errors.New("No user found")
	ErrBadPass = errors.New("Invald password")
)

func (repo *UserRepo) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	// dont do this un production :)
	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}
