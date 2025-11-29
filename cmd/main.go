package main

import (
	"phone-tokens/internal/certificates/service"
)

func main() {
	err := service.CreateOurCert()
	if err != nil {
		panic(err)
	}
}
