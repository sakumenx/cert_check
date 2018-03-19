package main

import(
  "net/http"
  "fmt"
  "bytes"
  "os"
  "time"
)

func main() {
  if len(os.Args) != 3 {
    fmt.Print("track not specified!\n")
    os.Exit(1)
  }
  var url string = os.Args[1]

  var webhook string = os.Args[2]

  resp, _ := http.Get(url)

  expireUTCTime := resp.TLS.PeerCertificates[0].NotAfter
  expireJSTTime := expireUTCTime.In(time.FixedZone("Asia/Tokyo", 9 * 60 * 60))

  fmt.Print(expireJSTTime)
  jsonStr := `{"text":"` + url + `\nSSL証明書期限は` + expireJSTTime.Format("2006/01/02  15:04") + `までです。"}`
  fmt.Print(jsonStr)
  req, _ := http.NewRequest(
    "POST",
    webhook,
    bytes.NewBuffer([]byte(jsonStr)),
  )
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Print(err)
  }
  defer resp.Body.Close()
}
