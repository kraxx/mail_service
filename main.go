package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // sets env variables from .env
	"log"
	"net/http"
	"net/smtp"
	"os"
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
	noReply          string `env:"CAMAGRU_NOREPLY"`
}

// Struct variables must be capitalized so json.Decoder can access and write to them
type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Function that takes form data and executes the mail delivery
type MailerFunc func(FormData) error

// Set up authentication info
func setupSmtpAuth() smtp.Auth {
	return smtp.PlainAuth(
		"",
		myEnv.defaultSmtpLogin,
		myEnv.defaultPassword,
		myEnv.smtpHostname,
	)
}

func sendKraxxSiteMail(data FormData) error {

	// Set up authentication info
	auth := setupSmtpAuth()

	// Headers delimited by newlines, separated from body by empty newline
	message := []byte(
		"To: " + myEnv.myContactEmail + "\r\n" +
			"Subject: Message via Portfolio: " + data.Name + "\r\n\n" +
			data.Message + "\r\n",
	)

	// Execute send email
	err := smtp.SendMail(
		myEnv.smtpHostname+myEnv.smtpPort, // SMTP address
		auth,
		data.Email,                     // send from
		[]string{myEnv.myContactEmail}, // send to
		message,                        // message body
	)
	return err
}

func sendCamagruMail(data FormData) error {

	// Set up authentication info
	auth := setupSmtpAuth()

	// Headers delimited by newlines, separated from body by empty newline
	message := []byte(
		"To: " + data.Email + "\r\n" +
			"Camagru - User Verification\r\n\n" +
			data.Message + "\r\n",
	)

	// Execute send email
	err := smtp.SendMail(
		myEnv.smtpHostname+myEnv.smtpPort, // SMTP address
		auth,
		myEnv.noReply,        // send from
		[]string{data.Email}, // send to
		message,              // message body
	)
	return err
}

func mailHandler(mailer MailerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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
		var data FormData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("Error decoding body: %s", err.Error())
			http.Error(w, err.Error(), 400)
			return
		}

		// Send email
		err = mailer(data)
		if err != nil {
			log.Printf("Error sending email: %s", err.Error())
			http.Error(w, err.Error(), 400)
			return
		}

		w.WriteHeader(200)
	}
}

// Load .env
func init() {

	/*
		Local

		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	*/

	// Production
	myEnv = Env{
		os.Getenv("MY_CONTACT_EMAIL"),
		":" + os.Getenv("PORT"),
		":" + os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_HOSTNAME"),
		os.Getenv("DEFAULT_SMTP_LOGIN"),
		os.Getenv("DEFAULT_PASSWORD"),
		os.Getenv("CAMAGRU_NOREPLY"),
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/send_mail", mailHandler(sendKraxxSiteMail)).Methods("POST")
	router.HandleFunc("/camagru_mail", mailHandler(sendCamagruMail)).Methods("POST")
	log.Printf("kraxx mail service listening on port %s", myEnv.port)
	log.Fatal(http.ListenAndServe(myEnv.port, router))
}
