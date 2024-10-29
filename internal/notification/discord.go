package notification

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ipanardian/price-api/internal/helpers"
	"github.com/ipanardian/price-api/internal/model/frame"
	"github.com/spf13/viper"
)

func SendPriceAlert(id string, message string) (err error) {
	body := frame.DiscordBody{
		Username: "Price Engine Alert",
		Embeds: []frame.DiscordEmbed{{
			Color: 15548997,
			Fields: []frame.DiscordField{
				{
					Name:   "Price Feed ID",
					Value:  id,
					Inline: true,
				},
				{
					Name:   "Message",
					Value:  message,
					Inline: false,
				},
				{
					Name:   "Time",
					Value:  helpers.CurrentTimeAsRFC822(false),
					Inline: false,
				},
			},
		}},
	}

	hookBroker := viper.GetString("DISCORD_MONITORING_URI")

	params, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = http.Post(hookBroker, "application/json", bytes.NewBuffer(params))
	if err != nil {
		return err
	}

	return
}
