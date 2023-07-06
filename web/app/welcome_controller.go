package app

import (
	"github.com/andrewarrow/feedback/router"
)

func HandleWelcome(c *router.Context, second, third string) {
	c.Title = "Feedback - Go on Rails"
	if second == "" && third == "" && c.Method == "GET" {
		handleWelcomeIndex(c)
		return
	}
	c.NotFound = true
}

func handleWelcomeIndex(c *router.Context) {
	list := []any{}
	feature := map[string]any{"name": "CRUD", "desc": "Create, retrieve, update, delete"}
	list = append(list, feature)
	feature = map[string]any{"name": "Turbo/Ajax", "desc": "Default is to replace a div via XMLHttpRequest."}
	list = append(list, feature)
	feature = map[string]any{"name": "HTML/CSS", "desc": "Default is to render application_layout.html with tailwindcss."}
	list = append(list, feature)
	feature = map[string]any{"name": "API/JSON", "desc": "Default is to return paginated lists of size 30."}
	list = append(list, feature)
	feature = map[string]any{"name": "Migrations", "desc": "Default is to run additive migrations on startup. Create table, add column, create index (CONCURRENTLY) but never drop or rename anything."}
	list = append(list, feature)
	feature = map[string]any{"name": "Schema", "desc": "There is one source of truth file called feedback.json that defines all models/tables."}
	list = append(list, feature)

	colAttributes := map[int]string{}
	colAttributes[0] = "w-1/4"
	colAttributes[1] = "w-3/4"

	m := map[string]any{}
	headers := []string{"feature", "description"}

	params := map[string]any{}
	m["headers"] = headers
	m["cells"] = c.MakeCells(list, headers, params, "_feedback")
	m["col_attributes"] = colAttributes

	send := map[string]any{}
	send["bottom"] = c.Template("table_show.html", m)
	c.SendContentInLayout("feedback.html", send, 200)
}
