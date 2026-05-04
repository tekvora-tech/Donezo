package types

type (
	Metadata struct {
		CurrentPage  int  `json:"current_page"`
		PageSize     int  `json:"page_size"`
		TotalPages   int  `json:"total_pages"`
		TotalRecords int  `json:"total_records"`
		HasNext      bool `json:"has_next"`
		HasPrev      bool `json:"has_prev"`
	}
)