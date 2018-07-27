package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // sets env variables from .env
	"log"
	"net/http"
	"net/smtp"
	// "os"
)

// Globally accessible env struct
var myEnv = Env{}

type Env struct {
	myContactEmail   string `env:"MY_CONTACT_EMAIL"`
	port             string `env:"PORT"`
	smtpPort         string `env:"SMTP_PORT"`
	smtpHostname     string `env:"SMTP_HOSTNAME"`
	defaultSmtpLogin string `env:"DEFAULT_SMTP_LOGIN"`
	defaultPassword  string `env:"DEFAULT_PASSWORD"`
}

// Struct variables must be capitalized so json.Decoder can access and write to them
type Email struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func mySendMail(email Email) error {

	// Set up authentication info
	auth := smtp.PlainAuth(
		"",
		myEnv.defaultSmtpLogin,
		myEnv.defaultPassword,
		myEnv.smtpHostname,
	)

	// Headers delimited by newlines, separated from body by empty newline
	message := []byte(
		"To: " + myEnv.myContactEmail + "\r\n" +
			"Subject: Message via Portfolio: " + email.Name + "\r\n\n" +
			email.Message + "\r\n",
	)

	// Execute send email
	err := smtp.SendMail(
		myEnv.smtpHostname+myEnv.smtpPort, // SMTP address
		auth,
		email.Email,                    // send from
		[]string{myEnv.myContactEmail}, // send to
		message, // message body
	)
	return err
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {

	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	if r.Body == nil {
		log.Println("Request has no body")
		http.Error(w, "Ain't nothin in the body", 400)
		return
	}
	var email Email
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		log.Printf("Error decoding body: %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	// Send email
	err = mySendMail(email)
	if err != nil {
		log.Printf("Error sending email: %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(email)
}

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// myEnv = Env{
	// 	os.Getenv("MY_CONTACT_EMAIL"),
	// 	":" + os.Getenv("PORT"),
	// 	":" + os.Getenv("SMTP_PORT"),
	// 	os.Getenv("SMTP_HOSTNAME"),
	// 	os.Getenv("DEFAULT_SMTP_LOGIN"),
	// 	os.Getenv("DEFAULT_PASSWORD"),
	// }
	myEnv = Env{
		"contact@jchow.club",
		":8000",
		":587",
		"smtp.mailgun.org",
		"postmaster@jchow.club",
		"2a612a906abc94054d3530dead0c6509-8889127d-806e569a",
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/send_mail", sendMailHandler).Methods("POST")
	log.Printf("kraxx mail service listening on port %s", myEnv.port)
	log.Fatal(http.ListenAndServe(myEnv.port, router))
}
