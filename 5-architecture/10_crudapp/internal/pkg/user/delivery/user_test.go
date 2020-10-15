package delivery

import (
	"crudapp/internal/pkg/models"
	"crudapp/internal/pkg/user/mock"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestGetByID(t *testing.T) {
	login := "somelogin"
	retModel := &models.User{
		ID:       0,
		Login:    login,
		Password: "password",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mock.NewMockRepository(ctrl)
	mock.EXPECT().GetByLogin(login).Times(1).Return(retModel, nil)

	uh := UserHandler{
		UserRepo: mock,
	}

	models, err := uh.GetUserByID(login)

	require.NoError(t, err)
	require.Equal(t, models, retModel)
}
