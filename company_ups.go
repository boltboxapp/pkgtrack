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
	"sort"
	"strings"
	"time"
)

var (
	ups_re_number1    = regexp.MustCompile(`\b(1Z ?[0-9A-Z]{3} ?[0-9A-Z]{3} ?[0-9A-Z]{2} ?[0-9A-Z]{4} ?[0-9A-Z]{3} ?[0-9A-Z]|[\dT]\d{3} ?\d{4} ?\d{3})\b`)
	ups_re_events     = regexp.MustCompile(`(?sU)class="dataTable">(.*)</table>`)
	ups_re_events_row = regexp.MustCompile(`(?sU)<tr.*>\s*<td.*>(.*)</td>\s*<td.*>(.*)</td>\s*<td.*>(.*)</td>\s*<td.*>(.*)</td>\s*</tr>`)
	ups_re_multispace = regexp.MustCompile(`(?s)\s+`)

	ups_tracking_url = "http://wwwapps.ups.com/WebTracking/track?track=yes&trackNums=%s"

	UPS Company = &upsCompany{}
)

type upsCompany struct{}

func (c *upsCompany) Name() string {
	return "UPS"
}

func (c *upsCompany) NumberRegexps() []*regexp.Regexp {
	return []*regexp.Regexp{
		ups_re_number1,
	}
}

func (c *upsCompany) IsTrackingAvailable() bool {
	return true
}

func (c *upsCompany) Find(body string) []Package {
	return findRegexp(c, body)
}

func (c *upsCompany) TrackingUrl(p Package) string {
	return fmt.Sprintf(ups_tracking_url, p.Number)
}

func (c *upsCompany) Track(p Package) (DeliveryStatus, error) {
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

	eventsTable := ups_re_events.FindStringSubmatch(string(body))
	if eventsTable == nil {
		return s, errors.New("Parse error: events table not found")
	}

	eventRows := ups_re_events_row.FindAllStringSubmatch(eventsTable[0], -1)
	s.Events = make([]DeliveryEvent, len(eventRows))
	for i, row := range eventRows {
		s.Events[i], err = c.parseEvent(
			strings.TrimSpace(row[2])+" "+strings.TrimSpace(row[3]),
			strings.TrimSpace(row[1]),
			strings.TrimSpace(row[4]),
		)
		if err != nil {
			return s, err
		}
	}
	sort.Sort(EventsByTime(s.Events))

	return s, nil
}

func (c *upsCompany) parseEvent(timestamp, location, status string) (e DeliveryEvent, err error) {
	e.Text = status
	e.Location = ups_re_multispace.ReplaceAllString(location, " ")

	timestamp = strings.Replace(timestamp, ".", "", -1)
	if e.Time, err = time.Parse("01/02/2006 3:04 PM", timestamp); err != nil {
		return e, err
	}

	switch {
	case strings.Contains(status, "Order Processed"):
		e.Type = EventInstructionData
	case strings.Contains(status, "Delivered"):
		e.Type = EventDelivered
	default:
		return e, errors.New("Unknown event status: " + status)
	}

	return e, err
}
