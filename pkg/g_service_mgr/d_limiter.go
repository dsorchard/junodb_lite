package service

type (
	ILimiter interface {
		LimitReached() bool
		Throttle()
	}
)
