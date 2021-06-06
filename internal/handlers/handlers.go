package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/tsawler/bookings-app/internal/config"
	"github.com/tsawler/bookings-app/internal/driver"
	"github.com/tsawler/bookings-app/internal/forms"
	"github.com/tsawler/bookings-app/internal/helpers"
	"github.com/tsawler/bookings-app/internal/models"
	"github.com/tsawler/bookings-app/internal/render"
	"github.com/tsawler/bookings-app/internal/repository"
	"github.com/tsawler/bookings-app/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{

	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// get the reservation stored in the session, here is start date, end date
	res, ok :=  m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}
	//get the room name
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	//store the roomName into res info
	res.Room.RoomName = room.RoomName
	//transfer to time.time format and store in the model structure
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	// parse the date to frontend
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	log.Println("trigger post ")
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout,sd)
	if err != nil {
		helpers.ServerError(w,err)
	}

	endDate, err := time.Parse(layout,ed)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w,err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	log.Println("walk here")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//insert into db
	newReservationID, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:    startDate,
		EndDate:      endDate,
		RoomID:       roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		
		helpers.ServerError(w,err)
		return
	}


	m.App.Session.Put(r.Context(), "reservation", reservation)


	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}


// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout,start)
	if err != nil {
		helpers.ServerError(w,err)
	}

	endDate, err := time.Parse(layout,end)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate) 

	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	//no room, show error and redirect to the search page
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(),"error","No availability")
		http.Redirect(w,r,"/search-availability", http.StatusSeeOther)
		return
	}
	//log.Print("rooms", len(rooms))
	//parse the date to front-end
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}
	//store the start/end date in the session, will be used for the next step: make reservation
	m.App.Session.Put(r.Context(),"reservation", res)

	//render the available room
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

// render the info in the page when there is not template
//	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		m.App.Session.Put(r.Context(), "error","Can't get reservation from session")
		http.Redirect(w,r,"/", http.StatusTemporaryRedirect) //redirect to home page if no make reservation 
		return
	}
	//once reservation complete, it's go the summary reservation page, and reservation data from c.context()
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})

	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// select id from url
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// get the reservation stored in the session
	m.App.Session.Get(r.Context(), "reservation")

	res, ok :=  m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	res.RoomID = roomID
	//update the reservation in the session 
	m.App.Session.Put(r.Context(), "reservation",res)
	//redirect to the make reservation 
	http.Redirect(w,r, "/make-reservation", http.StatusSeeOther)
}