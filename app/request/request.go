package request

type ReqUID struct {
	UID int64 `uri:"uid" form:"uid" json:"uid"`
}

type ReqPage struct {
	Page     int `uri:"page" form:"page" json:"page"`
	PageSize int `uri:"page_size" form:"page_size" json:"page_size"`
}

const pageSizeMin = 1
const pageSizeMax = 100

func (p *ReqPage) ValidatePageSize() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSize <= pageSizeMin {
		p.PageSize = pageSizeMin
	}

	if p.PageSize > pageSizeMax {
		p.PageSize = pageSizeMax
	}
}
