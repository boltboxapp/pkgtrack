// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pkgtrack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	dhl_re_number1    = regexp.MustCompile(`\b(JJD\d{18})\b`)
	dhl_re_number2    = regexp.MustCompile(`\b(\d{12})\b`)
	dhl_re_events     = regexp.MustCompile(`(?sU)id="events">.*<table>(.*)</table>`)
	dhl_re_events_row = regexp.MustCompile(`(?sU)<tr>\s*<td.*>(.*)</td>\s*<td.*>(.*)</td>\s*<td.*>(.*)</td>\s*</tr>`)

	dhl_tracking_url = "https://nolp.dhl.de/nextt-online-public/set_identcodes.do?lang=en&idc=%s"

	DHL Company = &dhlCompany{}
)

type dhlCompany struct{}

func (c *dhlCompany) Name() string {
	return "DHL"
}

func (c *dhlCompany) NumberRegexps() []*regexp.Regexp {
	return []*regexp.Regexp{
		dhl_re_number1,
		dhl_re_number2,
	}
}

func (c *dhlCompany) IsTrackingAvailable() bool {
	return true
}

func (c *dhlCompany) Find(body string) []Package {
	return findRegexp(c, body)
}

func (c *dhlCompany) TrackingUrl(p Package) string {
	return fmt.Sprintf(dhl_tracking_url, p.Number)
}

func (c *dhlCompany) Track(p Package) (DeliveryStatus, error) {
	s := DeliveryStatus{}

	resp, err := http.Get(c.TrackingUrl(p))
	if err != nil {
		return s, err
	}
	if resp.StatusCode != 200 {
		return s, errors.New("Unexpected status: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return s, err
	}

	eventsTable := dhl_re_events.FindStringSubmatch(string(body))
	if eventsTable == nil {
		return s, errors.New("Parse error: events table not found")
	}

	eventRows := dhl_re_events_row.FindAllStringSubmatch(eventsTable[0], -1)
	s.Events = make([]DeliveryEvent, len(eventRows))
	for i, row := range eventRows {
		s.Events[i], err = c.parseEvent(
			strings.TrimSpace(row[1]),
			strings.TrimSpace(row[2]),
			strings.TrimSpace(row[3]),
		)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func (c *dhlCompany) parseEvent(timestamp, location, status string) (e DeliveryEvent, err error) {
	e.Text = status
	e.Location = location
	if e.Location == "--" {
		e.Location = ""
	}

	if e.Time, err = time.Parse("Mon, 02.01.2006 15:04 h", timestamp); err != nil {
		return e, err
	}

	switch {
	case strings.Contains(status, "instruction data for this shipment have been provided"):
		e.Type = EventInstructionData
	case strings.Contains(status, "processed in the parcel center of origin"):
		e.Type = EventCenterOrigin
	case strings.Contains(status, "processed in the destination parcel center"):
		e.Type = EventCenterDestination
	case strings.Contains(status, "loaded onto the delivery vehicle"):
		e.Type = EventDelivery
	case strings.Contains(status, "on its way to the PACKSTATION"):
		e.Type = EventDelivery
	case strings.Contains(status, "ready for pick-up"):
		e.Type = EventReady
	case strings.Contains(status, "recipient has picked up"):
		e.Type = EventDelivered
	case strings.Contains(status, "successfully delivered"):
		e.Type = EventDelivered
	default:
		return e, errors.New("Unknown event status: " + status)
	}

	return e, err
}
