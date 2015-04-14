package mssg

type Connect struct {
	Type uint8  // Operation
	Id   uint32 // Node ID
	QVal uint32
}

type WorkReq struct {
	Type    uint8
	Arg1    string
	Arg2    string
	Arg3    uint32
	SrcIp   string
	SrcPort string
}

type WorkResp struct {
	Type    uint8
	Id      uint32
	Data    []byte // Can be json kind of thing or string
	SrcIp   string
	SrcPort string
	RTime   uint32
}
