// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pkgtrack

import "time"

type EventType string

const (
	EventInstructionData   EventType = "Instruction data provided"
	EventCenterOrigin      EventType = "Processed in the parcel center of origin"
	EventCenterDestination EventType = "Processed in the parcel center of destination"
	EventDelivery          EventType = "On delivery"
	EventReady             EventType = "Ready for pick-up"
	EventDelivered         EventType = "Package was delivered"
	EventUnknown           EventType = "Unknown event"
)

// DeliveryEvent contains a discrete event in the shipping progress of a package.
type DeliveryEvent struct {
	Time     time.Time
	Location string
	Type     EventType
	Text     string
}

// DeliveryStatus describes the shipment progress of a package.
type DeliveryStatus struct {
	Events []DeliveryEvent
}

// A type that sorts []DeliveryEvent by ascending time
type EventsByTime []DeliveryEvent

func (a EventsByTime) Len() int           { return len(a) }
func (a EventsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EventsByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
