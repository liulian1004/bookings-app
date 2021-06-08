package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

//listen to the channel
func listenForMail() {

	//m := <- app.MailChan //get a message from channel
	// fire a annoymous function
	// listen to the coming date all the time
	go func() {
		for {
			msg := <- app.MailChan
			sendMsq(msg)
		}
	} ()
}

//send email
func sendMsq(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 *time.Second

	client , err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject) // set up subject
	// check it has email template or not
	if m.Template == ""{
			email.SetBody(mail.TextHTML,string(m.Content)) //set up html content
	}else{
		// read the template and replace body with m.Content
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-template/%s",m.Template))
		if err != nil {
			app.ErrorLog.Panicln(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content,1)
		email.SetBody(mail.TextHTML, msgToSend)

	}

	
	err = email.Send(client)

	if err != nil {
		log.Println(err)
	}else{
		log.Println("Email send!")
	}


}