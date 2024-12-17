package dto

type ServicesResponse struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name      string `json:"name"`
	Container string `json:"container"`
	Port      string `json:"port"`
	State     string `json:"state"`
}
