package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	userURL    = "https://localhost:8085/api/v1/users"
	sessionURL = "https://localhost:8085/api/v1/session"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func getClient(clients ...*http.Client) *http.Client {
	if len(clients) == 0 {
		return newClient()
	}
	return clients[0]
}

func doPostRequest[REQ any](url string, reqData REQ, clients ...*http.Client) (int, error) {
	client := getClient(clients...)

	data, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println("error marshalling data")
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error posting data")
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func doDeleteRequest(url string, clients ...*http.Client) (int, error) {
	client := getClient(clients...)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func createAccount(user User, clients ...*http.Client) (int, error) {
	return doPostRequest[User](userURL, user, clients...)
}

func createSession(session Session, clients ...*http.Client) (int, error) {
	return doPostRequest[Session](sessionURL, session, clients...)
}

func deleteAccount(email string, clients ...*http.Client) error {
	deleteURL := userURL + "?email=" + email
	_, err := doDeleteRequest(deleteURL, clients...)
	return err
}

func newSimpleUser(username string) *User {
	return &User{
		Username: username,
		Email:    username + "@gmail.com",
		Password: "123456",
	}
}
