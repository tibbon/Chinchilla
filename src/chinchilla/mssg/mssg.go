package mssg

type Connect struct {
	Type uint8  // Operation
	Id   uint32 // Node ID
	QVal uint32
}

type WorkReq struct {
	Type  uint8
	Arg1  string
	Arg2  string
	Arg3  uint32
	SrcIp string
	Port  string
}

type WorkResp struct {
	Type  uint8
	Id    uint32
	Data  []byte
	SrcIp string
	Port  string
	RTime uint32
}
