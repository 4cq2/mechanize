package max

import (
   "encoding/json"
   "net/http"
   "strings"
   "time"
)

func (a *address) UnmarshalText(text []byte) error {
   split := strings.Split(string(text), "/")
   a.video_id = split[3]
   a.edit_id = split[4]
   return nil
}

func (d default_token) routes(path string) (*default_routes, error) {
   req, err := http.NewRequest(
      "", "https://default.any-amer.prd.api.discomax.com/cms/routes"+path, nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = "include=default"
   req.Header.Set("authorization", "Bearer " + d.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   route := new(default_routes)
   err = json.NewDecoder(resp.Body).Decode(route)
   if err != nil {
      return nil, err
   }
   return route, nil
}

type default_routes struct {
   Data struct {
      Attributes struct {
         Url address
      }
   }
   Included []include
}

type include struct {
   Attributes struct {
      AirDate time.Time
      EpisodeNumber int
      Name string
      SeasonNumber int
      Type string
   }
   Id string
}

type address struct {
   video_id string
   edit_id string
}

func (d default_routes) Show() string {
   return ""
}

func (default_routes) Season() int {
   return 0
}

func (default_routes) Episode() int {
   return 0
}

func (default_routes) Title() string {
   return ""
}

func (default_routes) Year() int {
   return 0
}
