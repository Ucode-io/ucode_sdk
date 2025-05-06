package ucodesdk

import (
	"time"
)

const BaseURL = "https://api.client.u-code.io"

type Config struct {
	AppId string
	// BaseURL        string
	FunctionName   string
	ProjectId      string
	RequestTimeout time.Duration
	BaseAuthUrl    string
	MQTTBroker     string
	MQTTUsername   string
	MQTTPassword   string
}
