package main

import(
	"flag"
	"os"
	"log"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
	"fmt"
)

var logger *log.Logger
func main() {
	var u = flag.String("url", "", "Url to fetch.")
	flag.Parse() // parse the flags
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	//errorlog, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
		//logger.Printf("error opening file: %v", err)
		//os.Exit(1)
	//}
	//defer errorlog.Close()
	//logger = log.New(errorlog, "applog: ", log.Lshortfile|log.LstdFlags)
	logger = log.New(os.Stderr, "", 0)

	get := func(url string) error {
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		//parts := strings.Split(url, "/")
		//last := parts[len(parts)-1]
		//file, err := os.OpenFile(last, os.O_CREATE|os.O_WRONLY, 0644)
		//if err != nil {
		//	return err
		//}
		//defer file.Close()
		//_, err = io.Copy(os.Stdout, response.Body)
		//n, err := io.Copy(file, resp.Body) // first var shows number of bytes
		if err != nil {
			return err
		}
		
		z := html.NewTokenizer(response.Body)

		for {
			tt := z.Next()

			switch {
				case tt == html.ErrorToken:
					// End of the document, we're done
				        return nil
				case tt == html.StartTagToken:
				        t := z.Token()

				   	isUl := t.Data == "ul"
				        if isUl {
						tn := z.Next()
						tul := z.Token()
						if tn == html.StartTagToken && tul.Data == "li" {
							ta := z.Next()
							tli := z.Token()
							if ta == html.StartTagToken && tli.Data == "a" {
								for _, a := range tli.Attr {
									if a.Key == "title" {
									        fmt.Println("Found title:", a.Val)
									        break
									}
								}
							}
						}
				        }
			    }
		}
		return nil
	}
	//validate flag input
	if _, err := url.Parse(*u); err != nil{
		logger.Fatalf("Invalid input url: %v \n", err)
	}
	
	err := get(*u)
	if err != nil {
		logger.Fatalf("Falied request: %s \n", err)
	}
}
