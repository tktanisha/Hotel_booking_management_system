package payloads

type CreateHotelPayload struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
