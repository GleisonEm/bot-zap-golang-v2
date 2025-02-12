package bot

import (
	"context"
	"fmt"
	"mime"
	"os"
	"time"

	"github.com/gleisonem/bot-zap-golang-v2/config"
	ServiceAppContext "github.com/gleisonem/bot-zap-golang-v2/contexts"
	domainMessage "github.com/gleisonem/bot-zap-golang-v2/domains/bot/message"
	domainBotTypes "github.com/gleisonem/bot-zap-golang-v2/domains/bot/types"
	servicesInternalBot "github.com/gleisonem/bot-zap-golang-v2/internal/bot/services"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/utils"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/whatsapp"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

type serviceMessage struct {
	WaCli *whatsmeow.Client
}

func NewMessageService(waCli *whatsmeow.Client) domainMessage.IMessageService {
	return &serviceMessage{
		WaCli: waCli,
	}
}

func (service serviceMessage) ConvertMessageAudioToText(
	ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageParams domainBotTypes.SendMessageParams,
) {
	dataWaRecipient, _ := whatsapp.ValidateJidWithLogin(service.WaCli, fromChat)
	dataWaRecipientSender, _ := whatsapp.ValidateJidWithLogin(service.WaCli, sender)
	audioMessage := &sendMessageParams.AudioMessage
	participantJID := sender

	reactLoading, errReactLoading := service.WaCli.SendMessage(context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "⏳"))
	fmt.Println("mandado react", reactLoading, errReactLoading)

	path, err := ServiceAppContext.Context.MessageService.ExtractMedia(context.Background(), config.PathStorages, *audioMessage)
	if err != nil {
		log.Errorf("Failed to download audio: %v", err)
		s, err2 := service.WaCli.SendMessage(
			context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
		)
		fmt.Println("mandado react", s, err2)
	} else {
		log.Infof("audio downloaded to %s", path)
		filePathAudio := fmt.Sprintf("%s/%d-%s%s", config.PathStorages, time.Now().Unix(), uuid.NewString(), ".wav")
		errConvertOgaToWav := utils.ConvertOgaToWav(path.MediaPath, filePathAudio)

		if errConvertOgaToWav != nil {
			s, err2 := service.WaCli.SendMessage(
				context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
			)
			fmt.Println("mandado react", s, err2)
			log.Errorf("Failed to convert audio: %v", errConvertOgaToWav)
		}
		response, errSendToRecogntionApiAudioFile := servicesInternalBot.SendToRecogntionApiAudioFile(filePathAudio)

		if errSendToRecogntionApiAudioFile != nil {
			s, err2 := service.WaCli.SendMessage(
				context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
			)
			fmt.Println("mandado react", s, err2)
			log.Errorf("Failed to convert audio to text: %v", errSendToRecogntionApiAudioFile)
		}

		if response.Text == "" {
			log.Errorf("Failed to convert audio to text: %v", "empty response")

			s, err2 := service.WaCli.SendMessage(
				context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
			)
			fmt.Println("mandado react", s, err2)
		}

		go service.WaCli.SendMessage(context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "✅"))

		msg := &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String(response.Text),
				ContextInfo: &waE2E.ContextInfo{
					StanzaID:    &stanzaID,
					Participant: proto.String(participantJID),
					QuotedMessage: &waE2E.Message{
						AudioMessage: sendMessageParams.AudioMessage,
					},
				},
			},
		}

		_, errorSender := service.WaCli.SendMessage(context.Background(), dataWaRecipient, msg)
		fmt.Println("mandado convert audio to text", errorSender)
	}
}

func (service serviceMessage) ExtractMedia(ctx context.Context, storageLocation string, mediaFile whatsmeow.DownloadableMessage) (extractedMedia domainMessage.ExtractedMessageMedia, err error) {
	if mediaFile == nil {
		logrus.Info("Skip download because data is nil")
		return extractedMedia, nil
	}

	data, err := service.WaCli.Download(mediaFile)
	if err != nil {
		logrus.Info("err download", err)
		return extractedMedia, err
	}

	switch media := mediaFile.(type) {
	case *waE2E.ImageMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	case *waE2E.AudioMessage:
		extractedMedia.MimeType = media.GetMimetype()
	case *waE2E.VideoMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	case *waE2E.StickerMessage:
		extractedMedia.MimeType = media.GetMimetype()
	case *waE2E.DocumentMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	}

	extensions, _ := mime.ExtensionsByType(extractedMedia.MimeType)
	extractedMedia.MediaPath = fmt.Sprintf("%s/%d-%s%s", storageLocation, time.Now().Unix(), uuid.NewString(), extensions[0])
	err = os.WriteFile(extractedMedia.MediaPath, data, 0600)
	if err != nil {
		logrus.Info("err write", err)
		return extractedMedia, err
	}
	return extractedMedia, nil
}
