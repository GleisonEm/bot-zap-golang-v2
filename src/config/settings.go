package config

import (
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
)

var (
	AppVersion             = "v4.17.0"
	AppPort                = "3002"
	AppDebug               = false
	AppOs                  = "AldinoKemal"
	AppPlatform            = waCompanionReg.DeviceProps_PlatformType(1)
	AppBasicAuthCredential string

	PathQrCode    = "statics/qrcode"
	PathSendItems = "statics/senditems"
	PathMedia     = "statics/media"
	PathStorages  = "storages"

	DBName = "whatsapp.db"

	WhatsappAutoReplyMessage    string
	WhatsappWebhook             string = "https://n8n.gemanuel.site/webhook/81a0a43c-571e-44d9-b2f2-257c3b6a403a"
	WhatsappWebhookSecundary    string = "https://n8n.gemanuel.site/webhook-test/81a0a43c-571e-44d9-b2f2-257c3b6a403a"
	WhatsappLogLevel                   = "ERROR"
	WhatsappSettingMaxFileSize  int64  = 50000000  // 50MB
	WhatsappSettingMaxVideoSize int64  = 100000000 // 100MB
	WhatsappTypeUser                   = "@s.whatsapp.net"
	WhatsappTypeGroup                  = "@g.us"
	WhatsappAccountValidation          = true
)
