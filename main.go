package main

import (
	"flag"
	"fmt"
	"github.com/gopher1980/dynql"
	"github.com/sclevine/agouti"
	"time"

	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	pages        map[string]*agouti.Page
	agoutiDriver *agouti.WebDriver
	debug        bool
)

func file2str(ilename string) string {
	file, err := os.Open(ilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return string(b)
}

type Param struct {
	Session  string `json:"session"`
	Value    string `json:"value"`
	Option   string `json:"option"`
	Browser  string `json:"browser"`
	Url      string `json:"url"`
	Selector string `json:"selector"`
	Strategy string `json:"strategy"`
	Duration int64  `json:"duration"`
}

type ReturnStep struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func start(name string, ptr interface{}, r *http.Request) interface{} {

	var err error
	p := ptr.(*Param)
	skey := p.Session
	if skey == "" {
		skey = r.URL.Query().Get("s")
	}

	if skey == "" {
		return ReturnStep{Message: "Session can not empty", Action: name}
	}

	if pages[skey] != nil {
		return ReturnStep{Message: "page is open", Action: name}
	}

	pages[skey], err = agoutiDriver.NewPage()

	if p.Duration == 0 {
		p.Duration = 40
	}

	go func() {
		time.Sleep(time.Duration(p.Duration) * time.Second)
		fmt.Println("page", skey, " Destroy ", pages[skey].String())
		pages[skey].Destroy()
		pages[skey] = nil
	}()

	if err != nil {
		return err
	}
	return ReturnStep{Message: "ok", Action: name}
}

func destroy(name string, ptr interface{}, r *http.Request) interface{} {

	p := ptr.(*Param)
	skey := p.Session
	if skey == "" {
		skey = r.URL.Query().Get("s")
	}
	page := pages[skey]

	return page.Destroy()

}

func Find(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.Find(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.Find(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByXPath(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByXPath(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByXPath(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByClass(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByClass(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByClass(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByID(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByID(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByID(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByLink(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByLink(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByLink(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByName(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByName(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByName(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByButton(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByButton(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByButton(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByLabel(page *agouti.Page, selector string) *agouti.Selection {
	i := 0
	elem := page.FindByLabel(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < 1000 && page != nil {
		time.Sleep(3 * time.Millisecond)
		elem = page.FindByLabel(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func FindByStrategy(page *agouti.Page, strategy string, selector string) *agouti.Selection {
	if strategy == "xpath" {
		return FindByXPath(page, selector)
	}
	if strategy == "name" {
		return FindByName(page, selector)
	}
	if strategy == "class" {
		return FindByClass(page, selector)
	}
	if strategy == "id" {
		return FindByID(page, selector)
	}
	if strategy == "link" {
		return FindByLink(page, selector)
	}
	if strategy == "button" {
		return FindByButton(page, selector)
	}
	if strategy == "label" {
		return FindByButton(page, selector)
	}
	if strategy == "label" {
		return FindByLabel(page, selector)
	}

	return Find(page, selector)
}

//  brew install chromedriver
func jsRun(name string, ptr interface{}, r *http.Request) interface{} {
	var err error
	p := ptr.(*Param)
	skey := p.Session
	if skey == "" {
		skey = r.URL.Query().Get("s")
	}
	page := pages[skey]

	if page == nil {
		return ReturnStep{Message: "page nil", Action: name}
	}
	vm := otto.New()

	vm.Set("session", skey)
	vm.Set("req", r)
	vm.Set("param", ptr.(*Param))
	vm.Set("page", page)
	vm.Set("action_name", name)

	vm.Set("sleep", time.Sleep)
	vm.Set("Second", time.Second)
	vm.Set("Millisecond", time.Millisecond)
	vm.Set("Microsecond", time.Microsecond)
	vm.Set("Nanosecond", time.Nanosecond)

	vm.Set("Find", FindByStrategy)

	code := file2str("actions/" + name + ".js")

	result, err := vm.Run(code + "\n  action()")
	if err != nil {
		return err
	}

	export, err := result.Export()
	if err != nil {
		return err
	}
	return export
}

//  brew install chromedriver
func main() {
	var driver string
	var addr string

	flag.StringVar(&driver, "driver", "phantomjs", "[drive, phantomjs, selenium, firefox ]")
	flag.StringVar(&addr, "addr", ":9090", "this is addr of this ListenAndServe")
	flag.BoolVar(&debug, "debug", false, "sql debug")

	flag.Parse()

	if driver == "chrome" {
		agoutiDriver = agouti.ChromeDriver()
	}
	if driver == "phantomjs" {
		agoutiDriver = agouti.PhantomJS()
	}

	if driver == "selenium" {
		agoutiDriver = agouti.Selenium()
	}

	if driver == "firefox" {
		agoutiDriver = agouti.Selenium(agouti.Browser("firefox"))
	}

	err := agoutiDriver.Start()
	if err != nil {
		panic(err)
	}
	pages = make(map[string]*agouti.Page)
	fmt.Println("OwlQA v0.0.1")
	fmt.Println(driver)
	dql := dynql.NewDQL()
	dql.Put("start", start, Param{})
	dql.Put("default", jsRun, Param{})
	dql.Put("destroy", destroy, Param{})

	r := mux.NewRouter()
	r.HandleFunc("/owl", dql.Run).Methods(http.MethodPost)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(addr, nil))
}
