package security

var safeAPI = true

// AllowUnsafeAPI() enables some UNSAFE behaviors which should NOT be used in PRODUCTION env.
func AllowUnsafeAPI() {
	safeAPI = false
}

func UnsafeAPI() bool {
	return !safeAPI
}
