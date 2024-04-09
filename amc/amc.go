package amc

import (
   "encoding/base64"
   "encoding/json"
   "errors"
   "net/http"
   "strconv"
   "strings"
)

func cache_hash() string {
   return base64.StdEncoding.EncodeToString([]byte("ff="))
}

type ContentCompiler struct {
   Data	struct {
      Children []struct {
         Properties json.RawMessage
         Type string
      }
   }
}

func (c ContentCompiler) Video() (*CurrentVideo, error) {
   for _, child := range c.Data.Children {
      if child.Type == "video-player-ap" {
         var s struct {
            CurrentVideo CurrentVideo
         }
         err := json.Unmarshal(child.Properties, &s)
         if err != nil {
            return nil, err
         }
         return &s.CurrentVideo, nil
      }
   }
   return nil, errors.New("video-player-ap")
}

type CurrentVideo struct {
   Meta struct {
      Airdate string // 1996-01-01T00:00:00.000Z
      EpisodeNumber int
      Season int `json:",string"`
      ShowTitle string
   }
   Text struct {
      Title string
   }
}

func (c CurrentVideo) Episode() int {
   return c.Meta.EpisodeNumber
}

func (c CurrentVideo) Show() string {
   return c.Meta.ShowTitle
}

func (c CurrentVideo) Season() int {
   return c.Meta.Season
}

func (c CurrentVideo) Title() string {
   return c.Text.Title
}

func (c CurrentVideo) Year() int {
   if v, _, ok := strings.Cut(c.Meta.Airdate, "-"); ok {
      if v, err := strconv.Atoi(v); err == nil {
         return v
      }
   }
   return 0
}

type DataSource struct {
   Key_Systems *struct {
      Widevine struct {
         License_URL string
      } `json:"com.widevine.alpha"`
   }
   Src string
   Type string
}

type Playback struct {
   header http.Header
   body struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []DataSource
         }
      }
   }
}

func (p Playback) HttpsDash() (*DataSource, bool) {
   for _, s := range p.body.Data.PlaybackJsonData.Sources {
      if strings.HasPrefix(s.Src, "https://") {
         if s.Type == "application/dash+xml" {
            return &s, true
         }
      }
   }
   return nil, false
}

func (Playback) RequestBody(b []byte) ([]byte, error) {
   return b, nil
}

func (p Playback) RequestHeader() (http.Header, error) {
   h := make(http.Header)
   h.Set("bcov-auth", p.header.Get("X-AMCN-BC-JWT"))
   return h, nil
}

func (p Playback) RequestUrl() (string, bool) {
   if v, ok := p.HttpsDash(); ok {
      return v.Key_Systems.Widevine.License_URL, true
   }
   return "", false
}

func (Playback) ResponseBody(b []byte) ([]byte, error) {
   return b, nil
}

type WebAddress struct {
   NID string
   Path string
}

func (w *WebAddress) Set(s string) error {
   var found bool
   _, w.Path, found = strings.Cut(s, "amcplus.com")
   if !found {
      return errors.New("amcplus.com")
   }
   _, w.NID, found = strings.Cut(w.Path, "--")
   if !found {
      return errors.New("--")
   }
   return nil
}

func (w WebAddress) String() string {
   return w.Path
}
