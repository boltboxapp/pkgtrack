// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pkgtrack

import "testing"

type parse struct {
	Body     string
	Expected []Package
}

func Contains(s []Package, e Package) bool {
	for _, v := range s {
		if v.Number == e.Number && v.Company == e.Company {
			return true
		}
	}
	return false
}

func TestParse(t *testing.T) {
	tests := []parse{
		{
			"Paketverfolgungsnummer JJD000390003018138123. Blah blah",
			[]Package{
				{Number: "JJD000390003018138123", Company: DHL},
			},
		},
		{
			"Paketverfolgungsnummer 314227130123. Blah blah",
			[]Package{
				{Number: "314227130123", Company: DHL},
			},
		},
		{
			"Kontrollnummer: 1ZA0Y3176800365123. Blah blah",
			[]Package{
				{Number: "1ZA0Y3176800365123", Company: UPS},
			},
		},
	}

	for _, test := range tests {
		numbers := Find(test.Body)
		if len(numbers) != len(test.Expected) {
			t.Errorf("expected %d tracking numbers, got %d", len(test.Expected), len(numbers))
		}
		t.Log(numbers)
		for _, exp := range test.Expected {
			if !Contains(numbers, exp) {
				t.Errorf("tracking number not found: %v", exp)
			}
		}
	}
}

func TestTracking(t *testing.T) {
	tests := []string{
		"JJD000390003018138895",
		"JJD000390003605245565",
		"1ZA0Y3176800365349",
	}

	for _, str := range tests {
		nums := Find(str)
		if len(nums) != 1 {
			t.Errorf("expected %d tracking numbers, got %d", 1, len(nums))
		}
		t.Log(nums)
		s, err := nums[0].Track()
		t.Log(s)
		if err != nil {
			t.Error("tracking error: ", err)
		}
	}
}
