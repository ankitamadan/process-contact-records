package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

var (
	dbString string
	//DB Initialized database
	DB *sql.DB
)

func init() {
	initDatabase()
}

func initDatabase() {
	var err error
	dbString = GetDatabaseString()

	DB, err = sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("unable to connect to ***** with error: %s", err.Error())
	}
}

func main() {
	log.Println("process-contact-flow-logs click...")
	lambda.Start(handler)
	log.Println("process-contact-flow-logs clunk...")
}

func insertContactTraceRecord(ctx context.Context, r events.KinesisEventRecord) {
	var err error
	// unmarshal the contact trace record
	var clog ContactTraceRecord
	err = json.Unmarshal(r.Kinesis.Data, &clog)
	if err != nil {
		log.Fatalf("Unable to unmarshal Contact Trace Record: %s %s", r.EventID, err.Error())
	}

	// the various tables
	err = InsertKinesisContactTraceRecords(r.Kinesis.Data, &clog, DB)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func handler(ctx context.Context, event events.KinesisEvent) error {
	var err error
	for _, r := range event.Records {

		insertContactTraceRecord(ctx, r)

	}
	return err
}
