package application

import (
	"net/http"
)

type ContextMenuData struct {
	Id   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Data any    `json:"data"`
}

const (
	ContextMenuOpen = 0
)

var contextmenuMethodNames = map[int]string{
	ContextMenuOpen: "Open",
}

func (m *MessageProcessor) processContextMenuMethod(method int, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case ContextMenuOpen:
		var data ContextMenuData
		err := params.ToStruct(&data)
		if err != nil {
			m.httpError(rw, "error parsing contextmenu message: %s", err.Error())
			return
		}
		window.openContextMenu(&data)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown contextmenu method: %d", method)
	}

	m.Info("Runtime:", "method", "ContextMenu."+contextmenuMethodNames[method])

}
