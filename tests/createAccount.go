package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"jwt_token/config"
	"log"
	"net/http"
)

func Test() {
	_, err := config.GetConfig()
	if err != nil {
		log.Fatal("config is not initialized!")
		return
	}

	httpposturl := "http://localhost:8080/api/v1/users"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"email": "dsasl@gmail.com",
		"password": "dsasl123",
		"username": "dsasl"
	}`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

}
