package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/gopher1980/dynql"
	"github.com/gopher1980/gormcrud"
	"github.com/sclevine/agouti"

	"io/ioutil"
	"os"

	"github.com/robertkrimen/otto"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	frequencyDefault = 10
	retryDefault     = 1000
)

var (
	pages        map[string]*agouti.Page
	agoutiDriver *agouti.WebDriver
	debug        bool
	db           *gorm.DB
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Storage is struct for save data for test
type Storage struct {
	ID        uint       `gorm:"primary_key" json:"id" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" gorm:"unique_index:storage_name"`
	Strategy  string     `json:"strategy"`
	Selector  string     `json:"selector"`
	Value     string     `json:"value"`
	//Frequency int64      `json:"frequency"`
	//Retry     int64      `json:"retry"`
}

// Actions is struct for save data for test
type Action struct {
	ID        uint       `gorm:"primary_key" json:"id" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" gorm:"unique_index:action_name"`
	Code      string     `json:"code"`
}

// Catalog is
type Catalog struct {
	ID        uint       `gorm:"primary_key" json:"id" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
}

// Param recive in method
type Param struct {
	Session   string `json:"session"`
	Value     string `json:"value"`
	Option    string `json:"option"`
	Browser   string `json:"browser"`
	URL       string `json:"url"`
	Condition string `json:"condition"`
	Selector  string `json:"selector"`
	Strategy  string `json:"strategy"`
	Duration  int64  `json:"duration"`
	Frequency int64  `json:"frequency"`
	Retry     int64  `json:"retry"`
	CheckPage bool   `json:"check_page"`
}

// ReturnStep is generic return
type ReturnStep struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func start(name string, ptr interface{}, r *http.Request, payload interface{}, parent interface{}) interface{} {

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

	if p.Duration != -1 {
		go func() {
			time.Sleep(time.Duration(p.Duration) * time.Second)
			fmt.Println("page", skey, " Destroy ", pages[skey].String())
			_ = pages[skey].Destroy()
			pages[skey] = nil
		}()
	}
	if err != nil {
		return err
	}
	return ReturnStep{Message: "ok", Action: name}
}

func destroy(name string, ptr interface{}, r *http.Request, payload interface{}, parent interface{}) interface{} {

	p := ptr.(*Param)
	skey := p.Session
	if skey == "" {
		skey = r.URL.Query().Get("s")
	}
	page := pages[skey]

	return page.Destroy()

}

func find(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.Find(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.Find(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByXPath(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByXPath(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByXPath(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByClass(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByClass(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByClass(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByID(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByID(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByID(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByLink(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByLink(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByLink(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByName(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByName(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByName(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByButton(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByButton(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByButton(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByLabel(page *agouti.Page, selector string, frequency int, retry int) *agouti.Selection {
	i := 0
	elem := page.FindByLabel(selector)
	count, err := elem.Count()
	for count == 0 && err != nil && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByLabel(selector)
		count, err = elem.Count()
		i++
	}
	return elem
}

func findByStrategy(page *agouti.Page, strategy string, selector string, frequency int, retry int) *agouti.Selection {

	if frequency == 0 {
		frequency = frequencyDefault
	}

	if retry == 0 {
		retry = retryDefault
	}

	if strategy == "xpath" {
		return findByXPath(page, selector, frequency, retry)
	}
	if strategy == "name" {
		return findByName(page, selector, frequency, retry)
	}
	if strategy == "class" {
		return findByClass(page, selector, frequency, retry)
	}
	if strategy == "id" {
		return findByID(page, selector, frequency, retry)
	}
	if strategy == "link" {
		return findByLink(page, selector, frequency, retry)
	}
	if strategy == "button" {
		return findByButton(page, selector, frequency, retry)
	}

	if strategy == "label" {
		return findByLabel(page, selector, frequency, retry)
	}
	if strategy == "storage" {
		storage, _, exists := findInStorage(selector)
		if !exists {
			return nil
		}
		return findByStrategy(page, storage.Strategy, storage.Selector, frequency, retry)
	}
	return find(page, selector, frequency, retry)
}

func waitHide(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.Find(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.Find(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByLink(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByLink(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByLink(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByLabel(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByLabel(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByLabel(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByID(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByID(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByID(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByXPath(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByXPath(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByXPath(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByClass(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByClass(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByXPath(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByName(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByName(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByName(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByButton(page *agouti.Page, selector string, frequency int, retry int) bool {
	i := 0
	elem := page.FindByButton(selector)
	count, _ := elem.Count()
	for count != 0 && i < retry && page != nil {
		time.Sleep(time.Duration(frequency) * time.Millisecond)
		elem = page.FindByButton(selector)
		count, _ = elem.Count()
		i++
	}
	if i >= retry {
		return false
	}
	return true
}

func waitHideByStrategy(page *agouti.Page, strategy string, selector string, frequency int, retry int) bool {

	if frequency == 0 {
		frequency = frequencyDefault
	}

	if retry == 0 {
		retry = retryDefault
	}

	if strategy == "xpath" {
		return waitHideByXPath(page, selector, frequency, retry)
	}
	if strategy == "name" {
		return waitHideByName(page, selector, frequency, retry)
	}
	if strategy == "class" {
		return waitHideByClass(page, selector, frequency, retry)
	}
	if strategy == "id" {
		return waitHideByID(page, selector, frequency, retry)
	}
	if strategy == "link" {
		return waitHideByLink(page, selector, frequency, retry)
	}
	if strategy == "button" {
		return waitHideByButton(page, selector, frequency, retry)
	}

	if strategy == "label" {
		return waitHideByLabel(page, selector, frequency, retry)
	}

	if strategy == "storage" {
		storage, _, exists := findInStorage(selector)
		if !exists {
			return false
		}
		return waitHideByStrategy(page, storage.Strategy, storage.Selector, frequency, retry)
	}

	return waitHide(page, selector, frequency, retry)
}

func findInStorage(name string) (Storage, int64, bool) {
	var count int64
	storage := Storage{}
	db.Where("name = ?", name).Find(&storage).Count(&count)
	exists := true
	if count == 0 {
		exists = false
	}

	return storage, count, exists

}

func waitShowByStrategy(page *agouti.Page, strategy string, selector string, frequency int, retry int) bool {
	elem := findByStrategy(page, strategy, selector, frequency, retry)
	if elem == nil {
		return false
	}
	count, err := elem.Count()
	if count == 0 || err != nil {
		return false
	}
	return true
}

func newClientResty() *resty.Client {
	return resty.New()
}
func unmarshal(str string) map[string]interface{} {
	ret := make(map[string]interface{})
	json.Unmarshal([]byte(str), &ret)
	return ret
}

//  brew install chromedriver
func jsRun(name string, ptr interface{}, r *http.Request, payload interface{}, parent interface{}) interface{} {
	var err error
	p := ptr.(*Param)
	skey := p.Session
	if skey == "" {
		skey = r.URL.Query().Get("s")
	}
	page := pages[skey]

	if page == nil && p.CheckPage {
		return ReturnStep{Message: "page nil", Action: name}
	}

	vm := otto.New()
	_ = vm.Set("payload", payload)
	if payload != nil {
		var copyPayload interface{}
		bytes, _ := json.Marshal(payload)
		_ = json.Unmarshal(bytes, &copyPayload)
		_ = vm.Set("payload", copyPayload)
	}
	_ = vm.Set("parent", parent)

	p.Condition = strings.ReplaceAll(p.Condition, "{", "")
	p.Condition = strings.ReplaceAll(p.Condition, "}", "")

	if p.Condition != "" {
		v, err := vm.Run(p.Condition)
		if err != nil {
			return err
		}

		if v.IsBoolean() {
			condition, _ := v.ToBoolean()
			if !condition {
				bytes, _ := json.Marshal(payload)
				mapPayload := make(map[string]interface{})
				_ = json.Unmarshal(bytes, &mapPayload)
				mapPayload["action"] = "!" + name
				mapPayload["message"] = "return false condition"
				return mapPayload
			}
		} else {
			return ReturnStep{Message: "Condition return not boolean", Action: name}
		}

	}

	_ = vm.Set("session", skey)
	_ = vm.Set("req", r)
	_ = vm.Set("param", ptr.(*Param))
	_ = vm.Set("page", page)
	_ = vm.Set("action_name", name)

	_ = vm.Set("sleep", time.Sleep)
	_ = vm.Set("Second", time.Second)
	_ = vm.Set("Millisecond", time.Millisecond)
	_ = vm.Set("Microsecond", time.Microsecond)
	_ = vm.Set("Nanosecond", time.Nanosecond)

	_ = vm.Set("WaitShowByStrategy", waitShowByStrategy)
	_ = vm.Set("WaitHideByStrategy", waitHideByStrategy)
	_ = vm.Set("Find", findByStrategy)
	_ = vm.Set("FindInStorage", findInStorage)
	_ = vm.Set("NewClientResty", newClientResty)

	_ = vm.Set("Unmarshal", unmarshal)

	code := "function action { return {message:'action no exists'}; };\n"
	method := "action"

	if strings.Contains(name, ".") {
		v := strings.Split(name, ".")
		name = v[0]
		method = v[1]
	}

	action := Action{}
	ret := db.Where("name = ?", name).Find(&action)
	if ret.RowsAffected == 0 {
		return ReturnStep{Message: "action not found", Action: name}
	}
	code = action.Code

	result, err := vm.Run(code + "\n  " + method + "()")
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

	flag.StringVar(&driver, "driver", "chrome", "[drive, phantomjs, selenium, firefox ]")
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

	flag.Parse()

	db, err = gorm.Open("sqlite3", "./storage.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)
	db.AutoMigrate(&Storage{})
	db.AutoMigrate(&Catalog{})
	db.AutoMigrate(&Action{})

	pages = make(map[string]*agouti.Page)
	fmt.Println("OwlQA v0.0.1")
	fmt.Println(driver)
	dql := dynql.NewDQL()
	dql.Put("start", start, Param{})
	dql.Put("default", jsRun, Param{})
	dql.Put("destroy", destroy, Param{})

	r := mux.NewRouter()

	rh := http.RedirectHandler("/site/queries.html", 307)
	r.Handle("/", rh)

	r.PathPrefix("/site/").Handler(http.StripPrefix("/site/", http.FileServer(http.Dir("site/"))))
	r.PathPrefix("/webassembly/").Handler(http.StripPrefix("/webassembly/", http.FileServer(http.Dir("webassembly/"))))

	gormcrud.MapMux(r, db).NewMap("/action", Action{}, []Action{}).Full()
	gormcrud.MapMux(r, db).NewMap("/catalog", Catalog{}, []Catalog{}).Full()
	gormcrud.MapMux(r, db).NewMap("/storage", Storage{}, []Storage{}).Full()
	r.HandleFunc("/run", dql.Run).Methods(http.MethodPost)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(addr, nil))
}
