package global

type JvmEnv struct {
	debug bool
}

var jvmEnv JvmEnv

func GlobalEnv() JvmEnv {
	return jvmEnv
}

func SetDebugMode(on bool) {
	jvmEnv.debug = on
}

func DebugMode() bool {
	return jvmEnv.debug
}
