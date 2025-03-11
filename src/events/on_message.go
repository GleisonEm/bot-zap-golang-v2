package events

import (
	"context"
	"fmt"

	"strings"

	"github.com/gleisonem/bot-zap-golang-v2/config"
	ServiceAppContext "github.com/gleisonem/bot-zap-golang-v2/contexts"
	domainBotTypes "github.com/gleisonem/bot-zap-golang-v2/domains/bot/types"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"go.mau.fi/whatsmeow/types/events"
)

type ExtractedMedia struct {
	MediaPath string `json:"media_path"`
	MimeType  string `json:"mime_type"`
	Caption   string `json:"caption"`
}

func OnMessage(evt *events.Message) {
	fmt.Println("RawMessage", evt.RawMessage)
	return
	messageText := ""
	if evt.Message.GetExtendedTextMessage().GetText() != "" {
		messageText = evt.Message.GetExtendedTextMessage().GetText()
	} else if evt.Message.GetConversation() != "" {
		messageText = evt.Message.GetConversation()
	}
	// fmt.Println("message text ", messageText, evt.Message)
	argument := utils.GetArgument(messageText)
	command := utils.ProcessCommand(messageText)
	sender := evt.Info.Sender.User + "@" + evt.Info.Sender.Server
	fromChat := evt.Info.Chat.String()
	stanzaID := evt.Info.ID

	// if fromChat == "558796485300-1461896371@g.us" {
	// 	return
	// }
	// fmt.Println("Received message ", string(evt.Info.ID), evt.Info.SourceString(), "is group:", evt.Info.IsGroup, evt.Message)

	downloadMedia := false

	if downloadMedia {
		img := evt.Message.GetImageMessage()
		if img != nil {
			path, err := ServiceAppContext.Context.MessageService.ExtractMedia(context.Background(), config.PathStorages, img)
			if err != nil {
				log.Errorf("Failed to download image: %v", err)
			} else {
				log.Infof("Image downloaded to %s", path)
			}
		} else {
			fmt.Println("img n é ", img)
		}

		video := evt.Message.GetVideoMessage()
		if video != nil {
			path, err := ServiceAppContext.Context.MessageService.ExtractMedia(context.Background(), config.PathStorages, video)
			if err != nil {
				log.Errorf("Failed to download image: %v", err)
			} else {
				log.Infof("Image downloaded to %s", path)
			}
		} else {
			fmt.Println("video n é ", video)
		}
	}
	// , "is user", evt.Info.Chat.IsUser(), "is broadcast", evt.Info.Chat.IsBroadcast(), "is server", evt.Info.Chat.IsServer(), "is status", evt.Info.Chat.IsStatus(), "is group", evt.Info.Chat.IsGroup(), "is user", evt.Info.Chat.IsUser(), "is broadcast", evt.Info.Chat.IsBroadcast(), "is server", evt.Info.Chat.IsServer(), "is status", evt.Info.Chat.IsStatus()
	// fmt.Println(argument, "\t", sender, "\t", evt.Info.Sender.Server, "\t", evt.Info.Chat)

	// fmt.Println("command", command, "argument", argument, "sender", sender, "fromChat", fromChat, "stanzaID", stanzaID, "messageText", messageText, evt.Message.GetConversation(), evt.Message.String())

	// if !strings.Contains(config.ChatsDevEnabled, fromChat) {
	// 	return
	// }
	// fmt.Println("COMANDO", command)
	imgMessage := evt.RawMessage.GetImageMessage()
	videoMessage := evt.RawMessage.GetVideoMessage()

	// if command == "!stickervideo" {
	// 	ctx := context.Background()
	// 	// filewebp := "/1725781806-2f93b821-a573-4423-9c58-b6fc39cf78a1-ezgif.com-resize2.webp"
	// 	filewebp := "/output7.webp"
	// 	dataWaImage, errDataWaImage := os.ReadFile(config.PathStorages + filewebp)
	// 	if errDataWaImage != nil {
	// 		fmt.Println("Failed to read image:", errDataWaImage)
	// 		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

	// 		return
	// 	}

	// 	// Upload the image
	// 	uploadedWebp, err := ServiceAppContext.Context.AppService.UploadInWhatsapp(ctx, dataWaImage)
	// 	if err != nil {
	// 		fmt.Println("Failed to upload image:", err)
	// 		go servicesInternalBot.SendReact(ctx, sender, stanzaID, "❌", fromChat)

	// 		return
	// 	}

	// 	quotedMessage := &waE2E.Message{}

	// 	servicesInternalBot.SendSticker(ctx, uploadedWebp, dataWaImage, fromChat, stanzaID, sender, quotedMessage, true)
	// }

	if imgMessage != nil {
		if imgMessage.GetCaption() == "!sticker" {
			go ServiceAppContext.Context.SendService.SendImageSticker(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageStickerParams{
				ImageMessage: imgMessage,
			})
		}
	}

	// if imgMessage != nil {
	// 	if imgMessage.GetCaption() == "!sticker" {
	// 		go ServiceAppContext.Context.SendService.SendImageSticker(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageStickerParams{
	// 			ImageMessage: imgMessage,
	// 		})
	// 	}
	// }

	if videoMessage != nil {
		if videoMessage.GetCaption() == "!sticker" {
			go ServiceAppContext.Context.SendService.SendVideoSticker(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageStickerParams{
				VideoMessage: videoMessage,
			})
		}
	}

	if command == "!sticker" {
		stickerImageMessage := evt.RawMessage.ExtendedTextMessage.ContextInfo.QuotedMessage.GetImageMessage()
		stickerVideoMessage := evt.RawMessage.ExtendedTextMessage.ContextInfo.QuotedMessage.GetVideoMessage()

		if stickerImageMessage != nil {
			go ServiceAppContext.Context.SendService.SendImageSticker(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageStickerParams{
				ImageMessage: stickerImageMessage,
			})

		} else if stickerVideoMessage != nil {
			go ServiceAppContext.Context.SendService.SendVideoSticker(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageStickerParams{
				VideoMessage: stickerVideoMessage,
			})
		}
		return
	}

	if command == "!audio" {
		parts := strings.Split(messageText, " ")
		originalCommand := parts[0]
		audioName := strings.Replace(messageText, originalCommand, "", -1)

		go ServiceAppContext.Context.SendService.SendAudioFunny(context.Background(), fromChat, sender, audioName, stanzaID, messageText)
	}

	if command == "!transcrever" {
		audioMessage := evt.RawMessage.ExtendedTextMessage.ContextInfo.QuotedMessage.GetAudioMessage()
		// fmt.Println("audio message transcrever", audioMessage)
		if audioMessage != nil {
			go ServiceAppContext.Context.MessageService.ConvertMessageAudioToText(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageParams{
				AudioMessage: audioMessage,
			})
		}
	}

	if evt.Info.IsGroup {
		if command == "@todes" {
			go ServiceAppContext.Context.SendService.SendMessage(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageParams{
				MentionAllUsers:     true,
				ExtendedTextMessage: evt.RawMessage.GetExtendedTextMessage(),
			})

		}
	} else {
		// audioMessage := evt.Message.GetAudioMessage()
		// // fmt.Println("audio message direto", audioMessage)
		// if audioMessage != nil {
		// 	go ServiceAppContext.Context.MessageService.ConvertMessageAudioToText(context.Background(), fromChat, sender, argument, stanzaID, messageText, domainBotTypes.SendMessageParams{
		// 		AudioMessage: audioMessage,
		// 	})
		// }
	}

	// if command == "@supremacy" {
	// 	go handleTodes(context.Background(), fromChat, sender, argument, stanzaID, messageText, DomainBot.SendMessageParams{
	// 		MentionAdmin: true,
	// 	})
	// }

	// if command == "!balinha" {
	// 	go ServiceAppContext.Context.SendService.SendMessage(context.Background(), fromChat, sender, argument, stanzaID, messageText, DomainBot.SendMessageParams{
	// 		Message: "Aí é com o famoso 😉",
	// 	})
	// }

	// if !evt.Info.IsGroup {
	// 	if command == "!documentacao" {
	// 		messageDoc := "Verificação documental completa."
	// 		if argument != "h1234" {
	// 			messageDoc = "Nenhum documento encontrado."
	// 		}

	// 		ServiceAppContext.Context.SendService.SendMessage(context.Background(), fromChat, sender, argument, stanzaID, messageText, DomainBot.SendMessageParams{
	// 			Message:         messageDoc,
	// 			IsQuotedMessage: true,
	// 		})
	// 	}

	// 	if command == "!localizacao" {
	// 		messageLoc := "TE"
	// 		if argument != "h1234" {
	// 			messageLoc = "Nenhuma localização encontrada para esse documento."
	// 		}

	// 		ServiceAppContext.Context.SendService.SendMessage(context.Background(), fromChat, sender, argument, stanzaID, messageText, DomainBot.SendMessageParams{
	// 			Message:         messageLoc,
	// 			IsQuotedMessage: true,
	// 		})
	// 	}
	// }
}

// func handleTodes(ctx context.Context, fromChat string, sender string, name string, stanzaID string, messageText string, sendMessageParams DomainBot.SendMessageParams) {
// 	ServiceAppContext.Context.SendService.SendMessage(ctx, fromChat, sender, name, stanzaID, messageText, sendMessageParams)
// }
