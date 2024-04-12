package dto

type PaginationReq struct {
	Page     uint `form:"page"`
	PageSize uint `form:"pageSize"`
}

type PaginationResp struct {
	Total int64 `json:"total"`
}
