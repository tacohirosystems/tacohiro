package discord

type (
	// https://docs.discord.com/developers/resources/user#user-object
	User struct {
		ID            UserID `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		GlobalName    string `json:"global_name"`
		Bot           bool   `json:"bot"`
		System        bool   `json:"system"`
	}

	UserID string
)
