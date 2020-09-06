package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	mockDatabase                  *sql.DB
	record                        []byte
	contactID                     = "a2808b12-0192-45f8-9242-a3a6ed5e2d91"
	agentUserName                 = "newAgent"
	agentConnectionAttempts       = 1
	connectedToSystem             = time.Now()
	customerEndpointNumberAddress = "+61730634595"
	disconnectedTimestamp         = time.Now()
	initialContactID              = "nil"
	initiationMethod              = "INBOUND"
	initiationTimestamp           = time.Now()
	lastUpdateTimestamp           = time.Now()
	nextContactID                 = "nil"
	previousContactID             = "nil"
	systemEndpointNumber          = "+61730634595"
	recordingLocation             = "newlocation"
	transferCompletedTimestamp    time.Time
	transferToEndpointNumber      = "+61730634595"
	arn                           = "arn:aws:connect:ap-southeast-2:983623687068:instance/fba6354d-98a4-42d9-9675-8987c604ed92/routing-profile/976b62e6-4dba-4dee-b44b-1b59da4afaad"
	routingName                   = "testName"
	numberOfHolds                 = 1
	longestHoldDuration           = 10
	customerHoldDuration          = 10
	connectedToAgentTimestamp     = time.Now()
	agentTransactionDuration      = 12
	agentContactWorkTimestamp     = time.Now()
	afterContactWorkTimestamp     = time.Now()
	afterContactWorkDuration      = 132
	agentInterationDuration       = 124
	dequeueTimestamp              = time.Now()
	duration                      = 10
	enqueueTimestamp              = time.Now()
	name                          = "newQueue"
)

func TestInsertJSONbSuccess(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mockDatabase := db
	defer db.Close()

	record = jsonRecord("test_data/interaction1.json")
	mock.ExpectBegin()
	isql1 := "INSERT INTO contact_trace_records (contact_id, record) VALUES ($1, $2) ON CONFLICT (contact_id) DO UPDATE SET record = $3"

	mock.ExpectExec(regexp.QuoteMeta(isql1)).
		WithArgs(contactID, record, record).
		WillReturnResult(sqlmock.NewResult(1, 1))

	transferredToEndpoint := TransferredToEndpoint{Address: transferToEndpointNumber}
	systemEndpoint := SystemEndpoint{Address: systemEndpointNumber}
	recording := Recording{Location: recordingLocation}
	customerEndpoint := CustomerEndpoint{Address: customerEndpointNumberAddress}
	routingProfile := RoutingProfile{Name: routingName}
	queue := Queue{DequeueTimestamp: dequeueTimestamp, Duration: duration, EnqueueTimestamp: enqueueTimestamp, Name: name}
	agentValue := Agent{Username: agentUserName, RoutingProfile: routingProfile, NumberOfHolds: numberOfHolds, LongestHoldDuration: longestHoldDuration, CustomerHoldDuration: customerHoldDuration,
		ConnectedToAgentTimestamp: connectedToAgentTimestamp, AgentInteractionDuration: agentInterationDuration, AfterContactWorkDuration: afterContactWorkDuration, AfterContactWorkEndTimestamp: afterContactWorkTimestamp, ARN: arn}

	var attributes = make(map[string]string)
	attributes["brandServiceNumber"] = "+61430634595"

	contactTraceRecord := ContactTraceRecord{ContactID: contactID, AgentConnectionAttempts: agentConnectionAttempts,
		ConnectedToSystemTimestamp: connectedToSystem, CustomerEndpoint: customerEndpoint, DisconnectTimestamp: disconnectedTimestamp,
		InitialContactID: &initialContactID, InitiationMethod: initiationMethod, InitiationTimestamp: initiationTimestamp,
		LastUpdateTimestamp: lastUpdateTimestamp, NextContactID: &nextContactID, PreviousContactID: &previousContactID, SystemEndpoint: systemEndpoint,
		Recording: &recording, TransferCompletedTimestamp: &transferCompletedTimestamp, TransferredToEndpoint: &transferredToEndpoint, Agent: &agentValue, Queue: queue, Attributes: attributes}

	var agentUsername *string = nil
	if contactTraceRecord.Agent != nil {
		agentUsername = &contactTraceRecord.Agent.Username
	}
	var transferredEndpointNumber *string = nil
	if contactTraceRecord.TransferredToEndpoint != nil {
		transferredEndpointNumber = &contactTraceRecord.TransferredToEndpoint.Address
	}
	var recordingLocation *string = nil
	if contactTraceRecord.Recording != nil {
		recordingLocation = &contactTraceRecord.Recording.Location
	}

	isql2 := `
	INSERT INTO contacts
	    ( contact_id, agent_username, agent_connection_attempts, connected_to_system ,
		customer_endpoint_number, disconnected, initial_contact_id, initiation_method,
		initiation, last_update, next_contact_id, previous_contact_id, system_endpoint_number,
		recording_location, transfer_completed, transferred_to_endpoint_number )
	VALUES
	    ( $1, $2, $3, $4 , $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16 )
	ON CONFLICT (contact_id) DO NOTHING 
	`

	mock.ExpectExec(regexp.QuoteMeta(isql2)).
		WithArgs(contactTraceRecord.ContactID,
			agentUsername,
			contactTraceRecord.AgentConnectionAttempts,
			contactTraceRecord.ConnectedToSystemTimestamp,
			contactTraceRecord.CustomerEndpoint.Address,
			contactTraceRecord.DisconnectTimestamp,
			contactTraceRecord.InitialContactID,
			contactTraceRecord.InitiationMethod,
			contactTraceRecord.InitiationTimestamp,
			contactTraceRecord.LastUpdateTimestamp,
			contactTraceRecord.NextContactID,
			contactTraceRecord.PreviousContactID,
			contactTraceRecord.SystemEndpoint.Address,
			recordingLocation,
			contactTraceRecord.TransferCompletedTimestamp,
			*transferredEndpointNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))

	isql4 := `
	INSERT INTO contact_queue
		( contact_id, dequeue, duration, enqueue, name )
	VALUES
	    ( $1, $2, $3, $4 , $5 )
	ON CONFLICT (contact_id) DO NOTHING 
	`

	mock.ExpectExec(regexp.QuoteMeta(isql4)).
		WithArgs(contactID,
			queue.DequeueTimestamp,
			queue.Duration,
			queue.EnqueueTimestamp,
			queue.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	isql3 := `
	INSERT INTO contact_agent
		( contact_id, username, routing_profile_name, number_of_holds, longest_host_duration,
		customer_hold_duration, connected_to_agent, agent_interaction_duration,
		after_contact_work_start, after_contact_work_end, after_contact_work_duration, agent_arn )
	VALUES
	    ( $1, $2, $3, $4 , $5, $6, $7, $8, $9, $10, $11, $12 )
	ON CONFLICT (contact_id) DO NOTHING 
	`

	mock.ExpectExec(regexp.QuoteMeta(isql3)).
		WithArgs(contactID,
			agentValue.Username,
			agentValue.RoutingProfile.Name,
			agentValue.NumberOfHolds,
			agentValue.LongestHoldDuration,
			agentValue.CustomerHoldDuration,
			agentValue.ConnectedToAgentTimestamp,
			agentValue.AgentInteractionDuration,
			agentValue.AfterContactWorkStartTimestamp,
			agentValue.AfterContactWorkEndTimestamp,
			agentValue.AfterContactWorkDuration,
			agentValue.ARN).
		WillReturnResult(sqlmock.NewResult(1, 1))

	isql5 := `
	INSERT INTO contact_attributes
		( contact_id, attribute, value )
	VALUES
	    ( $1, $2, $3 )
	ON CONFLICT (contact_id, attribute) DO NOTHING 
	`
	mock.ExpectPrepare(regexp.QuoteMeta(isql5)).ExpectExec().WithArgs(contactID, "brandServiceNumber", "+61430634595").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = InsertKinesisContactTraceRecords(record, &contactTraceRecord, mockDatabase)
	if err != nil {
		t.Errorf("error was not expected while inserting json %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func jsonRecord(filename string) []byte {
	json, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return json
}
