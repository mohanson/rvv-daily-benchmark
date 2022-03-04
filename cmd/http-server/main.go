package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/godump/acdb"
	"github.com/godump/doa"
)

type Conf struct {
	VMPath string
	VM     string
	Benchs []string
	URLs   []string
}

var (
	cDb   = acdb.Doc("./db")
	cConf = func() Conf {
		conf := Conf{}
		data := doa.Try(os.ReadFile("./conf.json")).([]byte)
		doa.Nil(json.Unmarshal(data, &conf))
		return conf
	}()
)

type Item struct {
	Time     time.Time
	Duration int64
	CommitID string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<head>"))
		w.Write([]byte("<style type=\"text/css\">"))
		w.Write([]byte("* { color: #666666 }"))
		w.Write([]byte("th,td { padding:5px; text-align:center; }"))
		w.Write([]byte("</style>"))
		w.Write([]byte("</head>"))
		w.Write([]byte("<body>"))
		for i := 0; i < len(cConf.Benchs); i++ {
			path := filepath.Base(cConf.Benchs[i])
			data := []Item{}
			doa.Nil(cDb.GetDecode(path, &data))
			w.Write([]byte(fmt.Sprintf("<h3>%s<small><a href=\"%s\">🔗</a></small></h3>", path, cConf.URLs[i])))
			w.Write([]byte("<table>"))
			w.Write([]byte("<thead><tr><th>Date</th><th>Duration(ms)</th><th>CommitID</th></tr></thead>"))
			w.Write([]byte("<tbody>"))
			for _, item := range data {
				w.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", item.Time.Format("2006-01-02 15:04:05"), item.Duration, item.CommitID[:8])))
			}
			w.Write([]byte("</tbody>"))
			w.Write([]byte("</table>"))
		}
		w.Write([]byte("<h2>CPU Info</h2>"))
		w.Write([]byte(fmt.Sprintf("<pre><code>%s</code></pre>", string(doa.Try(os.ReadFile("/proc/cpuinfo")).([]byte)))))
		w.Write([]byte("</code></pre>"))
		w.Write([]byte("</body>"))
	})
	log.Println("listen and server on :8080")
	doa.Nil(http.ListenAndServe(":8080", nil))
}
