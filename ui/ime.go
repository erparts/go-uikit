package ui

// IMEBridge is implemented on the Java side and registered from your mobile package.
// In this minimal integration, it only opens/closes the keyboard.
type IMEBridge interface {
	Show()
	Hide()
}

// WantsIME marks widgets that should trigger IME show/hide when focused.
type WantsIME interface {
	WantsIME() bool
}
