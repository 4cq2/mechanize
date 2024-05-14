package tubi

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (c Content) Marshal() ([]byte, error) {
   return json.Marshal(c)
}

type Content struct {
   Children        []*Content
   Detailed_Type   string
   Episode_Number  int `json:",string"`
   ID              int `json:",string"`
   Series_ID       int `json:",string"`
   Title           string
   Video_Resources []VideoResource
   Year            int
   parent          *Content
}

func (c Content) Episode() bool {
   return c.Detailed_Type == "episode"
}

func (c Content) Get(id int) (*Content, bool) {
   if c.ID == id {
      return &c, true
   }
   for _, child := range c.Children {
      if v, ok := child.Get(id); ok {
         return v, true
      }
   }
   return nil, false
}

func (c *Content) set(parent *Content) {
   c.parent = parent
   for _, child := range c.Children {
      child.set(c)
   }
}

type Namer struct {
   C *Content
}

func (n Namer) Episode() int {
   return n.C.Episode_Number
}

func (n Namer) Season() int {
   if v := n.C.parent; v != nil {
      return v.ID
   }
   return 0
}

func (n Namer) Show() string {
   if v := n.C.parent; v != nil {
      return v.parent.Title
   }
   return ""
}

// S01:E03 - Hell Hath No Fury
func (n Namer) Title() string {
   if _, v, ok := strings.Cut(n.C.Title, " - "); ok {
      return v
   }
   return n.C.Title
}

func (n Namer) Year() int {
   return n.C.Year
}

func (c *Content) New(id int) error {
   req, err := http.NewRequest("GET", "https://uapi.adrise.tv/cms/content", nil)
   if err != nil {
      return err
   }
   req.URL.RawQuery = url.Values{
      "content_id": {strconv.Itoa(id)},
      "deviceId":   {"!"},
      "platform":   {"android"},
      "video_resources[]": {
         "dash",
         "dash_widevine",
      },
   }.Encode()
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   text, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   return c.Unmarshal(text)
}

func (c *Content) Unmarshal(text []byte) error {
   err := json.Unmarshal(text, c)
   if err != nil {
      return err
   }
   c.set(nil)
   return nil
}
