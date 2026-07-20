package pagination

const (
	DEFAULT_LIMIT int32 = 5000
	MAX_PAGE_SIZE int32 = 200
)

type Pagination struct {
	PageSize   int32 `json:"page_limit"`
	PageNumber int32 `json:"page_number"`
}
