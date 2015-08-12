package main

import (
	"gopkg.in/airbrake/gobrake.v2"
)

const (
	AirBrakeTestAPIKey    = "652ef2cd73e75b8a8efcac293d58ef82"
	AirBrakeTestEnv       = "development"
	AirBrakeProjectID     = 114494
	AirBrakeEndPoint      = "https://api.airbrake.io/notifier_api/v2/notices"
	AirBrakeExpectedClass = "*airbrake.customErr"
)

type AirBrakeMiddleware struct {
	AirBrake *gobrake.Notifier
}

func NewAirBrakeMiddleware() *AirBrakeMiddleware {

	airbrake := gobrake.NewNotifier(AirBrakeProjectID, AirBrakeTestAPIKey)

	airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		notice.Context["environment"] = "production"
		return notice
	})

	return &AirBrakeMiddleware{
		AirBrake: airbrake,
	}
}
