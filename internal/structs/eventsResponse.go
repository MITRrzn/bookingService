package structs

type EventsResponse struct {
	Success string  `json:"success"`
	Events  []Event `json:"events"`
}
