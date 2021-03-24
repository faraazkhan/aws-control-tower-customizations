package main

import "time"

//EventBridgeEvent holds the Event Structure as forwarded by EventBridge
type EventBridgeEvent struct {
	Version    string
	ID         string
	DetailType string
	Source     string
	Account    string
	Time       time.Time
	Region     string
	Detail     CloudWatchEventDetails
}

//ManagedAccountStatus holds the status object for Managed Account Operations
type ManagedAccountStatus struct {
	OrganizationalUnit struct {
		OrganizationalUnitName string `json:"organizationalUnitName"`
		OrganizationalUnitID   string `json:"organizationalUnitId"`
	} `json:"organizationalUnit"`
	Account struct {
		AccountName string `json:"accountName"`
		AccountID   string `json:"accountId"`
	} `json:"account"`
	State              string `json:"state"`
	Message            string `json:"message"`
	Requestedtimestamp string `json:"requestedTimestamp"`
	Completedtimestamp string `json:"completedTimestamp"`
}

//CloudWatchEventDetails is based on the CloudTrail event structure for Successful Account Creates
type CloudWatchEventDetails struct {
	Eventversion string `json:"eventVersion"`
	UserIdentity struct {
		AccountID string `json:"accountId"`
		InvokedBy string `json:"invokedBy"`
	} `json:"userIdentity"`
	EventTime           time.Time   `json:"eventTime"`
	EventSource         string      `json:"eventSource"`
	EventName           string      `json:"eventName"`
	AwsRegion           string      `json:"awsRegion"`
	SourceIPAddress     string      `json:"sourceIPAddress"`
	UserAgent           string      `json:"userAgent"`
	RequestParameters   interface{} `json:"requestParameters"`
	ResponseElements    interface{} `json:"responseElements"`
	EventID             string      `json:"eventID"`
	ReadOnly            bool        `json:"readOnly"`
	EventType           string      `json:"eventType"`
	ManagementEvent     bool        `json:"managementEvent"`
	EventCategory       string      `json:"eventCategory"`
	RecipientAccountID  string      `json:"recipientAccountId"`
	ServiceEventDetails struct {
		UpdateManagedAccountStatus ManagedAccountStatus `json:"updateManagedAccountStatus"`
		CreateManagedAccountStatus ManagedAccountStatus `json:"createManagedAccountStatus"`
	} `json:"serviceEventDetails"`
}
