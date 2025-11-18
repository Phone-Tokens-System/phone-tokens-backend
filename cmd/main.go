package main

import "phone-tokens/internal/service"

func main() {
	err := service.CreateOurCert()
	if err != nil {
		panic(err)
	}
}
