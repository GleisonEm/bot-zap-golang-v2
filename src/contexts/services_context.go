package contexts

import (
	domainApp "github.com/gleisonem/bot-zap-golang-v2/domains/bot/app"
	domainGroup "github.com/gleisonem/bot-zap-golang-v2/domains/bot/group"
	domainMessage "github.com/gleisonem/bot-zap-golang-v2/domains/bot/message"
	domainSend "github.com/gleisonem/bot-zap-golang-v2/domains/bot/send"
)

type AppContext struct {
	AppService     domainApp.IAppService
	SendService    domainSend.ISendService
	MessageService domainMessage.IMessageService
	GroupService   domainGroup.IGroupService
}

var Context *AppContext

func InitServiceAppContext(
	appService domainApp.IAppService,
	sendService domainSend.ISendService,
	messageService domainMessage.IMessageService,
	groupService domainGroup.IGroupService,
) {
	Context = &AppContext{
		AppService:     appService,
		SendService:    sendService,
		MessageService: messageService,
		GroupService:   groupService,
	}
}
