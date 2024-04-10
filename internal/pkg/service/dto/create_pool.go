package dto

type CreatePoolCmd struct {
	Address        string
	StartBlock     uint64
	Token0         string
	Token0Decimals uint
	Token1         string
	Token1Decimals uint
}
