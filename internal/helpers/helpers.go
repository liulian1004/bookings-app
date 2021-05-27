package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/tsawler/bookings-app/internal/config"
)

var app *config.AppConfig

//set up app config
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with the status of", status)
	http.Error(w,http.StatusText(status), status)
}

//get the error message + error trace 
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s",err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}