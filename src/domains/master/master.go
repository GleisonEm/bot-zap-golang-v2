package master

import "context"

// IMainService defines the interface for main service operations.
type IMasterService interface {
	DownloadFile(ctx context.Context, request MasterRequestDownloadFile) (response MasterResponse, err error)
	RegisterLeadPipeDrive(ctx context.Context, request MasterRequestRegisterLeadPipeDrive) (response MasterResponse, err error)
}

// MasterRequest represents a sample request for main info.
type MasterRequestDownloadFile struct {
	FilePath string `json:"file_path" query:"file_path"`
}

// MasterResponse represents a sample response for main info.
type MasterResponse struct {
	Data []string `json:"data"`
}

type MasterRequestRegisterLeadPipeDrive struct {
	PropertyValue      string  `json:"property_value" form:"property_value"`
	PropertyInputValue string  `json:"property_input_value" form:"property_input_value"`
	Locality           string  `json:"locality" form:"locality"`
	LeadName           string  `json:"lead_name" form:"lead_name"`
	LeadPhone          string  `json:"lead_phone" form:"lead_phone"`
	LeadPotentialValue string  `json:"lead_potential_value" form:"lead_potential_value"`
	MonthlyIncome      *string `json:"monthly_income" form:"monthly_income"`
}

type PipedriveSearchPersonResponse struct {
	Success bool                      `json:"success"`
	Data    PipedriveSearchPersonData `json:"data"`
}

type PipedriveSearchPersonData struct {
	Items []PipedriveSearchPersonItem `json:"items"`
}

type PipedriveSearchPersonItem struct {
	ResultScore float64                     `json:"result_score"`
	Item        PipedriveSearchPersonDetail `json:"item"`
}

type PipedriveSearchPersonDetail struct {
	ID           int64                             `json:"id"`
	Type         string                            `json:"type"`
	Name         string                            `json:"name"`
	Phones       []string                          `json:"phones"`
	Emails       []string                          `json:"emails"`
	PrimaryEmail string                            `json:"primary_email"`
	VisibleTo    int                               `json:"visible_to"`
	Owner        PipedriveSearchPersonOwner        `json:"owner"`
	Organization PipedriveSearchPersonOrganization `json:"organization"`
}

type PipedriveSearchPersonOwner struct {
	ID int64 `json:"id"`
}

type PipedriveSearchPersonOrganization struct {
	ID      int64       `json:"id"`
	Name    string      `json:"name"`
	Address interface{} `json:"address"`
}
