package event

type CreateEventInput struct {
	Name     string `json:"name"`
	StartsAt string `json:"starts_at"`
}
