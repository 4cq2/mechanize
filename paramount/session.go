package paramount

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/hex"
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

/*
older versions fail
rg cbs.?appsecret
*/
var app_secrets = map[string]string{
   "04.1.0": "118b561316186716",
   "04.1.5": "",
   "04.1.8": "4c5bafd363ca317a",
   "04.3.0": "f33bfa06992390a5",
   "04.3.1": "",
   "04.3.5": "",
   "04.6.0": "",
   "04.8.6": "a958002817953588",
   "07.3.58": "4be3d46aecbcd26d",
   "08.0.00": "5d1d865f929d3daa",
   "08.0.10": "",
   "08.0.12": "",
   "08.0.14": "",
   "08.0.16": "",
   "08.0.20": "",
   "08.0.22": "",
   "08.0.24": "",
   "08.0.26": "",
   "08.0.28": "",
   "08.0.30": "",
   "08.0.32": "",
   "08.0.34": "",
   "08.0.36": "",
   "08.0.38": "",
   "08.0.40": "",
   "08.0.42": "",
   "08.0.44": "",
   "08.0.46": "",
   "08.0.48": "",
   "08.0.50": "",
   "08.0.52": "",
   "08.0.54": "",
   "08.1.10": "",
   "08.1.12": "",
   "08.1.14": "",
   "08.1.16": "add603b54be2fc3c",
   "08.1.18": "1705732089ff4d60",
   "08.1.20": "817774cbafb2b797",
   "08.1.22": "ddca2f16bfa3d937",
   "08.1.23": "c0966212aa651e8b",
   "08.1.26": "a75bd3a39bfcbc77",
   "08.1.28": "d0795c0dffebea73",
   "12.0.26": "f012987182d6f16c",
   "12.0.27": "79b7e56e442e65ed",
   "12.0.28": "439ba2e3622c344a",
}

func session_secret(guid, secret string) (*Session, error) {
   token, err := new_token(secret)
   if err != nil {
      return nil, err
   }
   var buf strings.Builder
   buf.WriteString("https://www.paramountplus.com/apps-api/v3.0/androidphone")
   buf.WriteString("/irdeto-control/anonymous-session-token.json")
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = "at=" + url.QueryEscape(token)
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   sess := new(Session)
   if err := json.NewDecoder(res.Body).Decode(sess); err != nil {
      return nil, err
   }
   sess.URL += guid
   return sess, nil
}

func new_token(secret string) (string, error) {
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   var src []byte
   src = append(src, '|')
   src = append(src, secret...)
   src = pad(src)
   var iv [aes.BlockSize]byte
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(src, src)
   var dst []byte
   dst = append(dst, 0, aes.BlockSize)
   dst = append(dst, iv[:]...)
   dst = append(dst, src...)
   return base64.StdEncoding.EncodeToString(dst), nil
}

func New_Session(guid string) (*Session, error) {
   return session_secret(guid, app_secrets["12.0.28"])
}

type Session struct {
   URL string
   LS_Session string
}

func (s Session) Request_URL() string {
   return s.URL
}

func (s Session) Request_Header() http.Header {
   head := make(http.Header)
   head.Set("Authorization", "Bearer " + s.LS_Session)
   return head
}

func (Session) Request_Body(buf []byte) ([]byte, error) {
   return buf, nil
}

func (Session) Response_Body(buf []byte) ([]byte, error) {
   return buf, nil
}
