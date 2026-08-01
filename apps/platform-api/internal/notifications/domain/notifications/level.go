package notifications

type Level string

const (
	LevelInfo    Level = "info"
	LevelSuccess Level = "success"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

var (
	levels = map[Level]struct{}{
		LevelInfo:    {},
		LevelSuccess: {},
		LevelWarning: {},
		LevelError:   {},
	}
)

func (l Level) String() string {
	return string(l)
}

func (l Level) IsValid() bool {
	_, ok := levels[l]
	return ok
}
