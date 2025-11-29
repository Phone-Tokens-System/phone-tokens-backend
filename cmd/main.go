package main

import (
	"fmt"
	"phone-tokens/internal/app/config/env"
	"phone-tokens/internal/sms_service/service"
)

func main() {
	config, err := env.LoadConfigEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", config)
	aeroService := service.NewAeroService(config.Email, config.ApiKey)

	sms, err := aeroService.GetSmsList()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sms)
	fmt.Println(sms[0])
}
