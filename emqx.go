package ucodesdk

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (u *object) ConnectToEMQX() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(u.config.MQTTBroker)
	opts.SetUsername(u.config.MQTTUsername) // Set your username
	opts.SetPassword(u.config.MQTTPassword) // Set your password

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
