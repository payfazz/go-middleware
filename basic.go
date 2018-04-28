package middleware

// Basic return middleware pack that contain Logger, Recovery, and KV
func Basic() []Func {
	return BuildList(
		NewLogger(nil),
		NewRecovery(10, nil),
		WithKV(),
	)
}
