package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
)

func Register(oauth2RedirectURL string) {
	domain.RegisterCredential(NewAirtableAPICredential())
	domain.RegisterCredential(NewAnthropicAPICredential())
	domain.RegisterCredential(NewChatgptAPICredential())
	domain.RegisterCredential(NewCoingeckoAPICredential())
	domain.RegisterCredential(NewDeeplAPICredential())
	domain.RegisterCredential(NewDeepseekAPICredential())
	domain.RegisterCredential(NewDiscordAPICredential())
	domain.RegisterCredential(NewElevenlabsAPICredential())
	domain.RegisterCredential(NewFacebookOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewGeminiAPICredential())
	domain.RegisterCredential(NewGmailOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewGoogleDocsOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewGoogleDriveOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewGoogleKeepOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewGoogleSheetsOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewInstagramOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewLinkedinOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewNotionOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewPiapiAPICredential())
	domain.RegisterCredential(NewSendgridAPICredential())
	domain.RegisterCredential(NewSlackOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewSoundcloudOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewSunomusicAPICredential())
	domain.RegisterCredential(NewTelegramAPICredential())
	domain.RegisterCredential(NewTiktokOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewVeoAPICredential())
	domain.RegisterCredential(NewWhatsappAPICredential())
	domain.RegisterCredential(NewXOAuth2Credential(oauth2RedirectURL))
	domain.RegisterCredential(NewYoutubeOAuth2Credential(oauth2RedirectURL))
}
