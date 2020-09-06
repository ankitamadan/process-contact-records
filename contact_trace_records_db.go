package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

//InsertKinesisContactTraceRecords Insert into contact_trace_records
func InsertKinesisContactTraceRecords(kinesisRecord []byte, clog *ContactTraceRecord, db *sql.DB) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Unable to begin transaction for inserting into contact_trace_records")
		return err
	}

	// add the JSON data to the table, the initial trace record
	isql := "INSERT INTO contact_trace_records (contact_id, record) VALUES ($1, $2) ON CONFLICT (contact_id) DO UPDATE SET record = $3"

	_, err = tx.Exec(isql, clog.ContactID, kinesisRecord, kinesisRecord)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to insert contact_trace_record with contactID %s: ", clog.ContactID)
		err := errors.Wrap(err, errorMessage)
		tx.Rollback()
		return err
	}

	err = InsertContact(tx, clog, DB)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = InsertContactQueue(tx, clog.ContactID, &clog.Queue, DB)
	if err != nil {
		tx.Rollback()
		return err
	}

	if clog.Agent != nil {
		err := InsertContactAgent(tx, clog.ContactID, clog.Agent, DB)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = InsertContactAttributes(tx, clog.ContactID, clog.Attributes, DB)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		errorMessage := fmt.Sprintf("Unable to insert contact_trace_record with contactID %s: ", clog.ContactID)
		err := errors.Wrap(err, errorMessage)
		return err
	}

	return nil
}

//InsertContact Insert into contacts
func InsertContact(tx *sql.Tx, clog *ContactTraceRecord, db *sql.DB) error {
	var err error
	//tx, err := db.Begin()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to begin transaction for inserting contacts")
		err := errors.Wrap(err, errorMessage)
		return err
	}

	// the base contact entry
	isql := `
	INSERT INTO contacts
	    ( contact_id, agent_username, agent_connection_attempts, connected_to_system ,
		customer_endpoint_number, disconnected, initial_contact_id, initiation_method,
		initiation, last_update, next_contact_id, previous_contact_id, system_endpoint_number,
		recording_location, transfer_completed, transferred_to_endpoint_number )
	VALUES
	    ( $1, $2, $3, $4 , $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16 )
	ON CONFLICT (contact_id) DO NOTHING 
	`
	var agentUsername *string = nil
	if clog.Agent != nil {
		agentUsername = &clog.Agent.Username
	}
	var transferredEndpointNumber *string = nil
	if clog.TransferredToEndpoint != nil {
		transferredEndpointNumber = &clog.TransferredToEndpoint.Address
	}
	var recordingLocation *string = nil
	if clog.Recording != nil {
		recordingLocation = &clog.Recording.Location
	}

	_, err = tx.Exec(
		isql,
		clog.ContactID,
		agentUsername,
		clog.AgentConnectionAttempts,
		clog.ConnectedToSystemTimestamp,
		clog.CustomerEndpoint.Address,
		clog.DisconnectTimestamp,
		clog.InitialContactID,
		clog.InitiationMethod,
		clog.InitiationTimestamp,
		clog.LastUpdateTimestamp,
		clog.NextContactID,
		clog.PreviousContactID,
		clog.SystemEndpoint.Address,
		recordingLocation,
		clog.TransferCompletedTimestamp,
		*transferredEndpointNumber,
	)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to insert contact with contactID %s: ", clog.ContactID)
		err := errors.Wrap(err, errorMessage)
		return err
	}

	return nil
}

//InsertContactAgent Insert into contact_agent
func InsertContactAgent(tx *sql.Tx, contactID string, agent *Agent, db *sql.DB) error {

	var err error

	isql := `
	INSERT INTO contact_agent
		( contact_id, username, routing_profile_name, number_of_holds, longest_host_duration,
		customer_hold_duration, connected_to_agent, agent_interaction_duration,
		after_contact_work_start, after_contact_work_end, after_contact_work_duration, agent_arn )
	VALUES
	    ( $1, $2, $3, $4 , $5, $6, $7, $8, $9, $10, $11, $12 )
	ON CONFLICT (contact_id) DO NOTHING 
	`
	_, err = tx.Exec(
		isql,
		contactID,
		agent.Username,
		agent.RoutingProfile.Name,
		agent.NumberOfHolds,
		agent.LongestHoldDuration,
		agent.CustomerHoldDuration,
		agent.ConnectedToAgentTimestamp,
		agent.AgentInteractionDuration,
		agent.AfterContactWorkStartTimestamp,
		agent.AfterContactWorkEndTimestamp,
		agent.AfterContactWorkDuration,
		agent.ARN,
	)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to insert contact agent information with contactID %s: ", contactID)
		err := errors.Wrap(err, errorMessage)
		return err
	}

	return nil
}

//InsertContactQueue Insert into contact_queue
func InsertContactQueue(tx *sql.Tx, contactID string, queue *Queue, db *sql.DB) error {

	var err error

	isql := `
	INSERT INTO contact_queue
		( contact_id, dequeue, duration, enqueue, name )
	VALUES
	    ( $1, $2, $3, $4 , $5 )
	ON CONFLICT (contact_id) DO NOTHING 
	`
	_, err = tx.Exec(
		isql,
		contactID,
		queue.DequeueTimestamp,
		queue.Duration,
		queue.EnqueueTimestamp,
		queue.Name,
	)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to insert contact queue information with contactID %s: ", contactID)
		err := errors.Wrap(err, errorMessage)
		return err
	}

	return nil
}

//InsertContactAttributes insert into contact_attributes
func InsertContactAttributes(tx *sql.Tx, contactID string, attribs map[string]string, db *sql.DB) error {

	isql := `
	INSERT INTO contact_attributes
		( contact_id, attribute, value )
	VALUES
	    ( $1, $2, $3 )
	ON CONFLICT (contact_id, attribute) DO NOTHING 
	`

	stmt, err := tx.Prepare(isql)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to prepare statement for contact attributes with contactID %s: ", contactID)
		err := errors.Wrap(err, errorMessage)
		return err
	}

	for k, v := range attribs {
		_, err = stmt.Exec(contactID, k, v)
		if err != nil {
			errorMessage := fmt.Sprintf("Unable to insert contact attributes with contactID %s: ", contactID)
			err := errors.Wrap(err, errorMessage)
			return err
		}
	}

	return nil
}
