package sep

type CheckWire struct {
	Ok bool `json:"ok"`
	N  int  `json:"n"`
}

var liveCheck = CheckWire{
	Ok: false,
	N:  1,
}

func HoldCheckLive(cur CheckWire) CheckWire {
	out := liveCheck
	liveCheck = cur
	return out
}
