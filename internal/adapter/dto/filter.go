package dto

type Filter struct {
	Key          string   `json:"key"`
	Type         string   `json:"type"`
	Options      []string `json:"options"`
	OptionSource string   `json:"option_source"`
}

type FilterResponse struct {
	Filters []Filter
}

type FilterRequest struct {
	Filters map[string]string `json:"filters"`
}
