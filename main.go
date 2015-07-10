package main

import (
	"encoding/json"
	"flag"
	"github.com/myafeier/ws/lib"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	// "time"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	// redisaddr = flag.String("redis", "127.0.0.1:6379", "redis service address")
	assets    = flag.String("assets", "templ", "path to assets")
	homeTempl *template.Template
	testTempl *template.Template
)

var conf *Config

type Config struct {
	Mysql struct {
		Connstr string
	}
}

// func defaultAssetPath() string{
// 	p,err:=build.Default.Import(path, srcDir, mode)
// }

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}
func testHandler(c http.ResponseWriter, req *http.Request) {
	testTempl.Execute(c, req.Host)
}

//初始化
func init() {
	r, err := os.Open("conf/config.json")
	if err != nil {
		log.Fatalln("载入配置文件失败", err)
	}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	testTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "test.html")))
	// rc := lib.Newredisc(*redisaddr, 0)
	// err := rc.StartAndGc()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	db := new(lib.Tips)
	err := db.NewTips(conf.Mysql.Connstr)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/test", testHandler)
	h := lib.NewHub()
	go h.Run()
	go h.Productmessage(db)
	go h.ProductHotMessage(db)

	http.Handle("/ws", lib.WsHandler{H: h})
	log.Println("Server is opening")
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln("Listen & Serve Error！")
	}

}
