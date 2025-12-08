package main

import (
	"fmt"
	"phone-tokens/internal/app/config/env"
	"phone-tokens/internal/sms_service/service/sms_aero"
)

func main() {
	config, err := env.LoadConfigEnv()
	if err != nil {
		panic(err)
	}
	aeroService := sms_aero.NewAeroService(config.Email, config.ApiKey)

	sms, err := aeroService.GetSmsList()
	if err != nil {
		panic(err)
	}
	for key, value := range sms {
		fmt.Println(key, value)
	}
	//smsSent, err := aeroService.SendSms(79232213113, "Phone tokens system. defense date: December 18")
	//if err != nil {
	//	fmt.Println(err)
	//}
}
