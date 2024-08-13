package draken

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strings"
)

type Namer struct {
   Movie *FullMovie
}

func (Namer) Episode() int {
   return 0
}

func (Namer) Season() int {
   return 0
}

func (Namer) Show() string {
   return ""
}

func (FullMovie) Error() string {
   return "FullMovie"
}

type FullMovie struct {
   DefaultPlayable struct {
      Id string
   }
   ProductionYear int `json:",string"`
   Title          string
}

const get_custom_id = `
query($customId: ID!) {
   viewer {
      viewableCustomId(customId: $customId) {
         ... on Movie {
            defaultPlayable {
               id
            }
            productionYear
            title
         }
      }
   }
}
`

func graphql_compact(s string) string {
   f := strings.Fields(s)
   return strings.Join(f, " ")
}

type Entitlement struct {
   Token string
}

type header struct {
   key   string
   value string
}

var magine_accesstoken = header{
   "magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663",
}

var magine_play_devicemodel = header{
   "magine-play-devicemodel", "firefox 111.0 / windows 10",
}

var magine_play_deviceplatform = header{
   "magine-play-deviceplatform", "firefox",
}

var magine_play_devicetype = header{
   "magine-play-devicetype", "web",
}

var magine_play_drm = header{
   "magine-play-drm", "widevine",
}

var magine_play_protocol = header{
   "magine-play-protocol", "dashs",
}

// this value is important, with the wrong value you get random failures
var x_forwarded_for = header{
   "x-forwarded-for", "95.192.0.0",
}

type Playback struct {
   Headers  map[string]string
   Playlist string
}

func (Poster) RequestUrl() (string, bool) {
   return "https://client-api.magine.com/api/playback/v1/widevine/license", true
}

func (Poster) UnwrapResponse(b []byte) ([]byte, error) {
   return b, nil
}

func (Poster) WrapRequest(b []byte) ([]byte, error) {
   return b, nil
}

type Poster struct {
   Login *AuthLogin
   Play *Playback
}

func (h *header) set(head http.Header) {
   head.Set(h.key, h.value)
}

func (p *Poster) RequestHeader() (http.Header, error) {
   head := http.Header{}
   magine_accesstoken.set(head)
   head.Set("authorization", "Bearer "+p.Login.Token)
   for key, value := range p.Play.Headers {
      head.Set(key, value)
   }
   return head, nil
}

func (f *FullMovie) New(custom_id string) error {
   body, err := func() ([]byte, error) {
      var s struct {
         Query     string `json:"query"`
         Variables struct {
            CustomId string `json:"customId"`
         } `json:"variables"`
      }
      s.Variables.CustomId = custom_id
      s.Query = graphql_compact(get_custom_id)
      return json.MarshalIndent(s, "", " ")
   }()
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/apiql/v2",
      bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   magine_accesstoken.set(req.Header)
   x_forwarded_for.set(req.Header)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var data struct {
      Data struct {
         Viewer struct {
            ViewableCustomId *FullMovie
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&data)
   if err != nil {
      return err
   }
   if v := data.Data.Viewer.ViewableCustomId; v != nil {
      *f = *v
      return nil
   }
   return FullMovie{}
}

func (n *Namer) Title() string {
   return n.Movie.Title
}

func (n *Namer) Year() int {
   return n.Movie.ProductionYear
}
