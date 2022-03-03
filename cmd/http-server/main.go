package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/godump/acdb"
	"github.com/godump/doa"
)

var (
	cDb = acdb.Doc("./db")
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
		w.Write([]byte("th,td { padding:5px; text-align:center; }"))
		w.Write([]byte("</style>"))
		w.Write([]byte("</head>"))
		w.Write([]byte("<body>"))
		for _, path := range doa.Try(os.ReadDir("./db")).([]fs.DirEntry) {
			data := []Item{}
			doa.Nil(cDb.GetDecode(path.Name(), &data))
			w.Write([]byte(fmt.Sprintf("<h2 style=\"color: #F4606C;\">%s</h2>", path.Name())))
			w.Write([]byte("<table style=\"color: #F4606C;\">"))
			w.Write([]byte("<thead><tr><th>Date</th><th>Duration(ms)</th><th>CommitID</th></tr></thead>"))
			w.Write([]byte("<tbody>"))
			for _, item := range data {
				w.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", item.Time.Format("2006-01-02 15:04:05"), item.Duration, item.CommitID[:8])))
			}
			w.Write([]byte("</tbody>"))
			w.Write([]byte("</table>"))
		}
		w.Write([]byte("</body>"))
	})
	log.Println("listen and server on :8080")
	doa.Nil(http.ListenAndServe(":8080", nil))
}
