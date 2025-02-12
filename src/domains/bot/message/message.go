package message

import (
	"context"

	domainBotTypes "github.com/gleisonem/bot-zap-golang-v2/domains/bot/types"
	"go.mau.fi/whatsmeow"
)

type IMessageService interface {
	ConvertMessageAudioToText(ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageParams domainBotTypes.SendMessageParams)
	ExtractMedia(ctx context.Context, storageLocation string, mediaFile whatsmeow.DownloadableMessage) (extractedMedia ExtractedMessageMedia, err error)
}

type GenericResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type RevokeRequest struct {
	MessageID string `json:"message_id" uri:"message_id"`
	Phone     string `json:"phone" form:"phone"`
}

type DeleteRequest struct {
	MessageID string `json:"message_id" uri:"message_id"`
	Phone     string `json:"phone" form:"phone"`
}

type ReactionRequest struct {
	MessageID string `json:"message_id" form:"message_id"`
	Phone     string `json:"phone" form:"phone"`
	Emoji     string `json:"emoji" form:"emoji"`
}

type UpdateMessageRequest struct {
	MessageID string `json:"message_id" uri:"message_id"`
	Message   string `json:"message" form:"message"`
	Phone     string `json:"phone" form:"phone"`
}
