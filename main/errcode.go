package main

var errorMap = map[int]string{
	0:      "success",
	-61001: "err data",
	-61002: "err server busy",
	-61003: "err not exist",
	-61005: "err req param",
}

var (
	Success       = 0
	ErrData       = -61001
	ErrInnerFault = -61002
	ErrNotExist   = -61003
	ErrReqParam   = -61005
)
