package request

type ReqRegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqUsername struct {
	Username string `uri:"username" form:"username" json:"username"`
}

type ReqEmail struct {
	Email string `uri:"email" form:"email" json:"email"`
}
