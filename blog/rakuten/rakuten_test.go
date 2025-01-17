package rakuten

import (
   "41.neocities.org/platform/mullvad"
   "41.neocities.org/text"
   "log"
   "net/http"
   "testing"
)

type transport struct{}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Print(req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type web_test struct {
   a           address
   address     string
   address_out string
   content_id  string
   key_id      string
   language    string
   location    string
   name        string
}

var web_tests = []web_test{
   {
      a: address{
         market_code: "cz", content_id: "transvulcania-the-people-s-run",
      },
      language:    "SPA",
      address:     "rakuten.tv/cz/movies/transvulcania-the-people-s-run",
      address_out: "cz/movies/transvulcania-the-people-s-run",
      name:        "Transvulcania, The People’s Run - 2024",
      location:    "cz",
   },
   {
      content_id:  "MGU1MTgwMDA2Y2Q1MDhlZWMwMGQ1MzVmZWM2YzQyMGQtbWMtMC0xNDEtMC0w",
      key_id:      "DlGAAGzVCO7ADVNf7GxCDQ==",
      address:     "rakuten.tv/fr/movies/infidele",
      language:    "ENG",
      address_out: "fr/movies/infidele",
      a: address{
         market_code: "fr", content_id: "infidele",
      },
      name:     "Infidèle - 2002",
      location: "fr",
   },
   {
      content_id:  "OWE1MzRhMWYxMmQ2OGUxYTIzNTlmMzg3MTBmZGRiNjUtbWMtMC0xNDctMC0w",
      key_id:      "mlNKHxLWjhojWfOHEP3bZQ==",
      language:    "ENG",
      address:     "rakuten.tv/se/movies/i-heart-huckabees",
      address_out: "se/movies/i-heart-huckabees",
      a: address{
         market_code: "se", content_id: "i-heart-huckabees",
      },
      name:     "I Heart Huckabees - 2004",
      location: "se",
   },
   {
      a: address{
         market_code: "uk",
         season_id:   "hell-s-kitchen-usa-15",
         content_id:  "hell-s-kitchen-usa-15-1",
      },
      language:    "ENG",
      address:     "rakuten.tv/uk/player/episodes/stream/hell-s-kitchen-usa-15/hell-s-kitchen-usa-15-1",
      address_out: "uk/player/episodes/stream/hell-s-kitchen-usa-15/hell-s-kitchen-usa-15-1",
      name:        "Hell's Kitchen USA - 15 1 - 18 Chefs Compete",
      location:    "gb",
   },
}

func TestAddress(t *testing.T) {
   for _, test := range web_tests {
      t.Run("Set", func(t *testing.T) {
         var a address
         err := a.Set(test.address)
         if err != nil {
            t.Fatal(err)
         }
         if a != test.a {
            t.Fatal(test)
         }
      })
      t.Run("String", func(t *testing.T) {
         if test.a.String() != test.address_out {
            t.Fatal(test.a)
         }
      })
   }
   t.Run("classification_id", func(t *testing.T) {
      var web address
      _, ok := web.classification_id()
      if ok {
         t.Fatal(web)
      }
   })
}

func TestMain(m *testing.M) {
   http.DefaultClient.Transport = transport{}
   m.Run()
}

func TestNamer(t *testing.T) {
   for _, test := range web_tests {
      class, _ := test.a.classification_id()
      var content *gizmo_content
      if test.a.season_id != "" {
         season, err := test.a.season(class)
         if err != nil {
            t.Fatal(err)
         }
         _, ok := season.content(&address{})
         if ok {
            t.Fatal(season)
         }
         content, _ = season.content(&test.a)
      } else {
         var err error
         content, err = test.a.movie(class)
         if err != nil {
            t.Fatal(err)
         }
      }
      if text.Name(namer{content}) != test.name {
         t.Fatal(content)
      }
   }
}

func TestContent(t *testing.T) {
   for _, test := range web_tests {
      class, _ := test.a.classification_id()
      var content *gizmo_content
      if test.a.season_id != "" {
         season, err := test.a.season(class)
         if err != nil {
            t.Fatal(err)
         }
         _, ok := season.content(&address{})
         if ok {
            t.Fatal(season)
         }
         content, _ = season.content(&test.a)
      } else {
         var err error
         content, err = test.a.movie(class)
         if err != nil {
            t.Fatal(err)
         }
      }
      t.Run("String", func(t *testing.T) {
         if content.String() == "" {
            t.Fatal(content)
         }
      })
      func() {
         err := mullvad.Connect(test.location)
         if err != nil {
            t.Fatal(err)
         }
         defer mullvad.Disconnect()
         t.Run("fhd", func(t *testing.T) {
            _, err = content.fhd(class, test.language).streamings()
            if err != nil {
               t.Fatal(err)
            }
         })
         t.Run("hd", func(t *testing.T) {
            _, err = content.hd(class, test.language).streamings()
            if err != nil {
               t.Fatal(err)
            }
         })
      }()
   }
}
