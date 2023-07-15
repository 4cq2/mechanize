package main

import (
   "154.pages.dev/media/youtube"
   "flag"
   "strings"
)

type flags struct {
   audio_q string
   audio_t string
   info bool
   r youtube.Request
   refresh bool
   request int
   video_q string
   video_t string
}

func main() {
   var f flags
   // a
   flag.Var(&f.r, "a", "address")
   // b
   flag.StringVar(&f.r.Video_ID, "b", "", "video ID")
   // aq
   flag.StringVar(&f.audio_q, "aq", "AUDIO_QUALITY_MEDIUM", "audio quality")
   // at
   flag.StringVar(&f.audio_t, "at", "opus", "audio type")
   // i
   flag.BoolVar(&f.info, "i", false, "information")
   // r
   {
      var b strings.Builder
      b.WriteString("0: Android\n")
      b.WriteString("1: Android embed\n")
      b.WriteString("2: Android check")
      flag.IntVar(&f.request, "r", 0, b.String())
   }
   // refresh
   flag.BoolVar(&f.refresh, "refresh", false, "create OAuth refresh token")
   // vq
   flag.StringVar(&f.video_q, "vq", "1080p", "video quality")
   // vt
   flag.StringVar(&f.video_t, "vt", "vp9", "video type")
   flag.Parse()
   if f.refresh {
      err := f.do_refresh()
      if err != nil {
         panic(err)
      }
   } else if f.r.Video_ID != "" {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
