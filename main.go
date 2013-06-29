package gemBadges

import (
    "net/http"
    "net/url"
    "appengine"
    "appengine/urlfetch"
    "fmt"
    "log"
    "io/ioutil"
    "image"
    "image/draw"
    "image/color"
    "image/png"
    "strconv"
    "encoding/hex"
    "github.com/emrekutlu/go-rubygems"
    "github.com/gorilla/mux"
    "code.google.com/p/freetype-go/freetype"
)

var(
  rightBg = color.RGBA{53, 187, 15, 255}
  errorBg = color.RGBA{255, 0, 0, 255}
  leftBg  = color.RGBA{63, 63, 63, 255}
)

func init() {
    r := mux.NewRouter()

    r.HandleFunc("/", root)
    gemsRouter := r.PathPrefix("/gems/{gem}").Subrouter()
    gemsRouter.HandleFunc("/", gems)
    gemsRouter.HandleFunc("/downloads.png", downloads)
    gemsRouter.HandleFunc("/version.png", version)

    http.Handle("/", r)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "code badges!")
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

    checkParams(r.URL.Query())

    rubygems.Initialize(client)
    gem, err := rubygems.NewGem(vars["gem"]).Get()

    textImage := image.NewRGBA(image.Rect(0, 0, 64, 18))
    addTextAndBg("downloads", textImage, &leftBg)

    var count string
    var countImage *image.RGBA
    if err == nil {
      count = strconv.Itoa(gem.Downloads)
      countImage = image.NewRGBA(image.Rect(0, 0, 4 + len(count) * 7, 18))
      addTextAndBg(count, countImage, &rightBg)
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

func version(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "image/png")
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    vars := mux.Vars(r)

    checkParams(r.URL.Query())

    rubygems.Initialize(client)
    gem, err := rubygems.NewGem(vars["gem"]).Get()

    textImage := image.NewRGBA(image.Rect(0, 0, 47, 18))
    addTextAndBg("version", textImage, &leftBg)

    var version string
    var versionImage *image.RGBA
    if err == nil {
      version = gem.Version
      versionImage = image.NewRGBA(image.Rect(0, 0, 4 + len(version) * 6, 18))
      addTextAndBg(version, versionImage, &rightBg)
    } else {
      version = "Gem not found!"
      versionImage = image.NewRGBA(image.Rect(0, 0, 87, 18))
      addTextAndBg(version, versionImage, &errorBg)
    }

    bgImage := image.NewRGBA(image.Rect(0, 0, textImage.Bounds().Dx() + versionImage.Bounds().Dx(), 18))
    draw.Draw(bgImage, textImage.Bounds(), textImage, image.ZP, draw.Over)
    draw.Draw(bgImage, bgImage.Bounds(), versionImage, image.Point{-textImage.Bounds().Dx(), 0}, draw.Over)

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

func hexToRGBA(hexParam string) (color.RGBA, error){
  rgbaArray, err := hex.DecodeString(hexParam)
  if err == nil {
    return color.RGBA{rgbaArray[0], rgbaArray[1], rgbaArray[2], 255} , err
  } else {
    return color.RGBA{0, 0, 0, 255}, err
  }
}

func checkParams(params url.Values) {
  leftBgParam := params.Get("left_bg")
  if len(leftBgParam) > 0 {
    newLeftBg, err := hexToRGBA(leftBgParam)
    if err == nil {
      leftBg = newLeftBg
    }
  }

  rightBgParam := params.Get("right_bg")
  if len(rightBgParam) > 0 {
    newRightBg, err := hexToRGBA(rightBgParam)
    if err == nil {
      rightBg = newRightBg
    }
  }
}
