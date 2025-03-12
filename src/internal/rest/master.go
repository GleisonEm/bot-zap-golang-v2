package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	domainMaster "github.com/gleisonem/bot-zap-golang-v2/domains/master"
	"github.com/gleisonem/bot-zap-golang-v2/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Master struct {
	Service domainMaster.IMasterService
}

func InitRestMaster(app *fiber.App) {
	// rest := Master{Service: service}
	app.Get("/master/download", DownloadFile)
	app.Post("/master/pipedrive/register", RegisterLeadPipeDrive)
}

func DownloadFile(c *fiber.Ctx) error {
	var request domainMaster.MasterRequestDownloadFile
	err := c.QueryParser(&request)
	utils.PanicIfNeeded(err)
	filePath := request.FilePath
	fileName := "data.wav"

	return c.Download(filePath, fileName)
}

func RegisterLeadPipeDrive(c *fiber.Ctx) error {

	var request domainMaster.MasterRequestRegisterLeadPipeDrive
	err := c.BodyParser(&request)
	utils.PanicIfNeeded(err)
	go func(request domainMaster.MasterRequestRegisterLeadPipeDrive) {
		logrus.Info("register lead body:", request)

		propertyValue, err := strconv.Atoi(request.PropertyValue)
		if err != nil {
			fmt.Println("Erro convert property valuye:", err)
			utils.PanicIfNeeded(err)
		}

		propertyInputValue, err := strconv.Atoi(request.PropertyInputValue)
		if err != nil {
			fmt.Println("Erro convert property input valuye:", err)
			utils.PanicIfNeeded(err)
		}

		monthlyIncome, err := strconv.Atoi(*request.MonthlyIncome)
		if err != nil {
			fmt.Println("Erro convert property input valuye:", err)
			monthlyIncome = 0
		}

		pathSearchPerson := "persons/search?term=" + request.LeadPhone + "&fields=phone&exact_match=true&limit=1"
		resp, err := sendRequest(http.MethodPost, pathSearchPerson, nil)

		if err != nil {
			logrus.Error("error doing request:", err, resp)
			utils.PanicIfNeeded(err)
		}

		var pipedriveResp domainMaster.PipedriveSearchPersonResponse
		err = json.NewDecoder(resp.Body).Decode(&pipedriveResp)
		if err != nil {
			logrus.Error("error decoding response:", err)
			utils.PanicIfNeeded(err)
		}

		var personId int64

		if pipedriveResp.Success && len(pipedriveResp.Data.Items) != 0 {
			logrus.Error("person not found in pipedrive or api error")
			personId = pipedriveResp.Data.Items[0].Item.ID
		} else {

			bodyCreatePerson := map[string]interface{}{
				"name":     request.LeadName,
				"owner_id": 23575609,
				"org_id":   7,
				"phone": []map[string]interface{}{
					{
						"label":   "mobile",
						"value":   request.LeadPhone, // Certifique-se de que essa variável é uma string
						"primary": true,
					},
				},
				"label": 27,
				"label_ids": []int64{
					27, // Certifique-se de que esta variável é uma string
				},
				"visible_to": "3",
			}
			resp, err := sendRequest(http.MethodPost, "persons", bodyCreatePerson)

			if err != nil {
				logrus.Error("error doing request create person:", err, resp)
				utils.PanicIfNeeded(err)
			}

			var createPersonResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&createPersonResp)
			if err != nil {
				logrus.Error("error decoding response create person:", err)
				utils.PanicIfNeeded(err)
			}

			if dataCreatePerson, ok := createPersonResp["data"].(map[string]interface{}); ok {
				if id, exists := dataCreatePerson["id"].(float64); exists {
					fmt.Println("ID:", int(id)) // Convertendo para int
					personId = int64(id)
				} else {
					fmt.Println("ID não encontrado.")
					errNotFoundId := errors.New("ID não encontrado")
					utils.PanicIfNeeded(errNotFoundId)
				}
			} else {
				fmt.Println("Campo 'data' não encontrado.")
				errNotFoundId := errors.New("ID não encontrado")
				utils.PanicIfNeeded(errNotFoundId)
			}
		}

		var leadPotentialValue string
		var optionsLeadMap = map[string]string{
			"quente": "1289601a-0e0a-440a-8e27-bab8062ec695",
			"morno":  "77ae55d8-f174-48d2-a6a6-e874101d9e4e",
			"frio":   "93ff3062-5fec-41a0-8508-83c943139602",
		}

		nameLower := strings.ToLower(request.LeadPotentialValue) // Converte para minúsculas
		leadPotentialValue, exists := optionsLeadMap[nameLower]
		if !exists {
			leadPotentialValue = "93ff3062-5fec-41a0-8508-83c943139602"
		}

		now := time.Now()
		futureDate := now.AddDate(0, 0, 7)
		formattedDate := futureDate.Format("2006-01-02")

		bodyCreateLead := map[string]interface{}{
			"title":    request.LeadName,
			"owner_id": 23575609,
			"label_ids": []string{
				leadPotentialValue, // Certifique-se de que esta variável é uma string
			},
			"person_id":       personId,
			"organization_id": 7,
			"value": map[string]interface{}{
				"amount":   propertyValue,
				"currency": "BRL",
			},
			"expected_close_date": formattedDate,
			"visible_to":          "3",
			"was_seen":            true,
			"origin_id":           nil,
			"channel":             1,
			"channel_id":          "whatsapp",
			"3878984867c57db15672e501024c9668e1a63998":          request.LeadPotentialValue,
			"b806ba8632c15df94c35c8655160af35861835b2":          request.Locality,
			"4fbc04c3a255cf36f7f8cbb454f7bccaa06e151e":          propertyInputValue,
			"4fbc04c3a255cf36f7f8cbb454f7bccaa06e151e_currency": "BRL",
		}

		if monthlyIncome != 0 {

			bodyCreateLead["dcc65cd9c497e5cb3eb426a3a4d069cb82dabac7"] = monthlyIncome
			bodyCreateLead["dcc65cd9c497e5cb3eb426a3a4d069cb82dabac7_currency"] = "BRL"
		}

		// $SELECTION_PLACEHOLDER$ modified:
		resp, err = sendRequest(http.MethodPost, "leads", bodyCreateLead)
		if err != nil {
			logrus.Error("error doing request:", err, resp)
			utils.PanicIfNeeded(err)
		}
	}(request)

	return c.JSON(utils.ResponseData{
		Status:  204,
		Code:    "SUCCESS",
		Message: "Success register lead",
	})
}

func sendRequest(method, path string, payload interface{}) (*http.Response, error) {
	var req *http.Request
	var err error
	client := &http.Client{Timeout: 10 * time.Second}

	pipeDriveApi := os.Getenv("PIPEDRIVE_API")
	if pipeDriveApi == "" {
		logrus.Info("not api url pipedrive found:", pipeDriveApi)
	}

	pipeDriveApiKey := os.Getenv("PIPEDRIVE_API_KEY")
	if pipeDriveApiKey == "" {
		logrus.Info("not api key pipedrive found:", pipeDriveApiKey)
	}

	url := fmt.Sprintf("%s/v1/%s?api_token=%s", pipeDriveApi, path, pipeDriveApiKey)

	if payload != nil {
		data, err := json.Marshal(payload)
		logrus.Info(path+" payload: ", string(data))

		if err != nil {
			logrus.Error("error marshal payload:", err)
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
		if err != nil {
			logrus.Error("error creating new request:", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			logrus.Error("error creating new request:", err)
			return nil, err
		}
	}
	return client.Do(req)
}
