package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ipanardian/price-api/internal/logger"
	"github.com/spf13/viper"
)

type Message struct {
	Level   Level
	Color   Color
	Title   string
	Message string
	Fields  []Fields
}

type Fields struct {
	Key   string
	Value interface{}
}

type Level string

const (
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
)

type Color int

const (
	ColorRed    Color = 0xFF0000
	ColorGreen  Color = 0x00FF00
	ColorYellow Color = 0xFFFF00
)

func Send(msg Message) {
	webhookURL := viper.GetString("DISCORD_MONITORING_URI")
	embed := map[string]interface{}{
		"title":       msg.Title,
		"description": msg.Message,
		"color":       int(msg.Color),
		"fields":      []map[string]interface{}{},
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	for _, field := range msg.Fields {
		value := fmt.Sprintf("%v", field.Value)
		embed["fields"] = append(embed["fields"].([]map[string]interface{}), map[string]interface{}{
			"name":   field.Key,
			"value":  value,
			"inline": false,
		})
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Sugar().Errorf("failed to marshal payload: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Log.Sugar().Errorf("failed to send HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		logger.Log.Sugar().Errorf("received non-200 status code: %d body: %v", resp.StatusCode, string(body))
		return
	}

	logger.Log.Sugar().Debugln("Successfully sent notification")
}
