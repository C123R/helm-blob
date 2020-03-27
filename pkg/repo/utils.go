package repo

import (
	"os"
	"path/filepath"
)

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

func getChartFileName(chartpath string) string {
	_, file := filepath.Split(chartpath)

	return file
}
