package main

import (
	"fmt"
	"jwt_token/pkg/handler"
	"jwt_token/pkg/service"
	//"lqdtests/tests"
	//"testing"
)

func RunTests() {
	fmt.Println("Running tests...")

	//t := &testing.T{}

	/*tests.TestCreateAccount(t)
	tests.TestRequireFields(t)
	tests.TestSession(t)*/

	fmt.Println("Tests completed!")
}

/*
func main() {
	config.Load(context.Background())

	fmt.Println("tests started...")
	RunTests()
}*/

func main() {
	services := service.NewService()
	handlers := handler.NewHandler(services)

	r := handlers.InitRoutes()
	r.Run(":8085")
}
