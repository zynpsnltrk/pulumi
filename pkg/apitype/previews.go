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

// PreviewState describes the state of a preview.
type PreviewState string

const (
	// PreviewStateAwaitingReview indicates that the preview is awaiting review.
	PreviewStateAwaitingReview = "waiting"

	// PreviewStateApproved indicates that the preview has been approved.
	PreviewStateApproved = "approved"

	// PreviewStateRejected indicates that the preview has been rejected.
	PreviewStateRejected = "rejected"
)

// GetPreviewResultsResponse describes the format of a response to GET preview/results/{previewID}
type GetPreviewResultsResponse struct {
	ID        string         `json:"id"`
	Succeeded bool           `json:"succeeded"`
	State     PreviewState   `json:"state"`
	Results   []PreviewEvent `json:"results"`
}
