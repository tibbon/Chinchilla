package mssg

type Connect struct {
	Type uint8  // Operation
	Id   uint32 // Node ID
}

type WorkReq struct {
	Type  uint8
	arg1  string
	arg2  string
	arg3  uint32
	SrcIp string
}

type WorkResp struct {
	Type  uint8
	Id    uint32
	data  []byte
	SrcIp string
}
