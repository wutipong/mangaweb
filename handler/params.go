package handler

import "github.com/julienschmidt/httprouter"

func ParseParam(params httprouter.Params, key string) string {
	// the catchAll param always start with / .
	return params.ByName(key)[1:]
}
