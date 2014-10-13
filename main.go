package main

import "fmt"
import "regexp"
import "net/http"
import "io/ioutil"
import "net/url"
import "log"
import "text/template"

type TorrentLink struct {
	Magnet string
	Title  string
	Show   Show
}

type Show struct {
	Id    int
	Title string
	Slug  string
}

var links []TorrentLink
var shows []Show

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	re := regexp.MustCompile(`"(magnet.+:[0-9]+)"`)

	links = make([]TorrentLink, 0)
	shows = make([]Show, 0)

	shows = append(shows, Show{23, "The Big Bang Theory", "the-big-bang-theory"})
	shows = append(shows, Show{330, "Modern Family", "modern-family"})

	for _, show := range shows {
		uri := fmt.Sprintf("https://eztv.it/shows/%v/%v/", show.Id, show.Slug)
		resp, _ := http.Get(uri)
		contents, _ := ioutil.ReadAll(resp.Body)

		magnets := (re.FindAllString(string(contents), -1))
		for _, magnet := range magnets {
			fmt.Println(magnet)
			u, _ := url.Parse(magnet)
			title := u.Query()["dn"][0]
			links = append(links, TorrentLink{magnet, title, show})
		}
	}

	http.Handle("/", http.HandlerFunc(EZTV))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	fmt.Println(links)
}

func EZTV(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, links)
}

const templateStr = `
<html>
<head>
<title>EZTV links</title>
</head>
<body>
<ul>{{ range .}}
	<li><a href={{ .Magnet }}>{{ .Show.Slug }} - {{ .Title }}</a></li>
{{ end }}
</ul>
</body>
</html>
`
