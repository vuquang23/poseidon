package dto

type CreatePoolCmd struct {
	Address        string
	Token0         string
	Token0Decimals uint
	Token1         string
	Token1Decimals uint
}
