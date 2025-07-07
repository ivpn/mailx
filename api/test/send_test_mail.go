package main

import (
	"flag"
	"log"
	"net/smtp"
)

func main() {
	send()
}

func send() {
	var alias string
	flag.StringVar(&alias, "alias", "", "Alias to use for sending a test email")
	flag.Parse()

	err := smtp.SendMail(
		"0.0.0.0:587",
		nil,
		"foo@bar.com",
		[]string{alias},
		[]byte(
			"From: Foo Bar <foo@bar.com>\r\n"+
				"Content-Type: text/plain; charset=us-ascii\r\n"+
				"Content-Transfer-Encoding: 7bit\r\n"+
				"Subject: Test mail\r\n"+
				"Date: Fri, 8 Mar 2024 11:14:26 +0800\r\n"+
				"To: Baz Quux <baz@quux.com>\r\n\r\n"+
				"This is the email body."),
	)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Email sent successfully")
}
