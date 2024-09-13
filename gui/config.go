package gui

type Config struct {
	Title  string
	SizeDp float32
}

func NewAppConfig() Config {
	return Config{
		Title:  "Application",
		SizeDp: 1600,
	}
}
