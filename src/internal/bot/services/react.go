package services

import (
	"context"

	ServiceAppContext "github.com/gleisonem/bot-zap-golang-v2/contexts"
)

func SendReact(ctx context.Context, sender string, stanzaID string, emojiReaction string, fromChat string) {
	ServiceAppContext.Context.AppService.SendInWhatsapp(
		ctx,
		"mandado react",
		ServiceAppContext.Context.AppService.MakeBuildReaction(
			ctx,
			sender,
			stanzaID,
			emojiReaction,
		),
		fromChat,
	)
}
