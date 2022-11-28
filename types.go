package main

// Block represents each 'item' in the blockchain
type Block1 struct {
	Index       int
	Timestamp   string
	LicenseInfo string
	Hash        string
	PrevHash    string
	Validator   string
}
