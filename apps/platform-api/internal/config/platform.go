package config

type PlatformOptions struct {
	Credentials PlatformCredentialsOptions `envPrefix:"CREDENTIALS_"`
}

type PlatformCredentialsOptions struct {
	AirtableAPI      string `env:"AIRTABLE_API"`
	AnthropicAPI     string `env:"ANTHROPIC_API"`
	ChatgptAPI       string `env:"CHATGPT_API"`
	CoingeckoAPI     string `env:"COINGECKO_API"`
	DeeplAPI         string `env:"DEEPL_API"`
	DeepseekAPI      string `env:"DEEPSEEK_API"`
	DiscordAPI       string `env:"DISCORD_API"`
	ElevenlabsAPI    string `env:"ELEVENLABS_API"`
	FacebookOAuth2   string `env:"FACEBOOK_OAUTH2"`
	GeminiAPI        string `env:"GEMINI_API"`
	GoogleOAuth2     string `env:"GOOGLE_OAUTH2"`
	InstagramOAuth2  string `env:"INSTAGRAM_OAUTH2"`
	LinkedinOAuth2   string `env:"LINKEDIN_OAUTH2"`
	MongodbAPI       string `env:"MONGODB_API"`
	NotionOAuth2     string `env:"NOTION_OAUTH2"`
	PiapiAPI         string `env:"PIAPI_API"`
	PostgresAPI      string `env:"POSTGRES_API"`
	RumbleAPI        string `env:"RUMBLE_API"`
	SendgridAPI      string `env:"SENDGRID_API"`
	SlackOAuth2      string `env:"SLACK_OAUTH2"`
	SoundcloudOAuth2 string `env:"SOUNDCLOUD_OAUTH2"`
	SunomusicAPI     string `env:"SUNOMUSIC_API"`
	TelegramAPI      string `env:"TELEGRAM_API"`
	TiktokOAuth2     string `env:"TIKTOK_OAUTH2"`
	VeoAPI           string `env:"VEO_API"`
	WhatsappAPI      string `env:"WHATSAPP_API"`
	XOAuth2          string `env:"X_OAUTH2"`
}

func (c *PlatformCredentialsOptions) ToCredentialConfigs() map[string]string {
	return map[string]string{
		"airtable_api":      c.AirtableAPI,
		"anthropic_api":     c.AnthropicAPI,
		"chatgpt_api":       c.ChatgptAPI,
		"coingecko_api":     c.CoingeckoAPI,
		"deepl_api":         c.DeeplAPI,
		"deepseek_api":      c.DeepseekAPI,
		"discord_api":       c.DiscordAPI,
		"elevenlabs_api":    c.ElevenlabsAPI,
		"facebook_oauth2":   c.FacebookOAuth2,
		"gemini_api":        c.GeminiAPI,
		"google_oauth2":     c.GoogleOAuth2,
		"instagram_oauth2":  c.InstagramOAuth2,
		"linkedin_oauth2":   c.LinkedinOAuth2,
		"mongodb_api":       c.MongodbAPI,
		"notion_oauth2":     c.NotionOAuth2,
		"piapi_api":         c.PiapiAPI,
		"postgres_api":      c.PostgresAPI,
		"rumble_api":        c.RumbleAPI,
		"sendgrid_api":      c.SendgridAPI,
		"slack_oauth2":      c.SlackOAuth2,
		"soundcloud_oauth2": c.SoundcloudOAuth2,
		"sunomusic_api":     c.SunomusicAPI,
		"telegram_api":      c.TelegramAPI,
		"tiktok_oauth2":     c.TiktokOAuth2,
		"veo_api":           c.VeoAPI,
		"whatsapp_api":      c.WhatsappAPI,
		"x_oauth2":          c.XOAuth2,
	}
}
