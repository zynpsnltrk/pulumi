package apitype

// PreviewEvent describes the contents of a single preview event.
type PreviewEvent AppendUpdateLogEntryRequest

// PutPreviewResultsRequest describes the format of a request to POST preview/results
type PutPreviewResultsRequest struct {
	Environment map[string]string `json:"environment"`
	Succeeded   bool              `json:"succeeded"`
	Results     []PreviewEvent    `json:"results"`
}

// PutPreviewResultsResponse describes the format of a response to POST preview/results
type PutPreviewResultsResponse struct {
	ID string `json:"id"`
}

// GetPreviewResultsResponse describes the format of a response to GET preview/results/{previewID}
type GetPreviewResultsResponse struct {
	ID        string         `json:"id"`
	Succeeded bool           `json:"succeeded"`
	Results   []PreviewEvent `json:"results"`
}
