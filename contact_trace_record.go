package main

import "time"

//ContactTraceRecord Contact Trace Record struct
type ContactTraceRecord struct {
	AWSAccountID               string                 `json:"AWSAccountId"`
	Agent                      *Agent                 `json:"Agent"`
	AgentConnectionAttempts    int                    `json:"AgentConnectionAttempts"`
	Attributes                 map[string]string      `json:"Attributes"`
	Channel                    string                 `json:"Channel"`
	ConnectedToSystemTimestamp time.Time              `json:"ConnectedToSystemTimestamp"`
	ContactID                  string                 `json:"ContactId"`
	CustomerEndpoint           CustomerEndpoint       `json:"CustomerEndpoint"`
	DisconnectTimestamp        time.Time              `json:"DisconnectTimestamp"`
	InitialContactID           *string                `json:"InitialContactId"`
	InitiationMethod           string                 `json:"InitiationMethod"`
	InitiationTimestamp        time.Time              `json:"InitiationTimestamp"`
	InstanceARN                string                 `json:"InstanceARN"`
	LastUpdateTimestamp        time.Time              `json:"LastUpdateTimestamp"`
	NextContactID              *string                `json:"NextContactId"`
	PreviousContactID          *string                `json:"PreviousContactId"`
	Queue                      Queue                  `json:"Queue"`
	Recording                  *Recording             `json:"Recording"`
	SystemEndpoint             SystemEndpoint         `json:"SystemEndpoint"`
	TransferCompletedTimestamp *time.Time             `json:"TransferCompletedTimestamp"`
	TransferredToEndpoint      *TransferredToEndpoint `json:"TransferredToEndpoint"`
}

//RoutingProfile Struct for routing profile
type RoutingProfile struct {
	ARN  string `json:"ARN"`
	Name string `json:"Name"`
}

//Agent struct for agent
type Agent struct {
	ARN                            string         `json:"ARN"`
	AfterContactWorkDuration       int            `json:"AfterContactWorkDuration"`
	AfterContactWorkEndTimestamp   time.Time      `json:"AfterContactWorkEndTimestamp"`
	AfterContactWorkStartTimestamp time.Time      `json:"AfterContactWorkStartTimestamp"`
	AgentInteractionDuration       int            `json:"AgentInteractionDuration"`
	ConnectedToAgentTimestamp      time.Time      `json:"ConnectedToAgentTimestamp"`
	CustomerHoldDuration           int            `json:"CustomerHoldDuration"`
	HierarchyGroups                interface{}    `json:"HierarchyGroups"`
	LongestHoldDuration            int            `json:"LongestHoldDuration"`
	NumberOfHolds                  int            `json:"NumberOfHolds"`
	RoutingProfile                 RoutingProfile `json:"RoutingProfile"`
	Username                       string         `json:"Username"`
}

//CustomerEndpoint Struct for customerendpoint
type CustomerEndpoint struct {
	Address string `json:"Address"`
	Type    string `json:"Type"`
}

//Queue Struct for queue
type Queue struct {
	ARN              string    `json:"ARN"`
	DequeueTimestamp time.Time `json:"DequeueTimestamp"`
	Duration         int       `json:"Duration"`
	EnqueueTimestamp time.Time `json:"EnqueueTimestamp"`
	Name             string    `json:"Name"`
}

//Recording struct for recording
type Recording struct {
	Location string `json:"Location"`
	Status   string `json:"Status"`
}

//SystemEndpoint Struct for systemendpoint
type SystemEndpoint struct {
	Address string `json:"Address"`
	Type    string `json:"Type"`
}

//TransferredToEndpoint struct for transferred to endpoint
type TransferredToEndpoint struct {
	Address string `json:"Address"`
	Type    string `json:"Type"`
}
