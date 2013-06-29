package gemBadges

import (
    "net/http"
    "fmt"
    "appengine"
    "appengine/urlfetch"
//    "os"
    "log"
    "image"
    "image/draw"
    "image/color"
    "image/png"
    "io/ioutil"
    "strconv"
    "github.com/emrekutlu/go-rubygems"
    "github.com/gorilla/mux"
    "code.google.com/p/freetype-go/freetype"
)

func init() {
    r := mux.NewRouter()

    r.HandleFunc("/", root)
    gemsRouter := r.PathPrefix("/gems/{gem}").Subrouter()
    gemsRouter.HandleFunc("/", gems)
    gemsRouter.HandleFunc("/downloads.png", downloads)

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
    w.Header().Set("Content-type", "image/png")
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    vars := mux.Vars(r)

    rubygems.Initialize(client)
    gem, err := rubygems.NewGem(vars["gem"]).Get()
/*
    file, err := os.Open("static/badge.png")
    if err != nil {
      log.Fatal(err)
    } else {
      defer file.Close()
    }

    // Decode the image.
    m, _, err := image.Decode(file)
   if err != nil {
     log.Fatal(err)
   }
*/
    m := image.NewRGBA(image.Rect(0, 0, 100, 20))
    blue := color.RGBA{0, 29, 204, 170}
    draw.Draw(m, m.Bounds(), &image.Uniform{C: blue}, image.ZP, draw.Src)

    // Freetype

// Read the font data.
   fontBytes, err := ioutil.ReadFile("Signika-Regular.ttf")
   if err != nil {
     log.Println(err)
     return
   }
   font, err := freetype.ParseFont(fontBytes)
   if err != nil {
     log.Println(err)
     return
   }

    var size float64
    size = 12
    freetypeContext := freetype.NewContext()
    freetypeContext.SetDPI(72)
    freetypeContext.SetFont(font)
    freetypeContext.SetFontSize(size)
    freetypeContext.SetClip(m.Bounds())
    freetypeContext.SetDst(m)
    freetypeContext.SetSrc(image.White)

    pt := freetype.Pt(40, int(freetypeContext.PointToFix32(size)>>8))

   for _, s := range []string{strconv.Itoa(gem.Downloads)} {
     _, err = freetypeContext.DrawString(s, pt)
     if err != nil {
       log.Println(err)
       return
     }
     pt.Y += freetypeContext.PointToFix32(size * 1.5)
   }


    if err != nil {
      fmt.Fprint(w, err)
    } else {
      log.Println(gem.Downloads)
      png.Encode(w, m)
    }
}

