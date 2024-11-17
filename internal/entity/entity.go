package entity

import "time"

// User represents a user in the system (either client or contractor).
type User struct {
	ID       string `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	Role     string `db:"role" json:"role"` // 'client' or 'contractor'
	Email    string `db:"email" json:"email"`
}

// Tender represents a tender created by a client.
type Tender struct {
	ID          string    `db:"id" json:"id"`
	ClientID    string    `db:"client_id" json:"client_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	Budget      float64   `db:"budget" json:"budget"`
	Status      string    `db:"status" json:"status"`
}

type TenderReq struct {
	ClientID    string    `db:"client_id" json:"client_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	Budget      float64   `db:"budget" json:"budget"`
}
type TenderReq1 struct {
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	Budget      float64   `db:"budget" json:"budget"`
}

type TenderRepoReq struct {
	ClientID    string    `db:"client_id" json:"client_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	Budget      float64   `db:"budget" json:"budget"`
	Status      string    `db:"status" json:"status"`
}

// Bid represents a bid made by a contractor for a tender.
type Bid struct {
	ID           string  `db:"id" json:"id"`
	TenderID     string  `db:"tender_id" json:"tender_id"`
	ContractorID string  `db:"contractor_id" json:"contractor_id"`
	Price        float64 `db:"price" json:"price"`
	DeliveryTime int     `db:"delivery_time" json:"delivery_time"` // in days
	Comments     string  `db:"comments" json:"comments"`
	Status       string  `db:"status" json:"status"` // e.g., 'pending'
}

type BidReq struct {
	TenderID     string  `db:"tender_id" json:"tender_id"`
	ContractorID string  `db:"contractor_id" json:"contractor_id"`
	Price        float64 `db:"price" json:"price"`
	DeliveryTime int     `db:"delivery_time" json:"delivery_time"` // in days
	Comments     string  `db:"comments" json:"comments"`
	Status       string  `db:"status" json:"status"` // e.g., 'pending'
}

type Bid1 struct {
	Price        float64 `db:"price" json:"price"`
	DeliveryTime int     `db:"delivery_time" json:"delivery_time"` // in days
	Comments     string  `db:"comments" json:"comments"`
	Status       string  `db:"status" json:"status"` // e.g., 'pending'
}

// Notification represents a notification for a user.
type Notification struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	Message    string    `db:"message" json:"message"`
	RelationID string    `db:"relation_id" json:"relation_id"` // can refer to related tender/bid ID
	Type       string    `db:"type" json:"type"`               // notification type
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Error struct {
	Error string `db:"error" json:"error"`
}
type Error1 struct {
	Status  int    `db:"status" json:"status"`
	Message string `db:"message" json:"message"`
}
type RegisterReq struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
	Role     string `db:"role" json:"role"`
}

type RegisterRes struct {
	UserId   string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
}

type LogInReq struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type LogInRes struct {
	Token    string `db:"token" json:"token"`
	UserId   string `db:"user_id" json:"user_id"`
	ExpireAt int    `db:"expire_at" json:"expire_at"`
}

type Message struct {
	Message string `db:"message" json:"message"`
}

type StatusRequest struct {
	Status string `db:"status" json:"status"`
}

type ListBidReq struct {
	ClientID           string   `db:"client_id" json:"client_id"`
	TenderID           string   `db:"tender_id" json:"tender_id"`
	PriceFilter        *float64 `db:"price_filter" json:"price_filter,omitempty"`
	DeliveryTimeFilter *int     `db:"delivery_time_filter" json:"delivery_time_filter,omitempty"`
	Comments           string   `db:"comments" json:"comments"`
	Status             string   `db:"status" json:"status"`
}

type UpdateTender struct {
	Id          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	Budget      float64   `db:"budget" json:"budget"`
	Status      string    `db:"status" json:"status"`
}

type Awarded struct {
	TenderID string `db:"tender_id" json:"tender_id"`
	BideId   string `db:"bide_id" json:"bide_id"`
}

type AwardedRes struct {
	TenderID     string `db:"tender_id" json:"tender_id"`
	BideId       string `db:"bide_id" json:"bide_id"`
	ContractorID string `db:"contractor_id" json:"contractor_id"`
}
