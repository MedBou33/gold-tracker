package alert

func Crossed(previous, current, target float64) bool {
	crossedUp := previous < target && current >= target
	crossedDown := previous >= target && current < target
	return crossedUp || crossedDown
}
