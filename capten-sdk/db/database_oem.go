package db

type OemName int

const (
	UNKNOWN OemName = iota + 1

	// Postgres
	POSTGRES
	// Future enhancements
	// MS MySQL
	// MD MariaDB
)

var oemNames = map[string]OemName{
	UNKNOWN.String():  UNKNOWN,
	POSTGRES.String(): POSTGRES,
}

func (o OemName) String() string {
	return [...]string{
		"UNKNOWN",
		"POSTGRES",
	}[o-1]
}

func (o OemName) EnumIndex() int {
	return int(o)
}

func GetEnum(key string) (OemName, bool) {
	if val, ok := oemNames[key]; ok {
		return val, ok
	}
	return UNKNOWN, false
}

func GetSupportedDatabaseOEMs() []string {
	supportedDBs := make([]string, len(oemNames))
	for oemName := range oemNames {
		if oemName != UNKNOWN.String() {
			supportedDBs = append(supportedDBs, oemName)
		}
	}
	return supportedDBs
}
