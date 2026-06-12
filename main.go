package main

import (
	"flag"
	"fmt"
	"os"

	routes "trading-office/trading_office_backend/route"
)

// @title Trading Office AI Dashboard API
// @version 0.1.0
// @description Trading Office AI Dashboard API documentation
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	env := flag.String("env", "dev", "environment: dev, qa, uat, prod")
	flag.Parse()

	os.Setenv("APP_ENV", *env)
	fmt.Printf("[main] environment: %s\n", *env)

	app, wg, err := routes.Bootstrap()
	if err != nil {
		panic(err)
	}

	go app.Start()
	wg.Wait()
}
