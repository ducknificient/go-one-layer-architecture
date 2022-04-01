package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gedesukra/goutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/subosito/gotenv"
)

func init() {
	errEnv := gotenv.Load()
	if errEnv != nil {
		fmt.Println("Error load .env file : ", errEnv)
	}
}

func GetPgDataSource(host string, port string, user string, pwd string, dbname string, sslmode string, poolmaxcon string, poolmaxconnidletime string) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s pool_max_conns=%s pool_max_conn_idle_time=%s",
		host, port, user, pwd, dbname, sslmode, poolmaxcon, poolmaxconnidletime)
}

func main() {

	_appip = os.Getenv("appip")
	_appport = os.Getenv("appport")
	_production = os.Getenv("production")
	_dbhost = os.Getenv("dbhost")
	_dbport = os.Getenv("dbport")
	_dbname = os.Getenv("dbname")
	_dbusr = os.Getenv("dbusr")
	_dbpwd = os.Getenv("dbpwd")
	_dbsslmode = os.Getenv("dbsslmode")
	_dbpoolmaxcon = os.Getenv("dbpoolmaxcon")
	_dbpoolmaxconnidletime = os.Getenv("dbpoolmaxconnidletime")

	/*Connect to Postgresql*/
	dbpool, err := pgxpool.Connect(context.Background(), GetPgDataSource(_dbhost, _dbport, _dbusr, _dbpwd, _dbname, _dbsslmode, _dbpoolmaxcon, _dbpoolmaxconnidletime))
	if err != nil {
		panic(err.Error())
	}
	defer dbpool.Close()

	ipPort := ":" + _appport
	if _production == "true" {
		ipPort = _appip + ":" + _appport
	}

	if _production == "false" {
		fmt.Println("App listening on " + ipPort)
	}

	err = http.ListenAndServe(ipPort, getRouter(dbpool))
	if err != nil {
		panic("Unable start server, ListenAndServe:" + goutils.GetStacktraceError(err))
	}

}
