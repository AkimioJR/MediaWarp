package alist

type FsGetRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Page     uint32 `json:"page"`
	PerPage  uint32 `json:"per_page"`
	Refresh  bool   `json:"refresh"`
}
