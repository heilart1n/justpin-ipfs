package infura

type addEvent struct {
	Name  string
	Hash  string `json:",omitempty"`
	Bytes int64  `json:",omitempty"`
	Size  string `json:",omitempty"`
}
