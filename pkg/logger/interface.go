package logger

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	KVLog(k string, v interface{})
	Fatal(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
}
