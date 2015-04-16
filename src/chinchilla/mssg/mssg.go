package mssg

type Connect struct {
	Type uint8  "json:`type`" // Operation
	Id   uint32 // Node ID
	QVal uint32 // Value of q (0 is default)
}

type WorkReq struct {
	Type uint8
	Arg1 string
	Arg2 string
	Arg3 uint32
	WId  uint32
}

type WorkResp struct {
	Type  uint8
	Id    uint32
	Data  []byte // Can be json kind of thing or string
	WId   uint32
	RTime uint32
}
