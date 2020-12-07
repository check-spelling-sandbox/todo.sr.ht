// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type ACL interface {
	IsACL()
}

type Entity interface {
	IsEntity()
}

type EventDetail interface {
	IsEventDetail()
}

type Subscription interface {
	IsSubscription()
}

type ACLCursor struct {
	Results []ACL   `json:"results"`
	Cursor  *string `json:"cursor"`
}

type Assignment struct {
	EventType EventType `json:"eventType"`
	Assigner  Entity    `json:"assigner"`
	Assignee  Entity    `json:"assignee"`
}

func (Assignment) IsEventDetail() {}

type Comment struct {
	EventType     EventType    `json:"eventType"`
	Text          string       `json:"text"`
	Authenticity  Authenticity `json:"authenticity"`
	SuperceededBy *Comment     `json:"superceededBy"`
}

func (Comment) IsEventDetail() {}

type DefaultACL struct {
	Browse  bool `json:"browse"`
	Submit  bool `json:"submit"`
	Comment bool `json:"comment"`
	Edit    bool `json:"edit"`
	Triage  bool `json:"triage"`
}

func (DefaultACL) IsACL() {}

type DefaultACLs struct {
	Anonymous ACL `json:"anonymous"`
	Submitter ACL `json:"submitter"`
	LoggedIn  ACL `json:"logged_in"`
}

type EmailAddress struct {
	ID            int       `json:"id"`
	Created       time.Time `json:"created"`
	CanonicalName string    `json:"canonicalName"`
	Mailbox       string    `json:"mailbox"`
	Name          string    `json:"name"`
}

func (EmailAddress) IsEntity() {}

type Event struct {
	ID      int           `json:"id"`
	Created time.Time     `json:"created"`
	Changes []EventDetail `json:"changes"`
	Entity  Entity        `json:"entity"`
	Ticket  *Ticket       `json:"ticket"`
	Tracker *Tracker      `json:"tracker"`
}

type EventCursor struct {
	Results []*Event `json:"results"`
	Cursor  *string  `json:"cursor"`
}

type ExternalUser struct {
	ID            int       `json:"id"`
	Created       time.Time `json:"created"`
	CanonicalName string    `json:"canonicalName"`
	ExternalID    string    `json:"externalId"`
	ExternalURL   *string   `json:"externalUrl"`
}

func (ExternalUser) IsEntity() {}

type Label struct {
	ID              int           `json:"id"`
	Created         time.Time     `json:"created"`
	Name            string        `json:"name"`
	Tracker         *Tracker      `json:"tracker"`
	BackgroundColor string        `json:"backgroundColor"`
	TextColor       string        `json:"textColor"`
	Tickets         *TicketCursor `json:"tickets"`
}

type LabelCursor struct {
	Results []*Label `json:"results"`
	Cursor  *string  `json:"cursor"`
}

type LabelUpdate struct {
	EventType EventType `json:"eventType"`
	Label     *Label    `json:"label"`
}

func (LabelUpdate) IsEventDetail() {}

type StatusChange struct {
	EventType     EventType        `json:"eventType"`
	OldStatus     TicketStatus     `json:"oldStatus"`
	NewStatus     TicketStatus     `json:"newStatus"`
	OldResolution TicketResolution `json:"oldResolution"`
	NewResolution TicketResolution `json:"newResolution"`
}

func (StatusChange) IsEventDetail() {}

type SubscriptionCursor struct {
	Results []Subscription `json:"results"`
	Cursor  *string        `json:"cursor"`
}

type Ticket struct {
	ID           int              `json:"id"`
	Created      time.Time        `json:"created"`
	Updated      time.Time        `json:"updated"`
	Submitter    Entity           `json:"submitter"`
	Tracker      *Tracker         `json:"tracker"`
	Ref          string           `json:"ref"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Status       TicketStatus     `json:"status"`
	Resolution   TicketResolution `json:"resolution"`
	Authenticity Authenticity     `json:"authenticity"`
	Labels       []*Label         `json:"labels"`
	Assignees    []Entity         `json:"assignees"`
}

type TicketCursor struct {
	Results []*Ticket `json:"results"`
	Cursor  *string   `json:"cursor"`
}

type TicketMention struct {
	EventType EventType `json:"eventType"`
	Mentioned *Ticket   `json:"mentioned"`
}

func (TicketMention) IsEventDetail() {}

type TicketSubscription struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	Entity  Entity    `json:"entity"`
	Ticket  *Ticket   `json:"ticket"`
}

func (TicketSubscription) IsSubscription() {}

type Tracker struct {
	ID          int           `json:"id"`
	Created     time.Time     `json:"created"`
	Updated     time.Time     `json:"updated"`
	Owner       Entity        `json:"owner"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Tickets     *TicketCursor `json:"tickets"`
	Labels      *LabelCursor  `json:"labels"`
	Acls        *ACLCursor    `json:"acls"`
	DefaultACLs *DefaultACLs  `json:"defaultACLs"`
}

type TrackerACL struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	Tracker *Tracker  `json:"tracker"`
	Entity  Entity    `json:"entity"`
	Browse  bool      `json:"browse"`
	Submit  bool      `json:"submit"`
	Comment bool      `json:"comment"`
	Edit    bool      `json:"edit"`
	Triage  bool      `json:"triage"`
}

func (TrackerACL) IsACL() {}

type TrackerCursor struct {
	Results []*Tracker `json:"results"`
	Cursor  *string    `json:"cursor"`
}

type TrackerSubscription struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	Entity  Entity    `json:"entity"`
	Tracker *Tracker  `json:"tracker"`
}

func (TrackerSubscription) IsSubscription() {}

type User struct {
	ID            int            `json:"id"`
	Created       time.Time      `json:"created"`
	Updated       time.Time      `json:"updated"`
	CanonicalName string         `json:"canonicalName"`
	Username      string         `json:"username"`
	Email         string         `json:"email"`
	URL           *string        `json:"url"`
	Location      *string        `json:"location"`
	Bio           *string        `json:"bio"`
	Trackers      *TrackerCursor `json:"trackers"`
}

func (User) IsEntity() {}

type UserMention struct {
	EventType EventType `json:"eventType"`
	Mentioned Entity    `json:"mentioned"`
}

func (UserMention) IsEventDetail() {}

type Version struct {
	Major           int        `json:"major"`
	Minor           int        `json:"minor"`
	Patch           int        `json:"patch"`
	DeprecationDate *time.Time `json:"deprecationDate"`
}

type AccessKind string

const (
	AccessKindRo AccessKind = "RO"
	AccessKindRw AccessKind = "RW"
)

var AllAccessKind = []AccessKind{
	AccessKindRo,
	AccessKindRw,
}

func (e AccessKind) IsValid() bool {
	switch e {
	case AccessKindRo, AccessKindRw:
		return true
	}
	return false
}

func (e AccessKind) String() string {
	return string(e)
}

func (e *AccessKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessKind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessKind", str)
	}
	return nil
}

func (e AccessKind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AccessScope string

const (
	AccessScopeProfile       AccessScope = "PROFILE"
	AccessScopeTrackers      AccessScope = "TRACKERS"
	AccessScopeTickets       AccessScope = "TICKETS"
	AccessScopeACLS          AccessScope = "ACLS"
	AccessScopeEvents        AccessScope = "EVENTS"
	AccessScopeSubscriptions AccessScope = "SUBSCRIPTIONS"
)

var AllAccessScope = []AccessScope{
	AccessScopeProfile,
	AccessScopeTrackers,
	AccessScopeTickets,
	AccessScopeACLS,
	AccessScopeEvents,
	AccessScopeSubscriptions,
}

func (e AccessScope) IsValid() bool {
	switch e {
	case AccessScopeProfile, AccessScopeTrackers, AccessScopeTickets, AccessScopeACLS, AccessScopeEvents, AccessScopeSubscriptions:
		return true
	}
	return false
}

func (e AccessScope) String() string {
	return string(e)
}

func (e *AccessScope) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessScope(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessScope", str)
	}
	return nil
}

func (e AccessScope) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Authenticity string

const (
	AuthenticityAuthentic   Authenticity = "AUTHENTIC"
	AuthenticityUnauthentic Authenticity = "UNAUTHENTIC"
	AuthenticityTampered    Authenticity = "TAMPERED"
)

var AllAuthenticity = []Authenticity{
	AuthenticityAuthentic,
	AuthenticityUnauthentic,
	AuthenticityTampered,
}

func (e Authenticity) IsValid() bool {
	switch e {
	case AuthenticityAuthentic, AuthenticityUnauthentic, AuthenticityTampered:
		return true
	}
	return false
}

func (e Authenticity) String() string {
	return string(e)
}

func (e *Authenticity) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Authenticity(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Authenticity", str)
	}
	return nil
}

func (e Authenticity) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EventType string

const (
	EventTypeCreated         EventType = "CREATED"
	EventTypeComment         EventType = "COMMENT"
	EventTypeStatusChange    EventType = "STATUS_CHANGE"
	EventTypeLabelAdded      EventType = "LABEL_ADDED"
	EventTypeLabelRemoved    EventType = "LABEL_REMOVED"
	EventTypeAssignedUser    EventType = "ASSIGNED_USER"
	EventTypeUnassignedUser  EventType = "UNASSIGNED_USER"
	EventTypeUserMentioned   EventType = "USER_MENTIONED"
	EventTypeTicketMentioned EventType = "TICKET_MENTIONED"
)

var AllEventType = []EventType{
	EventTypeCreated,
	EventTypeComment,
	EventTypeStatusChange,
	EventTypeLabelAdded,
	EventTypeLabelRemoved,
	EventTypeAssignedUser,
	EventTypeUnassignedUser,
	EventTypeUserMentioned,
	EventTypeTicketMentioned,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeCreated, EventTypeComment, EventTypeStatusChange, EventTypeLabelAdded, EventTypeLabelRemoved, EventTypeAssignedUser, EventTypeUnassignedUser, EventTypeUserMentioned, EventTypeTicketMentioned:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

func (e *EventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventType", str)
	}
	return nil
}

func (e EventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TicketResolution string

const (
	TicketResolutionUnresolved  TicketResolution = "UNRESOLVED"
	TicketResolutionFixed       TicketResolution = "FIXED"
	TicketResolutionImplemented TicketResolution = "IMPLEMENTED"
	TicketResolutionWontFix     TicketResolution = "WONT_FIX"
	TicketResolutionByDesign    TicketResolution = "BY_DESIGN"
	TicketResolutionInvalid     TicketResolution = "INVALID"
	TicketResolutionDuplicate   TicketResolution = "DUPLICATE"
	TicketResolutionNotOurBug   TicketResolution = "NOT_OUR_BUG"
)

var AllTicketResolution = []TicketResolution{
	TicketResolutionUnresolved,
	TicketResolutionFixed,
	TicketResolutionImplemented,
	TicketResolutionWontFix,
	TicketResolutionByDesign,
	TicketResolutionInvalid,
	TicketResolutionDuplicate,
	TicketResolutionNotOurBug,
}

func (e TicketResolution) IsValid() bool {
	switch e {
	case TicketResolutionUnresolved, TicketResolutionFixed, TicketResolutionImplemented, TicketResolutionWontFix, TicketResolutionByDesign, TicketResolutionInvalid, TicketResolutionDuplicate, TicketResolutionNotOurBug:
		return true
	}
	return false
}

func (e TicketResolution) String() string {
	return string(e)
}

func (e *TicketResolution) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TicketResolution(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TicketResolution", str)
	}
	return nil
}

func (e TicketResolution) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TicketStatus string

const (
	TicketStatusReported   TicketStatus = "REPORTED"
	TicketStatusConfirmed  TicketStatus = "CONFIRMED"
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	TicketStatusPending    TicketStatus = "PENDING"
	TicketStatusResolved   TicketStatus = "RESOLVED"
)

var AllTicketStatus = []TicketStatus{
	TicketStatusReported,
	TicketStatusConfirmed,
	TicketStatusInProgress,
	TicketStatusPending,
	TicketStatusResolved,
}

func (e TicketStatus) IsValid() bool {
	switch e {
	case TicketStatusReported, TicketStatusConfirmed, TicketStatusInProgress, TicketStatusPending, TicketStatusResolved:
		return true
	}
	return false
}

func (e TicketStatus) String() string {
	return string(e)
}

func (e *TicketStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TicketStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TicketStatus", str)
	}
	return nil
}

func (e TicketStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
