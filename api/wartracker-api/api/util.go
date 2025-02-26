package api

import (
	"net/http"
	"strconv"
)

func GetQueryBool(r *http.Request, k string) bool {
	var b bool = false

	q := r.URL.Query()[k]
	if len(q) > 0 {
		b, _ = strconv.ParseBool(q[0])
	}

	return b
}

func GetQueryInt(r *http.Request, k string) int {
	var i int = 0

	q := r.URL.Query()[k]
	if len(q) > 0 {
		i, _ = strconv.Atoi(q[0])
	}

	return i
}
