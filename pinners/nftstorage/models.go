package nftstorage

type value struct {
	Cid     string
	Size    int64  `json:",omitempty"`
	Created string `json:",omitempty"`
	Type    string `json:",omitempty"`
}

type er struct {
	Name, Message string
}

type addEvent struct {
	Ok    bool
	Value value
	Error er
}
