package master

import "context"

// IMainService defines the interface for main service operations.
type IMasterService interface {
	DownloadFile(ctx context.Context, request MasterRequestDownloadFile) (response MasterResponse, err error)
}

// MainRequest represents a sample request for main info.
type MasterRequestDownloadFile struct {
	FilePath string `json:"file_path" query:"file_path"`
}

// MainResponse represents a sample response for main info.
type MasterResponse struct {
	Data []string `json:"data"`
}
