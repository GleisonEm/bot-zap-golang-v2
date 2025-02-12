package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ServiceAppContext "github.com/gleisonem/bot-zap-golang-v2/contexts"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func SendSticker(ctx context.Context, uploadedWebp whatsmeow.UploadResponse, dataWaImage []byte, fromChat string, stanzaID string, participantJID string, QuotedMessage *waE2E.Message, isAnimated bool) {
	// // Upload the image
	// fmt.Println("mimetype", http.DetectContentType(dataWaImage))
	// // Create the sticker message
	msg := &waE2E.Message{StickerMessage: &waE2E.StickerMessage{
		URL:               proto.String(uploadedWebp.URL),
		DirectPath:        proto.String(uploadedWebp.DirectPath),
		MediaKey:          uploadedWebp.MediaKey,
		Mimetype:          proto.String(http.DetectContentType(dataWaImage)),
		FileEncSHA256:     uploadedWebp.FileEncSHA256,
		FileSHA256:        uploadedWebp.FileSHA256,
		FileLength:        proto.Uint64(uploadedWebp.FileLength),
		MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
		IsAnimated:        proto.Bool(isAnimated),
		StickerSentTS:     proto.Int64(time.Now().Unix()),
		IsLottie:          proto.Bool(false),
		// ContextInfo: &waE2E.ContextInfo{
		// 	StanzaID:      &stanzaID,
		// 	Participant:   proto.String(participantJID),
		// 	QuotedMessage: QuotedMessage,
		// },
	}}
	fmt.Println("figurinha ", msg)
	ServiceAppContext.Context.AppService.SendInWhatsapp(context.Background(), "mandado sticker", msg, fromChat)
}
