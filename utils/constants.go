package utils

const (
	IniName            = ".dhcore.ini"
	CurrentEnvironment = "current_environment"
	configFile         = "config.json"
	ApiLevelKey        = "dhcore_api_level"

	// API level the current version of the CLI was developed for
	MinApiLevel = 10

	// Individual commands; 0 means no restriction
	LoginMin  = 10
	LoginMax  = 0
	CreateMin = 10
	CreateMax = 0
	ListMin   = 10
	ListMax   = 0
	GetMin    = 10
	GetMax    = 0
	UpdateMin = 10
	UpdateMax = 0
	DeleteMin = 10
	DeleteMax = 0
)
