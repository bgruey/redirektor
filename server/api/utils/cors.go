package utils

import (
	"net/http"
)

// leftover from another project, not used
func EnableAllCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
}
