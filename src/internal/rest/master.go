package rest

import (
	domainMaster "github.com/gleisonem/bot-zap-golang-v2/domains/master"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Master struct {
	Service domainMaster.IMasterService
}

func InitRestMaster(app *fiber.App) {
	// rest := Master{Service: service}
	app.Get("/master/download", DownloadFile)
}

func DownloadFile(c *fiber.Ctx) error {
	var request domainMaster.MasterRequestDownloadFile
	err := c.QueryParser(&request)
	utils.PanicIfNeeded(err)
	filePath := request.FilePath
	fileName := "data.wav"

	return c.Download(filePath, fileName)
}
