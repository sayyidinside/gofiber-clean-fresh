package model

type (
	QueryGet struct {
		Page     string `query:"page"`
		Limit    string `query:"limit"`
		SearchBy string `query:"search_by"`
		Search   string `query:"search"`
		FilterBy string `query:"filter_by"`
		Filter   string `query:"filter"`
		OrderBy  string `query:"order_by"`
		Order    string `query:"order"`
	}
)
