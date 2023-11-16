package utils

type Pagination struct {
	Limit        int    `json:"limit"`
	CurrentPage  int    `json:"currentPage"`
	NextPage     string `json:"nextPage"`
	PreviousPage string `json:"previousPage"`
	SortField    string `json:"sortField"`
	SortOrder    string `json:"sortOrder"`
	TotalRecords int    `json:"totalRecords"`
}
