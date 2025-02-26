package handler

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"text/template"

	"wartracker/ui/wartracker-ui/site"
)

type Handler struct {
}

var SiteDataDir string

func (h Handler) Default(w http.ResponseWriter, r *http.Request) {
	var sfile string

	p := path.Clean(r.URL.Path)
	dir, tmpl := path.Split(p)
	fmt.Println(dir)
	fmt.Println(tmpl)
	if tmpl == "" {
		if dir == "/" {
			sfile = SiteDataDir + "root.json"
		}
		tmpl = "root"
	} else if dir == "/" {
		sfile = SiteDataDir + tmpl + ".json"
	} else {
		rdir := strings.Split(dir, "/")[1]
		fmt.Println(rdir)
		sfile = SiteDataDir + rdir + ".json"
	}
	tmpl += ".gohtml"

	s := site.NewSite(sfile)

	//	w.Write([]byte(build.Build + "\n"))
	//	state := fmt.Sprintf("Dir: %s\nTemplate: %s\nSite:\n%#v\n", dir, tmpl, s)
	//	w.Write([]byte(state))

	RenderTemplate(w, tmpl, s)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, s *site.Site) {
	t, err := template.ParseFiles(
		"templates/head.gohtml",
		"templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = t.ExecuteTemplate(w, tmpl, s)
	if err != nil {
		return
	}
}
