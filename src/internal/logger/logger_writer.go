package logger

type dualLoggerWriter struct {
	_write func([]byte) (int, error)
}

func (dlw *dualLoggerWriter) Write(p []byte) (n int, err error) {
	return dlw._write(p)
}
