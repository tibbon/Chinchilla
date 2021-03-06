package mssg

import (
	"time"
)

type Connect struct {
	Type uint8   `json:"type"` // Operation
	Id   uint32  // Node ID
	QVal float64 // Value of q (0 is default)
}

type WorkReq struct {
	Type  uint8
	Arg1  string
	Arg2  string
	Arg3  uint32
	WId   uint32
	STime time.Time
}

type WorkResp struct {
	Type  uint8        `json:"type"`
	Id    uint32       `json:"id"`
	Data  WorkRespData `json:"data"` // Can be json kind of thing or string
	WId   uint32       `json:"work_id"`
	PTime float64      `json:"return_time"`
	QVal  float64      `json:"q_val"`
}

type WorkRespData struct {
	Desc string `json:"description"`
}
