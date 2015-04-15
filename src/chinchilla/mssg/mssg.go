package mssg

import (
	"net/http"
)

type Connect struct {
	Type uint8  // Operation
	Id   uint32 // Node ID
	QVal uint32
}

type WorkReq struct {
	Type uint8
	Arg1 string
	Arg2 string
	Arg3 uint32
	W    http.ResponseWriter
}

type WorkResp struct {
	Type  uint8
	Id    uint32
	Data  []byte // Can be json kind of thing or string
	W     http.ResponseWriter
	RTime uint32
}
