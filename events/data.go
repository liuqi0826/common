package events

type SystemEventData struct {
}
type MouseEventData struct {
	X          int
	Y          int
	MovementX  int
	MovementY  int
	StageX     int
	StageY     int
	AltKey     bool
	CtrlKey    bool
	ShiftKey   bool
	ButtonDown bool
	CommandKey bool
	Delta      int
}
type KeyboardEventData struct {
	CharCode    int
	KeyCode     int
	KeyLocation int
	ScanCode    int
	Key         string
	AltKey      bool
	CtrlKey     bool
	ShiftKey    bool
	CommandKey  bool
}
type WindowEventData struct {
	WindowTitle     string
	WindowWidth     int
	WindowHeight    int
	WindowPositionX int
	WindowPositionY int
	Resizable       bool
	MonitorIndex    int
}
