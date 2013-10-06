package hello

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"appengine"
	"appengine/urlfetch"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/fetch", fetch)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	username := ""
	if u != nil {
		username = u.String()
	}
	if err := guestbookTemplate.Execute(w, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
  <body>
    <div>User: {{.}}</div>
    <form action="/fetch" method="get">
      <div><input type="url" name="url"/></div>
      <div><input type="submit" value="Fetch URL"></div>
    </form>
  </body>
</html>
`

func fetch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url := r.FormValue("url")
	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, v := range resp.Header {
		for _, v2 := range v {
			w.Header().Add(k, v2);
		}
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, string(body[:]))
}
