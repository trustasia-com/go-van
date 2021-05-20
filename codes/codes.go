// Package codes provides ...
package codes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// list languages
const (
	LangZhCN = "zh-cn"
	LangEnUS = "en-us"
	// and more...

	// system embeded codes
	// same: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
	// adjust for custom status
	OK                 Code = 0  // 0
	Canceled           Code = 1  // 1
	Unknown            Code = 2  // 2
	InvalidArgument    Code = 3  // 3
	DeadlineExceeded   Code = 4  // 4
	NotFound           Code = 5  // 5
	AlreadyExist       Code = 6  // 6
	PermissionDenied   Code = 7  // 7
	ResourceExhausted  Code = 8  // 8
	FailedPrecondition Code = 9  // 9
	Aborted            Code = 10 // 10
	OutOfRange         Code = 11 // 11
	Unimplemented      Code = 12 // 12
	Internal           Code = 13 // 13
	Unavailable        Code = 14 // 14
	DataLoss           Code = 15 // 15
	Unauthenticated    Code = 16 // 16

	_maxCode = 17
)

var (
	// globalI18n global instance
	globalI18n = &i18nInstance{}

	// code2Desc embeded code description
	code2Desc = map[string]map[Code]string{
		LangZhCN: {
			OK:                 "OK",
			Canceled:           "操作被取消",
			Unknown:            "未知错误",
			InvalidArgument:    "请求参数无效",
			DeadlineExceeded:   "操作超时",
			NotFound:           "资源未找到",
			AlreadyExist:       "资源已存在",
			PermissionDenied:   "权限不足",
			ResourceExhausted:  "资源耗尽",
			FailedPrecondition: "前置条件不足",
			Aborted:            "操作被终止",
			OutOfRange:         "索引越界",
			Unimplemented:      "方法未实现",
			Internal:           "内部错误",
			Unavailable:        "服务繁忙",
			DataLoss:           "数据丢失",
			Unauthenticated:    "未经认证",
		},
		LangEnUS: {
			OK:                 "OK",
			Canceled:           "Canceled",
			Unknown:            "Unknown",
			InvalidArgument:    "InvalidArgument",
			DeadlineExceeded:   "DeadlineExceeded",
			NotFound:           "NotFound",
			AlreadyExist:       "AlreadyExist",
			PermissionDenied:   "PermissionDenied",
			ResourceExhausted:  "ResourceExhausted",
			FailedPrecondition: "FailedPrecondition",
			Aborted:            "Aborted",
			OutOfRange:         "OutOfRange",
			Unimplemented:      "Unimplemented",
			Internal:           "Internal",
			Unavailable:        "Unavailable",
			DataLoss:           "DataLoss",
			Unauthenticated:    "Unauthenticated",
		},
	}
)

// Code error code
type Code uint32

// StatusCode http status code
func (c Code) StatusCode() int {
	switch c {
	case OK:
		return http.StatusOK
	case PermissionDenied:
		return http.StatusForbidden
	case Internal:
		return http.StatusInternalServerError
	case Unauthenticated:
		return http.StatusUnauthorized
	case FailedPrecondition:
		return http.StatusPreconditionFailed
	case AlreadyExist:
		return http.StatusConflict
	case DeadlineExceeded:
		return http.StatusRequestTimeout
	case Unimplemented:
		return http.StatusNotImplemented
	case Unavailable:
		return http.StatusServiceUnavailable
	}
	// other codes
	return http.StatusBadRequest
}

// GRPCCodeFromStatus converts a HTTP error code into the corresponding gRPC response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func GRPCCodeFromStatus(code Code) Code {
	switch code {
	case http.StatusOK:
		return OK
	case http.StatusBadRequest:
		return InvalidArgument
	case http.StatusUnauthorized:
		return Unauthenticated
	case http.StatusForbidden:
		return PermissionDenied
	case http.StatusNotFound:
		return NotFound
	case http.StatusConflict:
		return Aborted
	case http.StatusTooManyRequests:
		return ResourceExhausted
	case http.StatusInternalServerError:
		return Internal
	case http.StatusNotImplemented:
		return Unimplemented
	case http.StatusServiceUnavailable:
		return Unavailable
	case http.StatusGatewayTimeout:
		return DeadlineExceeded
	}
	return Unknown
}

// Tr translate code to description
func (c Code) Tr(lang string, args ...interface{}) string {
	if globalI18n.translator == nil {
		return "codes: warning: please specific translator"
	}
	// judge language
	found := false
	for _, l := range globalI18n.supportedLang {
		if l != lang {
			found = true
			break
		}
	}
	if !found {
		lang = globalI18n.supportedLang[0]
	}
	// tansplate code
	if c < _maxCode {
		codes, ok := code2Desc[lang]
		if !ok {
			return "codes: warning: unsupported lang " + lang
		}
		desc := codes[c]
		for _, arg := range args {
			desc += fmt.Sprintf("|%v", arg)
		}
		return desc
	}
	return globalI18n.translator.Tr(lang, c, args...)
}

// String convert to string
func (c Code) String() string {
	return "(" + strconv.FormatInt(int64(c), 10) + ")"
}

// GetSupportedLang get supported lang
func GetSupportedLang() []string {
	return globalI18n.supportedLang
}

// SetTranslator change default translator
func SetTranslator(trans translator) error {
	langs := trans.SupportedLang()
	if len(langs) == 0 {
		return errors.New("codes: warning: not found supported lang")
	}
	globalI18n.translator = trans
	// use first language as default
	globalI18n.supportedLang = langs
	return nil
}
