package handler

import (
	"github.com/wutipong/mangaweb/errors"
	"github.com/wutipong/mangaweb/log"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var errorTemplate *template.Template = nil

func initError() {
	var err error
	errorTemplate, err = template.New("error.gohtml").
		Funcs(HtmlTemplateFuncMap()).
		ParseFiles(
			"template/error.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Get().Sugar().Panic(err)
		os.Exit(-1)
	}
}

type ErrorResponse struct {
	Title   string
	Code    uint
	Message string
}

func WriteError(w http.ResponseWriter, err error) {
	var exErr errors.Error
	if e, ok := err.(errors.Error); ok {
		exErr = e
	} else {
		exErr = errors.ErrUnknown.Wrap(err)
	}
	log.Get().Error(
		"Error",
		zap.Error(err),
		zap.Uint("error_code", exErr.Code),
		zap.String("cause", exErr.Cause.Error()))

	data := ErrorResponse{
		Title:   "Error",
		Code:    exErr.Code,
		Message: exErr.Error(),
	}

	builder := strings.Builder{}

	if errorTemplate == nil {
		initError()
	}

	err = errorTemplate.Execute(&builder, data)
	if err != nil {
		WriteJson(w, err)
		return
	}

	WriteHtml(w, builder.String())
}
