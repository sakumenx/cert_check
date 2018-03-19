package main

import(
  "net/http"
  "fmt"
  "bytes"
  "os"
  "time"
  "strconv"
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

  currentUTCTime := time.Now()
  currentJSTTime := currentUTCTime.In(time.FixedZone("Asia/Tokyo", 9 * 60 * 60))

  timeSub := expireJSTTime.Sub(currentJSTTime)
  daysSub := int(timeSub.Hours()) / 24
  fmt.Println(daysSub)

  fmt.Print(expireJSTTime)
  postText := url + "\nSSL証明書期限は" + expireJSTTime.Format("2006/01/02 15:04") + "までです。"
  if daysSub < 60 {
    postText += " `あと残り" + strconv.Itoa(daysSub) + "日`"
  }
  jsonStr := "{\"text\":\"" + postText + "\"}"
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
