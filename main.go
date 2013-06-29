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

var(
  countBg = color.RGBA{53, 187, 15, 255}
  errorBg = color.RGBA{255, 0, 0, 255}
  textBg  = color.RGBA{63, 63, 63, 255}
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

    textImage := image.NewRGBA(image.Rect(0, 0, 64, 18))
    addTextAndBg("downloads", textImage, &textBg)

    var count string
    var countImage *image.RGBA
    if err == nil {
      count = strconv.Itoa(gem.Downloads)
      countImage = image.NewRGBA(image.Rect(0, 0, 4 + len(count) * 7, 18))
      addTextAndBg(count, countImage, &countBg)
    } else {
      count = "Gem not found!"
      countImage = image.NewRGBA(image.Rect(0, 0, 87, 18))
      addTextAndBg(count, countImage, &errorBg)
    }

    bgImage := image.NewRGBA(image.Rect(0, 0, textImage.Bounds().Dx() + countImage.Bounds().Dx(), 18))
    draw.Draw(bgImage, textImage.Bounds(), textImage, image.ZP, draw.Over)
    draw.Draw(bgImage, bgImage.Bounds(), countImage, image.Point{-textImage.Bounds().Dx(), 0}, draw.Over)

    png.Encode(w, bgImage)
}

func addTextAndBg(txt string, img *image.RGBA, bgColor *color.RGBA) {

  draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.ZP, draw.Src)

   fontBytes, err := ioutil.ReadFile("Signika-Regular.ttf")
   if err != nil {
     log.Println(err)
   }
   font, err := freetype.ParseFont(fontBytes)
   if err != nil {
     log.Println(err)
   }

    var size float64
    size = 12
    freetypeContext := freetype.NewContext()
    freetypeContext.SetDPI(72)
    freetypeContext.SetFont(font)
    freetypeContext.SetFontSize(size)
    freetypeContext.SetClip(img.Bounds())
    freetypeContext.SetDst(img)
    freetypeContext.SetSrc(image.White)

    pt := freetype.Pt(4, int(freetypeContext.PointToFix32(size)>>8))

   _, err = freetypeContext.DrawString(txt, pt)
   if err != nil {
     log.Println(err)
   }
}
