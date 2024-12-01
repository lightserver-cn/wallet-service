package request

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

const ContentTypeKey = "Content-Cate"
const ContentTypeJSON = "application/json"

type Response struct {
	ctx *gin.Context
}

type ResponseEntity struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    any    `json:"data"`
}

type WechatResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewResponse(ctx *gin.Context) Response {
	return Response{ctx: ctx}
}

func (r Response) Header() {
	r.ctx.Header(ContentTypeKey, ContentTypeJSON)
}

func (r Response) Success() {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    struct{}{},
	})
}

func (r Response) SuccessData(data any) {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    data,
	})
}

func (r Response) SuccessMsg(msg string) {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: 0,
		ErrMsg:  msg,
		Data:    struct{}{},
	})
}

func (r Response) SuccessDataMsg(data any, msg string) {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: 0,
		ErrMsg:  msg,
		Data:    data,
	})
}

func (r Response) ErrorCodeMsg(errCode int, msg string) {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: errCode,
		ErrMsg:  msg,
		Data:    nil,
	})
}

func (r Response) ErrorValidator(err error) {
	r.ctx.JSON(http.StatusOK, ResponseEntity{
		ErrCode: ErrCodeValidateErr,
		ErrMsg:  err.Error(),
		Data:    nil,
	})
}
