package schema

// PagingList 分页列表
type PagingList struct {
	PageIndex      int         `json:"pageIndex"`
	PageSize       int         `json:"pageSize"`
	TotalCount     int         `json:"totalCount"`
	PageTotalCount int         `json:"pageTotalCount"`
	Items          interface{} `json:"items"`
}
