package contracts

import "strconv"

type PaginatedRequest struct {
	Page int `json:"page" query:"page"`
	Size int `json:"size" query:"size"`
}

type PaginatedRequestOrdered struct {
	PaginatedRequest
	Sort  string `json:"sort" query:"sort"`
	Order string `json:"order" query:"order"`
}

type PaginationSetter interface {
	SetPage(page int)
	SetSize(size int)
}

type PaginationGetter[T any] interface {
	GetPage() int
	GetSize() int
	GetTotal() int
	GetItems() []T
}

type PaginatedResponse[T any] struct {
	Page  int `json:"page" validate:"min=0"`
	Size  int `json:"size" validate:"min=0"`
	Total int `json:"total"`
	Items []T `json:"items"`
}

type PaginatedResponseOrdered[T any] struct {
	PaginatedResponse[T]
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func (req *PaginatedRequest) ToQueryParams() map[string]string {
	params := make(map[string]string, 2)
	if req.Page > 0 {
		params["page"] = strconv.Itoa(req.Page)
	}
	if req.Size > 0 {
		params["size"] = strconv.Itoa(req.Size)
	}
	return params
}

func (req *PaginatedRequestOrdered) ToQueryParams() map[string]string {
	params := req.PaginatedRequest.ToQueryParams()
	if req.Sort != "" {
		params["sort"] = req.Sort
	}
	if req.Order != "" {
		params["order"] = req.Order
	}
	return params
}

func (req *PaginatedRequest) SetPage(page int) {
	req.Page = page
}

func (req *PaginatedRequest) SetSize(size int) {
	req.Size = size
}

func (res *PaginatedResponse[T]) GetPage() int {
	return res.Page
}

func (res *PaginatedResponse[T]) GetSize() int {
	return res.Size
}

func (res *PaginatedResponse[T]) GetTotal() int {
	return res.Total
}

func (res *PaginatedResponse[T]) GetItems() []T {
	return res.Items
}
