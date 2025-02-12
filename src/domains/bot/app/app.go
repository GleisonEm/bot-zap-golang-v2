package app

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

type IAppService interface {
	UploadInWhatsapp(ctx context.Context, dataWaImage []byte) (whatsmeow.UploadResponse, error)
	MakeBuildReaction(ctx context.Context, sender string, stanzaID string, emojiReaction string) *waE2E.Message
	SendInWhatsapp(ctx context.Context, contextLog string, message *waE2E.Message, fromChat string)
	Login(ctx context.Context) (response LoginResponse, err error)
	LoginWithCode(ctx context.Context, phoneNumber string) (loginCode string, err error)
	Logout(ctx context.Context) (err error)
	Reconnect(ctx context.Context) (err error)
	FirstDevice(ctx context.Context) (response DevicesResponse, err error)
	FetchDevices(ctx context.Context) (response []DevicesResponse, err error)
}

type DevicesResponse struct {
	Name   string `json:"name"`
	Device string `json:"device"`
}

type LoginResponse struct {
	ImagePath string        `json:"image_path"`
	Duration  time.Duration `json:"duration"`
	Code      string        `json:"code"`
}
