package middleware

// Basic return middleware pack that contain Logger, Recovery, and WithKV
func Basic() []Func {
	return BuildList(
		NewLogger(nil),
		NewRecovery(10, nil),
		WithKV(),
	)
}
