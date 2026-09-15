package event

type SuccessResponse struct {
	Success string `json:"success"`
	Event   Event  `json:"event"`
}

type EventsResponse struct {
	Success string  `json:"success"`
	Events  []Event `json:"events"`
}
