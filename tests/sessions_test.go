package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	user := User{"CreateTest", "securePass123", "create_test@gmail.com"}
	statusCode, err := createAccount(user)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	defer deleteAccount(user.Email)

	t.Run("missing email", func(t *testing.T) {
		t.Parallel()
		reqWithoutEmail := map[string]string{
			"username": "NoEmailUser",
			"password": "somepass",
		}
		statusCode, err := doPostRequest(userURL, reqWithoutEmail)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("missing password", func(t *testing.T) {
		t.Parallel()
		reqWithoutPassword := map[string]string{
			"username": "NoPassUser",
			"email":    "nopass@example.com",
		}
		statusCode, err := doPostRequest(userURL, reqWithoutPassword)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("missing username", func(t *testing.T) {
		t.Parallel()
		reqWithoutUsername := map[string]string{
			"email":    "nouser@example.com",
			"password": "somepass",
		}
		statusCode, err := doPostRequest(userURL, reqWithoutUsername)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})
}

func TestRequireFields(t *testing.T) {
	t.Parallel()

	user := User{"test", "qwerty", "test@gmail.com"}
	statusCode, err := createAccount(user)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	defer deleteAccount(user.Email)

	t.Run("email required", func(t *testing.T) {
		t.Parallel()
		reqWithoutEmail := map[string]string{
			"password": "qwerty",
		}
		statusCode, err := doPostRequest[map[string]string](sessionURL, reqWithoutEmail)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})

	t.Run("password required", func(t *testing.T) {
		t.Parallel()
		reqWithoutPassword := map[string]string{
			"email": "test@gmail.com",
		}
		statusCode, err := doPostRequest(sessionURL, reqWithoutPassword)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, statusCode)
	})
}

func TestSession(t *testing.T) {
	t.Parallel()

	user := User{"SessionTest", "qwerty1", "SessionTest@gmail.com"}
	statusCode, err := createAccount(user)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	defer deleteAccount(user.Email)

	ttc := []struct {
		name           string
		session        Session
		wantStatusCode int
	}{
		{
			name: "correct request",
			session: Session{
				Email:    "SessionTest@gmail.com",
				Password: "qwerty1",
			},
			wantStatusCode: http.StatusOK,
		},

		{
			name: "incorrect password",
			session: Session{
				Email:    "SessionTest@gmail.com",
				Password: "password",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "incorrect email",
			session: Session{
				Email:    "email@gmail.com",
				Password: "qwerty1",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range ttc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			statusCode, err := createSession(tc.session)
			require.NoError(t, err)
			require.Equal(t, tc.wantStatusCode, statusCode)
		})
	}
}
