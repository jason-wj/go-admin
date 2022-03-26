package config

type Application struct {
	ReadTimeout   int
	WriterTimeout int
	Host          string
	Port          int64
	Name          string
	Mode          string
	EnableDP      bool
	FileRootPath  string
	AmpKey        string
	IsSingleLogin bool
}

var ApplicationConfig = new(Application)
