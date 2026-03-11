package common

import (
	"net/http"
	"strconv"
)

func GetPathInt(r *http.Request, name string, defaultValue int) int {
	v := r.PathValue(name)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}
