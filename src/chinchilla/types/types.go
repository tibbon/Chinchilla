package types

import (
	"chinchilla/mssg"
	"encoding/gob"
	"net/http"
	"sync"
)

type Queue struct {
	QVal     float64
	Enc      *gob.Encoder
	Sent     bool
	Reqs     []mssg.WorkReq
	AvgTimes map[uint8]float64
}
type Job struct {
	W   http.ResponseWriter
	Sem chan struct{}
}

type MapJ struct {
	M map[uint32]Job
	L *sync.RWMutex
}

type MapQ struct {
	M map[uint32]Queue
	L *sync.RWMutex
}

type Stack struct {
	S []uint32
	L *sync.RWMutex
}
