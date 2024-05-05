package enum

type VideoStatus int

const (
	VideoStatusUnknown VideoStatus = iota
	VideoStatusProcessing
	VideoStatusProcessed
	VideoStatusFailed
	VideoStatusDisabled
)

func (s VideoStatus) String() string {
	switch s {
	case VideoStatusProcessing:
		return "processing"
	case VideoStatusProcessed:
		return "processed"
	case VideoStatusFailed:
		return "failed"
	case VideoStatusDisabled:
		return "disabled"
	default:
		panic("unhandled default case")
	}
}
