package main

import "log"
import "github.com/tobz/go-vertica"

func main() {
	log.Printf("trying to connect to vertica...")

	dsn := "dashboard:Peegh3duQuookoh6Haiwoo7a@tcp(vertica-prod-17h-01.node.us-east-1.prod.localytics.io:5433)/localytics_production?tls=skip-verify"
	conn, err := govertica.NewConnection(dsn)
	if err != nil {
		log.Fatalf("caught error getting new connection: %s", err)
	}

	err = conn.Connect()
	if err != nil {
		log.Fatalf("caught error while connecting: %s", err)
	}

	log.Printf("connected successfully!")
}
