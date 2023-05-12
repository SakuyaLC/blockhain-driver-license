package model

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Info      string
	Hash      string
	PrevHash  string
	Validator string
}
