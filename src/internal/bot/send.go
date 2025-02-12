package bot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gleisonem/bot-zap-golang-v2/config"
	ServiceAppContext "github.com/gleisonem/bot-zap-golang-v2/contexts"
	"github.com/gleisonem/bot-zap-golang-v2/domains/bot/app"
	domainSend "github.com/gleisonem/bot-zap-golang-v2/domains/bot/send"
	domainBotTypes "github.com/gleisonem/bot-zap-golang-v2/domains/bot/types"
	servicesInternalBot "github.com/gleisonem/bot-zap-golang-v2/internal/bot/services"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/utils"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/whatsapp"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

type serviceSend struct {
	WaCli      *whatsmeow.Client
	appService app.IAppService
}

func NewSendService(waCli *whatsmeow.Client, appService app.IAppService) domainSend.ISendService {
	return &serviceSend{
		WaCli:      waCli,
		appService: appService,
	}
}

func (service serviceSend) SendAudioFunny(ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string) {
	dataWaRecipient, _ := whatsapp.ValidateJidWithLogin(service.WaCli, fromChat)
	dataWaRecipientSender, _ := whatsapp.ValidateJidWithLogin(service.WaCli, sender)

	fmt.Println("dataWaRecipient", dataWaRecipient.String())
	participantJID := sender
	audioDownloaded, errAudioDownloaded := servicesInternalBot.SearchAudioFunnyReturnFile(name)

	if errAudioDownloaded != nil {
		s, err2 := service.WaCli.SendMessage(
			context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
		)
		fmt.Println("send message error download audio in errAudioDownloaded", s, err2)
	}

	audioMimeType := http.DetectContentType(audioDownloaded)

	audioUploaded, err := service.WaCli.Upload(context.Background(), audioDownloaded, whatsmeow.MediaAudio)
	if err != nil {
		fmt.Sprintf("Failed to upload audio: %v", err)
		s, err2 := service.WaCli.SendMessage(
			context.Background(), dataWaRecipient, service.WaCli.BuildReaction(dataWaRecipient, dataWaRecipientSender, stanzaID, "❌"),
		)
		fmt.Println("mandado react", s, err2)
	}

	duration := 2 // 2 minutos

	// Converta a duração para uint32 e pegue o endereço
	seconds := uint32(duration)
	secondsPtr := &seconds
	fmt.Println("dataWaRecipient.String()", dataWaRecipient.String(), "\t", participantJID)

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(audioUploaded.URL),
			DirectPath:    proto.String(audioUploaded.DirectPath),
			Mimetype:      proto.String(audioMimeType),
			FileLength:    proto.Uint64(audioUploaded.FileLength),
			FileSHA256:    audioUploaded.FileSHA256,
			FileEncSHA256: audioUploaded.FileEncSHA256,
			Seconds:       secondsPtr,
			MediaKey:      audioUploaded.MediaKey,
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: &waE2E.Message{
					Conversation: proto.String(messageText),
				},
				StanzaID:    &stanzaID,
				Participant: proto.String(participantJID), // O participante é quem enviou a mensagem original
			},
		},
	}

	s, err2 := service.WaCli.SendMessage(context.Background(), dataWaRecipient, msg)
	fmt.Println("mandado audio funny", s, err2)
}

func (service serviceSend) SendMessage(
	ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageParams domainBotTypes.SendMessageParams,
) {
	dataWaRecipient, _ := whatsapp.ValidateJidWithLogin(service.WaCli, fromChat)
	participantJID := sender

	msg := &waE2E.Message{}
	// if err != nil {
	// 	return nil, err
	// }
	if sendMessageParams.Message != "" && sendMessageParams.IsQuotedMessage {

		println("passei no if mensagem com audio", sendMessageParams.IsQuotedMessage)

		msg = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String(sendMessageParams.Message),
				ContextInfo: &waE2E.ContextInfo{
					StanzaID:    &stanzaID,
					Participant: proto.String(participantJID),
					QuotedMessage: &waE2E.Message{
						AudioMessage: sendMessageParams.AudioMessage,
					},
				},
			},
		}
	}

	if sendMessageParams.Message != "" && !sendMessageParams.IsQuotedMessage {
		println("passei no if mensagem simples", sendMessageParams.IsQuotedMessage)
		msg = &waE2E.Message{Conversation: proto.String(sendMessageParams.Message)}
	}
	// Remove the if statement since the mentionAllUsers field is not defined in the SendMessageParams struct.
	if sendMessageParams.MentionAllUsers {
		groupInfo, _ := service.WaCli.GetGroupInfo(dataWaRecipient)
		fmt.Println(groupInfo.Participants, groupInfo.ParticipantVersionID)

		var mentionedJids []string
		var textMentionedJids string
		for _, participant := range groupInfo.Participants {
			jid := participant.JID // substitua isso pela lógica correta para obter o JID
			fmt.Println(participant.JID, participant.JID.User)
			textMentionedJids += " @" + jid.User
			mentionedJids = append(mentionedJids, jid.String())
		}

		fmt.Println("mandando mentioned", mentionedJids, textMentionedJids)
		msg = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String(textMentionedJids),
				ContextInfo: &waE2E.ContextInfo{
					StanzaID:    &stanzaID,
					Participant: proto.String(participantJID),
					QuotedMessage: &waE2E.Message{
						ExtendedTextMessage: *&sendMessageParams.ExtendedTextMessage,
					},
					MentionedJID: mentionedJids,
				},
			},
		}
		fmt.Println("mandando mentioned")
	}

	if sendMessageParams.MentionAdmin {
		groupInfo, _ := service.WaCli.GetGroupInfo(dataWaRecipient)
		fmt.Println(groupInfo.Participants, groupInfo.ParticipantVersionID)

		var mentionedJids []string
		var textMentionedJids string
		for _, participant := range groupInfo.Participants {

			if participant.IsAdmin {
				jid := participant.JID // substitua isso pela lógica correta para obter o JID
				fmt.Println(participant.JID, participant.JID.User)
				textMentionedJids += " @" + jid.User
				mentionedJids = append(mentionedJids, jid.String())
			}
		}

		fmt.Println("mandando mentioned", mentionedJids, textMentionedJids)
		msg = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String("Supremacia Zaroquiana" + textMentionedJids),
				ContextInfo: &waE2E.ContextInfo{
					// StanzaID:                  request.ReplyMessageID,
					// Participant: proto.String(dataWaRecipient.String()),
					MentionedJID: mentionedJids,
				},
			},
		}
		fmt.Println("mandando mentioned")
	}

	s, err2 := service.WaCli.SendMessage(context.Background(), dataWaRecipient, msg)
	fmt.Println("mandado send", s, err2, msg.String())
}

func (service serviceSend) SendImageSticker(ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageStickerParams domainBotTypes.SendMessageStickerParams) {
	mediaMessage := &sendMessageStickerParams.ImageMessage
	path, err := ServiceAppContext.Context.MessageService.ExtractMedia(context.Background(), config.PathStorages, *mediaMessage)

	go servicesInternalBot.SendReact(ctx, sender, stanzaID, "⌛", fromChat)

	if err != nil {
		log.Errorf("Failed to download image to sticker:", err)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	imageToWebp, errImageToWebp := servicesInternalBot.ConvertJpegToWebp(config.PathStorages, path.MediaPath)

	if errImageToWebp != nil {
		log.Errorf("Failed to convert image png to sticker:", errImageToWebp)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	dataWaImage, errDataWaImage := os.ReadFile(imageToWebp)
	if errDataWaImage != nil {
		fmt.Println("Failed to read image:", errDataWaImage)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	// Upload the image
	uploadedImage, err := service.WaCli.Upload(ctx, dataWaImage, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Failed to upload image:", err)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	quotedMessage := &waE2E.Message{
		ImageMessage: sendMessageStickerParams.ImageMessage,
	}

	servicesInternalBot.SendSticker(ctx, uploadedImage, dataWaImage, fromChat, stanzaID, sender, quotedMessage, false)

	go servicesInternalBot.SendReact(ctx, sender, stanzaID, "✅", fromChat)
}

func (service serviceSend) SendVideoSticker(ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageStickerParams domainBotTypes.SendMessageStickerParams) {
	mediaMessage := &sendMessageStickerParams.VideoMessage
	path, err := ServiceAppContext.Context.MessageService.ExtractMedia(context.Background(), config.PathStorages, *mediaMessage)
	log.Debug("entrei video sticker")
	go servicesInternalBot.SendReact(ctx, sender, stanzaID, "⌛", fromChat)

	if err != nil {
		log.Errorf("Failed to download video to sticker:", err)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	outputPathMp4ToGif := fmt.Sprintf("%s/%d-%s%s", config.PathStorages, time.Now().Unix(), uuid.NewString(), ".gif")
	log.Debug("TIME INIT CONVERT VIDEO to gif")
	errVideoToGif := servicesInternalBot.ConvertMp4ToGif(path.MediaPath, outputPathMp4ToGif)
	log.Debug("TIME FINISH CONVERT VIDEO to gif")
	if errVideoToGif != nil {
		log.Errorf("Failed convert mp4 to gif:", errVideoToGif)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	outputPathGifToWebp := fmt.Sprintf("%s/%d-%s%s", config.PathStorages, time.Now().Unix(), uuid.NewString(), ".webp")
	log.Debug("TIME INIT CONVERT gif to webp")
	errVideoToWebp := servicesInternalBot.ConvertGifToWebp(outputPathMp4ToGif, outputPathGifToWebp)
	log.Debug("TIME FINISH CONVERT gif to webp")
	if errVideoToWebp != nil {
		log.Errorf("Failed convert gif to webp:", errVideoToWebp)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	dataWaImage, errDataWaImage := os.ReadFile(outputPathGifToWebp)
	if errDataWaImage != nil {
		fmt.Println("Failed to read image:", errDataWaImage)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	// Upload the image
	uploadedWebp, err := service.WaCli.Upload(ctx, dataWaImage, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Failed to upload image:", err)
		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

		return
	}

	quotedMessage := &waE2E.Message{
		VideoMessage: sendMessageStickerParams.VideoMessage,
	}

	servicesInternalBot.SendSticker(ctx, uploadedWebp, dataWaImage, fromChat, stanzaID, sender, quotedMessage, true)

	go servicesInternalBot.SendReact(ctx, sender, stanzaID, "✅", fromChat)
}

func (service serviceSend) getMentionFromText(ctx context.Context, messages string) (result []string) {
	mentions := utils.ContainsMention(messages)
	for _, mention := range mentions {
		// Get JID from phone number
		if dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, mention); err == nil {
			result = append(result, dataWaRecipient.String())
		}
	}
	return result
}
