package core

type Priority int

const (
	PriorityHigh Priority = 1
	PriorityMid  Priority = 2
	PriorityLow  Priority = 3
)

type Format string

const (
	FormatJPEG Format = "jpeg"
	FormatPNG  Format = "png"
	FormatWEBP Format = "webp"
)
