package v2

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/shurcooL/httpgzip"
	"github.com/swaggest/swgui"
)

var staticServer = httpgzip.FileServer(assets, httpgzip.FileServerOptions{})

// Handler handle swagger UI request
type Handler struct {
	swgui.Config

	ConfigJson template.JS

	tpl          *template.Template
	staticServer http.Handler
}

// NewHandler returns a HTTP handler for swagger UI
func NewHandler(title, swaggerJSONPath string, basePath string) *Handler {
	h := &Handler{}
	h.Title = title
	h.SwaggerJSON = swaggerJSONPath
	h.BasePath = basePath

	j, _ := json.Marshal(h.Config)
	h.ConfigJson = template.JS(j)
	h.tpl, _ = template.New("index").Parse(indexTpl)
	h.staticServer = http.StripPrefix(basePath, staticServer)
	return h
}

// NewHandlerWithConfig returns a HTTP handler for swagger UI
func NewHandlerWithConfig(config swgui.Config) *Handler {
	h := &Handler{
		Config: config,
	}
	j, _ := json.Marshal(h.Config)
	h.ConfigJson = template.JS(j)
	h.tpl, _ = template.New("index").Parse(indexTpl)
	h.staticServer = http.StripPrefix(h.BasePath, staticServer)
	return h
}

// ServeHTTP implement http.Handler interface, to handle swagger UI request
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != h.BasePath {
		h.staticServer.ServeHTTP(w, r)
		return
	}

	if err := h.tpl.Execute(w, h); err != nil {
		http.NotFound(w, r)
	}
}
