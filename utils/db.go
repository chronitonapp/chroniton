package utils

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	DB_PARAMS_TEMPLATE = "dbname=%v host=%v port=%v user=%v password=%v sslmode=disable"
)

var (
	ORM gorm.DB
)

type urlInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	Path     string
}

func ParseUrl(u *url.URL) urlInfo {
	uinfo := urlInfo{
		Host:     u.Host,
		Username: u.User.Username(),
		Path:     u.Path,
	}
	uinfo.Password, _ = u.User.Password()
	uinfo.Path = uinfo.Path[1:]
	sep := strings.Index(uinfo.Host, ":")
	port, _ := strconv.Atoi(uinfo.Host[sep+1:])
	uinfo.Port = port
	uinfo.Host = uinfo.Host[:sep]

	return uinfo
}

func init() {
	var err error
	fmt.Printf("Connecting to database...")
	rawDbURL := os.Getenv("CHRONITON_DB_URL")
	dbURL, _ := url.Parse(rawDbURL)
	parsedDbURL := ParseUrl(dbURL)
	ORM, err = gorm.Open("postgres", fmt.Sprintf(
		DB_PARAMS_TEMPLATE, parsedDbURL.Path,
		parsedDbURL.Host, parsedDbURL.Port,
		parsedDbURL.Username, parsedDbURL.Password,
	))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
	ORM.LogMode(true)
}
