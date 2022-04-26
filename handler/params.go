package handler

import "github.com/julienschmidt/httprouter"

func ParseParam(params httprouter.Params, key string) string {
	return params.ByName(key)[1:]
}
