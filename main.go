package gemBadges

import (
    "net/http"
    "fmt"
    "appengine"
    "appengine/urlfetch"
    "github.com/emrekutlu/go-rubygems"
    "github.com/gorilla/mux"
)

func init() {
    r := mux.NewRouter()

    r.HandleFunc("/", root)
    gemsRouter := r.PathPrefix("/gems/{gem}").Subrouter()
    gemsRouter.HandleFunc("/", gems)
    gemsRouter.HandleFunc("/downloads", downloads)

    http.Handle("/", r)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "ruby badges!")
}

func gems(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    vars := mux.Vars(r)

    rubygems.Initialize(client)
    gem, err := rubygems.NewGem(vars["gem"]).Get()

    if err != nil {
      fmt.Fprint(w, err)
    } else {
      fmt.Fprint(w, gem)
    }
}

func downloads(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    vars := mux.Vars(r)

    rubygems.Initialize(client)
    gem, err := rubygems.NewGem(vars["gem"]).Get()

    if err != nil {
      fmt.Fprint(w, err)
    } else {
      fmt.Fprint(w, gem.Downloads)
    }
}

