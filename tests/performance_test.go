package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

func usersGenerator() func() User {
	login := "user-number-"
	userNumber := 0
	return func() User {
		userNumber++
		username := login + strconv.Itoa(userNumber)
		email := username + "@gmail.com"
		return User{
			Username: username,
			Email:    email,
			Password: "123456",
		}
	}
}

func generateUsers(number int) <-chan User {
	usersCh := make(chan User)
	nextUser := usersGenerator()
	go func() {
		defer close(usersCh)
		for i := 0; i < number; i++ {
			usersCh <- nextUser()
		}
	}()
	return usersCh
}

func TestCreateManyUsers(t *testing.T) {

	usersCh := generateUsers(1000)
	workers := 100

	var isError bool
	var once sync.Once

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := newClient()
			for {
				user, ok := <-usersCh
				if !ok {
					return
				}
				statusCode, err := createAccount(user, client)

				if err != nil || statusCode != http.StatusOK {
					fmt.Println(fmt.Sprintf("Name: %v: StatusCode: %d Error: %v", user, statusCode, err))
					once.Do(func() {
						isError = true
					})
				}
			}

		}()
	}
	wg.Wait()

	require.False(t, isError)
}

func TestCreateUserPerformance(t *testing.T) {

	user := newSimpleUser("TestCreateUserPerformance")

	start := time.Now()
	statusCode, err := createAccount(*user)
	dur := time.Since(start)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
	require.Less(t, dur, 200*time.Millisecond)
}

func TestCreateSessionPerformance(t *testing.T) {

	user := newSimpleUser("TestCreateSessionPerformance")
	client := newClient()

	statusCode, err := createAccount(*user, client)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)

	session := Session{Email: user.Email, Password: user.Password}

	start := time.Now()

	statusCode, err = createSession(session, client)

	dur := time.Since(start)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
	require.Less(t, dur, 100*time.Millisecond)
}
