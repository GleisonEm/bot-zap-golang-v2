package structs

import (
	"go.mau.fi/whatsmeow/proto/waE2E"
)

type SendMessageParams struct {
	MentionAllUsers     bool
	Message             string
	MentionUsers        []string
	IsQuotedMessage     bool
	AudioMessage        *waE2E.AudioMessage
	TextMessage         *waE2E.MessageContextInfo
	ExtendedTextMessage *waE2E.ExtendedTextMessage
	IsQuotedType        string
	MentionAdmin        bool
}

type SendMessageStickerParams struct {
	MentionAllUsers bool
	Message         string
	MentionUsers    []string
	IsQuotedMessage bool
	ImageMessage    *waE2E.ImageMessage
	VideoMessage    *waE2E.VideoMessage
	TypeMedia       string
}
