package quickjs

type GlobalConfig struct {
	ManualFree bool
}

var globalConfig GlobalConfig

func SetManualFree() {
	globalConfig.ManualFree = true
}
