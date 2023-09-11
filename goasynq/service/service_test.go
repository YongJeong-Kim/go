package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := newTestService(t)
	err := srv.CreateUser(&CreateUserParam{
		Name: "fefefe",
		After: func(name string) error {
			return nil
		},
	})
	require.NoError(t, err)
}
