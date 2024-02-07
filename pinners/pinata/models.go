package pinata

type addEvent struct {
	IpfsHash  string
	PinSize   int64  `json:",omitempty"`
	Timestamp string `json:",omitempty"`
}
