package user

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtSecretTestKey = "thisIsASecretSoDontTellAnyone"

func TestGenerateJWT(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", jwtSecretTestKey)
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		expectedUserID := "abcdefghijklmnopqrstuvwxyz"

		uc := UserClaims{
			UserID: "abcdefghijklmnopqrstuvwxyz",
			MapClaims: jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			},
		}

		// call to test generateJWT
		signedJWT, err := generateJWT(uc)
		require.NoError(err)

		// Parse the signedJWT to confirm its contents
		var actualUserClaims UserClaims
		actualToken, err := jwt.ParseWithClaims(
			signedJWT,
			&actualUserClaims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretTestKey), nil
			},
		)
		require.NoError(err)

		// The actual token parsed will not be nil even if it fails, it will still contain the data from the signed jwt.
		// The same occurs for the actual claims returned, it will also contain the claims/payload data stored in the jwt.
		// Thus, it's pointless for me to assert NotEmpty or NotNil
		assert.True(actualToken.Valid)
		assert.Equal(expectedUserID, actualUserClaims.UserID)
		assert.IsType(&UserClaims{}, actualToken.Claims)
	})
}

func TestVerifyUserJWT(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", jwtSecretTestKey)
	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		expirationDate := time.Now().Add(time.Hour * 24).Unix()
		expectedUserID := "abcdefghijklmnopqrstuvwxyz"
		uc := UserClaims{
			UserID: "abcdefghijklmnopqrstuvwxyz",
			MapClaims: jwt.MapClaims{
				"exp": expirationDate,
			},
		}

		signedJWT, err := generateJWT(uc)
		require.NoError(err)

		var actualUserClaims UserClaims
		err = verifyUserJWT(signedJWT, &actualUserClaims)
		require.NoError(err)

		assert.Equal(expectedUserID, actualUserClaims.UserID)
		exp, ok := actualUserClaims.MapClaims["exp"].(time.Time)
		if ok {
			assert.Equal(expirationDate, exp)
		}
	})

	t.Run("FailTamperedJWT", func(t *testing.T) {
		require := require.New(t)

		uc := UserClaims{
			UserID: "abcdefghijklmnopqrstuvwxyz",
			MapClaims: jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			},
		}

		signedJWT, err := generateJWT(uc)
		require.NoError(err)

		tamperedJWT := signedJWT + "a"

		err = verifyUserJWT(tamperedJWT, &UserClaims{})
		require.Error(err)
	})
}
