package token

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func AccessTokenForTest() string {
	return "v4.local.1rgk1mWhxlIZTibar2LtKRHAe_Vz-dalZVpa0SQ4EtmPh9XOVOmm5W_QckylWvn-2m9q294_Y0DQ8vu9r8Ete7u9BAKis6Gb2Vvy7GoVWRWSJ4y2kmPYtGRpzmBTOURtP6OaFwVAb6yQDjk4gR7-KNVmxSO32NSck0PTwYkvPYs9zv6ZZhBST31ij7GXrfOUzqI8cbV9ftaCmj_-IjHbLbIj4C-y9v480C8mbl5sLgS2X9in5WIivA"
}

func TestVerify(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		check       func(*Payload, error)
	}{
		{
			name:        "OK",
			accessToken: AccessTokenForTest(),
			check: func(p *Payload, err error) {
				require.NoError(t, err)
				require.Equal(t, p.Issuer, Issuer)
				require.WithinDuration(t, time.Now(), p.ExpiredAt, 50*12*30*24*time.Hour)
				require.WithinDuration(t, time.Now(), p.IssuedAt, 50*12*30*24*time.Hour)
				require.Equal(t, "aaa", p.UserID)
				require.WithinDuration(t, time.Now(), p.NotBefore, 50*12*30*24*time.Hour)
			},
		},
		{
			name:        "token has expired",
			accessToken: "v4.local.XY6G0ykh5-mEfG6TY8r6IN7NoRz1t1hr50dkECT29qdbIg0Uk9Wjvb8rwU9jh_3icefXt6Hwhv5wEnrNZG3Cd5ZC1D6TbvFFZij9JF6mBzjuUOQieDH7PheUM_mR3-4JVnf3l8bH_XrkUwiWGfsxFbjJWgGJn_71TwvteRlFqW0ZI753GmdcrvOlHSAOx0Kxh8L5W0BmPn6kyAV4SQ1V-kKGY65hLeUz-AdXCkSXJz-6p2Ws6GZm5Q",
			check: func(p *Payload, err error) {
				require.Error(t, err)
				require.Equal(t, "parse failed: this token has expired", err.Error())
			},
		},
		{
			name:        "empty user id",
			accessToken: "v4.local.nRd-atiiG8-dlx0nYq32bS5T197oWFczcOLRb9Rlk9pUX-RjAzFCdEVJwDNIMqkhnrhRaLanAlCf5Sq7R5x30N8uCKmexOYB4vp2wNsrad6iTkCM84XmzL2Npc3LjhXUNHZVAjTl1Kyb4NxcQaHZLi8PrODOFesiWYHgkU9RjJnjsYj5Ksx7lq4yqP-DV2Q87_67j2yv5kf3j7AOzEJAH2udPFkU2J1a9esbLWPAQxoRlryWEg",
			check: func(p *Payload, err error) {
				require.Error(t, err)
				require.Equal(t, "cannot empty user id", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v TokenVerifier = NewPasetoVerifier(keyHex)

			p, err := v.Verify(tc.accessToken)
			tc.check(p, err)
		})
	}
}
