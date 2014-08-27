// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package pkgtrack provides functions to extract shipment tracking numbers of
// various logistics companies and fetch their delivery status from the web.
// Currently supported are: DHL, UPS
package pkgtrack

import (
	"errors"
	"regexp"
)

var (
	ErrNotAvailable = errors.New("Package tracking not available for this company")
)

// Company describes a logistics company, their tracking numbers and web service.
type Company interface {
	Name() string                    // Name of the company
	Find(body string) []Package      // Searches a string for tracking numbers
	NumberRegexps() []*regexp.Regexp // Returns a list of regexps for the tracking numbers

	IsTrackingAvailable() bool               // Whether online package tracking is available
	TrackingUrl(p Package) string            // Returns a URL to track a specific number
	Track(p Package) (DeliveryStatus, error) // Fetches the latest delivery status from the web
}

var companies = []Company{
	DHL,
	UPS,
}

// Package describes a single shipment.
type Package struct {
	Number  string
	Company Company
}

// Track fetches the current delivery status from the web.
func (p Package) Track() (DeliveryStatus, error) {
	return p.Company.Track(p)
}

// TrackingUrl returns the specific tracking URL for this package.
func (p Package) TrackingUrl() string {
	return p.Company.TrackingUrl(p)
}

// Find searches a text for tracking numbers of all supported logistics companies.
func Find(body string) []Package {
	pkgs := make([]Package, 0)
	for _, c := range companies {
		pkgs = append(pkgs, c.Find(body)...)
	}
	return pkgs
}

// RegisterCompany registers a new logistics company for use in the
// default set of functions.
func RegisterCompany(c Company) {
	companies = append(companies, c)
}

func findRegexp(c Company, body string) []Package {
	pkgs := make([]Package, 0)

	for _, re := range c.NumberRegexps() {
		extracted := re.FindAllStringSubmatch(body, -1)
		for _, num := range extracted {
			pkgs = append(pkgs, Package{
				Number:  num[len(num)-1],
				Company: c,
			})
		}
	}

	return pkgs
}
