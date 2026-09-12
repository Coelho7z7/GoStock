package models

import "time"

type Movement struct {
	ID            int
	ProductID     int
	UserID        int
	Type          string
	Quantity      int
	Date          time.Time
	Product       string
	User          string
	FormattedDate string
	FormattedTime string
	FormattedType string
}
