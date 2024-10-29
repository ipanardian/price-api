package frame

type Discord struct {
	Token string
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordEmbed struct {
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Fields      []DiscordField `json:"fields"`
}

type DiscordBody struct {
	Username string         `json:"username"`
	Embeds   []DiscordEmbed `json:"embeds"`
}

func (d *Discord) Publish(mesg map[string]interface{}) (ok bool, err error) {
	return
}
