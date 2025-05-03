package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Company struct {
	CompanyID   int            `json:"companyID"`
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	Url         sql.NullString `json:"url"`
	HqCity      sql.NullString `json:"hqCity"`
	HqState     sql.NullString `json:"hqState"`
}

type Role struct {
	RoleID              int            `json:"roleID"`
	CompanyID           int            `json:"companyID"`
	Name                sql.NullString `json:"name"`
	Url                 sql.NullString `json:"url"`
	Description         sql.NullString `json:"description"`
	CoverLetter         sql.NullString `json:"coverLetter"`
	ApplicationLocation sql.NullString `json:"applicationLocation"`
	AppliedDate         sql.NullString `json:"appliedDate"`
	ClosedDate          sql.NullString `json:"closedDate"`
	PostedRangeMin      sql.NullInt64  `json:"postedRangeMin"`
	PostedRangeMax      sql.NullInt64  `json:"postedRangeMax"`
	Equity              sql.NullBool   `json:"equity"`
	WorkCity            sql.NullString `json:"workCity"`
	WorkState           sql.NullString `json:"workState"`
	Location            sql.NullString `json:"location"`
	Status              sql.NullString `json:"status"`
	Discovery           sql.NullString `json:"discovery"`
	Referral            sql.NullBool   `json:"referral"`
	Notes               sql.NullString `json:"notes"`
}

type Interview struct {
	InterviewID int            `json:"interviewID"`
	RoleID      int            `json:"roleID"`
	Date        sql.NullString `json:"date"`
	Start       sql.NullString `json:"start"`
	End         sql.NullString `json:"end"`
	Notes       sql.NullString `json:"notes"`
	Type        sql.NullString `json:"type"`
}

type Contact struct {
	ContactID int            `json:"contactID"`
	CompanyID int            `json:"companyID"`
	FirstName sql.NullString `json:"firstName"`
	LastName  sql.NullString `json:"lastName"`
	Role      sql.NullString `json:"role"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	Linkedin  sql.NullString `json:"linkedin"`
	Notes     sql.NullString `json:"notes"`
}

type InterviewsContacts struct {
	InterviewsContactId int `json:"interviewsContactId"`
	InterviewId         int `json:"interviewId"`
	ContactId           int `json:"contactId"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("PRAGMA foreign_keys = ON;")

	http.HandleFunc("/companies", handleCompanies)
	http.HandleFunc("/roles", handleRoles)
	http.HandleFunc("/interviews", handleInterviews)
	http.HandleFunc("/contacts", handleContacts)
	http.HandleFunc("/interviews_contacts", handleInterviewsContacts)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCompanies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT * FROM Companies")
		var list []Company
		for rows.Next() {
			var x Company
			rows.Scan(&x.CompanyID, &x.Name, &x.Description, &x.Url, &x.HqCity, &x.HqState)
			list = append(list, x)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var x Company
		json.NewDecoder(r.Body).Decode(&x)
		res, _ := db.Exec("INSERT INTO Companies (name, description, url, hqCity, hqState) VALUES (?, ?, ?, ?, ?)", x.Name, x.Description, x.Url, x.HqCity, x.HqState)
		id, _ := res.LastInsertId()
		x.CompanyID = int(id)
		json.NewEncoder(w).Encode(x)
	case http.MethodPut:
		var x Company
		json.NewDecoder(r.Body).Decode(&x)
		db.Exec("UPDATE Companies SET name=?, description=?, url=?, hqCity=?, hqState=? WHERE companyID=?", x.Name, x.Description, x.Url, x.HqCity, x.HqState, x.CompanyID)
		json.NewEncoder(w).Encode(x)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM Companies WHERE companyID=?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT * FROM Roles")
		var list []Role
		for rows.Next() {
			var x Role
			rows.Scan(&x.RoleID, &x.CompanyID, &x.Name, &x.Url, &x.Description, &x.CoverLetter, &x.ApplicationLocation, &x.AppliedDate, &x.ClosedDate, &x.PostedRangeMin, &x.PostedRangeMax, &x.Equity, &x.WorkCity, &x.WorkState, &x.Location, &x.Status, &x.Discovery, &x.Referral, &x.Notes)
			list = append(list, x)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var x Role
		json.NewDecoder(r.Body).Decode(&x)
		res, _ := db.Exec(`INSERT INTO Roles (companyID, name, url, description, coverLetter, applicationLocation, appliedDate, closedDate, postedRangeMin, postedRangeMax, equity, workCity, workState, location, status, discovery, referral, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			x.CompanyID, x.Name, x.Url, x.Description, x.CoverLetter, x.ApplicationLocation, x.AppliedDate, x.ClosedDate, x.PostedRangeMin, x.PostedRangeMax, x.Equity, x.WorkCity, x.WorkState, x.Location, x.Status, x.Discovery, x.Referral, x.Notes)
		id, _ := res.LastInsertId()
		x.RoleID = int(id)
		json.NewEncoder(w).Encode(x)
	case http.MethodPut:
		var x Role
		json.NewDecoder(r.Body).Decode(&x)
		db.Exec(`UPDATE Roles SET companyID=?, name=?, url=?, description=?, coverLetter=?, applicationLocation=?, appliedDate=?, closedDate=?, postedRangeMin=?, postedRangeMax=?, equity=?, workCity=?, workState=?, location=?, status=?, discovery=?, referral=?, notes=? WHERE roleID=?`,
			x.CompanyID, x.Name, x.Url, x.Description, x.CoverLetter, x.ApplicationLocation, x.AppliedDate, x.ClosedDate, x.PostedRangeMin, x.PostedRangeMax, x.Equity, x.WorkCity, x.WorkState, x.Location, x.Status, x.Discovery, x.Referral, x.Notes, x.RoleID)
		json.NewEncoder(w).Encode(x)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM Roles WHERE roleID=?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleInterviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT * FROM Interviews")
		var list []Interview
		for rows.Next() {
			var x Interview
			rows.Scan(&x.InterviewID, &x.RoleID, &x.Date, &x.Start, &x.End, &x.Notes, &x.Type)
			list = append(list, x)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var x Interview
		json.NewDecoder(r.Body).Decode(&x)
		res, _ := db.Exec("INSERT INTO Interviews (roleID, date, start, end, notes, type) VALUES (?, ?, ?, ?, ?, ?)", x.RoleID, x.Date, x.Start, x.End, x.Notes, x.Type)
		id, _ := res.LastInsertId()
		x.InterviewID = int(id)
		json.NewEncoder(w).Encode(x)
	case http.MethodPut:
		var x Interview
		json.NewDecoder(r.Body).Decode(&x)
		db.Exec("UPDATE Interviews SET roleID=?, date=?, start=?, end=?, notes=?, type=? WHERE interviewID=?", x.RoleID, x.Date, x.Start, x.End, x.Notes, x.Type, x.InterviewID)
		json.NewEncoder(w).Encode(x)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM Interviews WHERE interviewID=?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleContacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT * FROM Contacts")
		var list []Contact
		for rows.Next() {
			var x Contact
			rows.Scan(&x.ContactID, &x.CompanyID, &x.FirstName, &x.LastName, &x.Role, &x.Email, &x.Phone, &x.Linkedin, &x.Notes)
			list = append(list, x)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var x Contact
		json.NewDecoder(r.Body).Decode(&x)
		res, _ := db.Exec("INSERT INTO Contacts (companyID, firstName, lastName, role, email, phone, linkedin, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", x.CompanyID, x.FirstName, x.LastName, x.Role, x.Email, x.Phone, x.Linkedin, x.Notes)
		id, _ := res.LastInsertId()
		x.ContactID = int(id)
		json.NewEncoder(w).Encode(x)
	case http.MethodPut:
		var x Contact
		json.NewDecoder(r.Body).Decode(&x)
		db.Exec("UPDATE Contacts SET companyID=?, firstName=?, lastName=?, role=?, email=?, phone=?, linkedin=?, notes=? WHERE contactID=?", x.CompanyID, x.FirstName, x.LastName, x.Role, x.Email, x.Phone, x.Linkedin, x.Notes, x.ContactID)
		json.NewEncoder(w).Encode(x)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM Contacts WHERE contactID=?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleInterviewsContacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, _ := db.Query("SELECT * FROM InterviewsContacts")
		var list []InterviewsContacts
		for rows.Next() {
			var x InterviewsContacts
			rows.Scan(&x.InterviewsContactId, &x.InterviewId, &x.ContactId)
			list = append(list, x)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var x InterviewsContacts
		json.NewDecoder(r.Body).Decode(&x)
		res, _ := db.Exec("INSERT INTO InterviewsContacts (interviewId, contactId) VALUES (?, ?)", x.InterviewId, x.ContactId)
		id, _ := res.LastInsertId()
		x.InterviewsContactId = int(id)
		json.NewEncoder(w).Encode(x)
	case http.MethodPut:
		var x InterviewsContacts
		json.NewDecoder(r.Body).Decode(&x)
		db.Exec("UPDATE InterviewsContacts SET interviewId=?, contactId=? WHERE interviewsContactId=?", x.InterviewId, x.ContactId, x.InterviewsContactId)
		json.NewEncoder(w).Encode(x)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		db.Exec("DELETE FROM InterviewsContacts WHERE interviewsContactId=?", id)
		w.WriteHeader(http.StatusNoContent)
	}
}
