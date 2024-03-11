package service

import (
	"github.com/stretchr/testify/require"
	tkmock "github.com/yongjeong-kim/go/goapigw/account/token/mock"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := tkmock.NewMockTokenMaker(ctrl)
	s := NewAccountService(m)

	testCases := []struct {
		name       string
		buildStubs func()
		check      func()
	}{
		{
			name: "OK",
			buildStubs: func() {
				m.EXPECT().Create(gomock.Eq("aaa"), time.Minute).
					Times(1).
					Return("asdf", nil)
			},
			check: func() {
				tk, err := s.Login("aaa", "1234", time.Minute)
				require.NoError(t, err)
				require.Equal(t, "asdf", tk)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs()
			tc.check()
		})
	}
}
