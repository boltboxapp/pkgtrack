// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pkgtrack_test

import (
	"fmt"

	"github.com/xconstruct/pkgtrack"
)

func ExampleFind() {
	packages := pkgtrack.Find("my text with a number 1ZA0Y3176800365123")
	if len(packages) == 0 {
		fmt.Println("no numbers found")
		return
	}

	for _, pkg := range packages {
		fmt.Println("number:", pkg.Number)
		fmt.Println("company:", pkg.Company.Name())
	}
	// Output:
	// number: 1ZA0Y3176800365123
	// company: UPS
}

func ExampleFind_specificCompany() {
	packages := pkgtrack.UPS.Find("my text with a number 1ZA0Y3176800365123")
	if len(packages) == 0 {
		fmt.Println("no numbers found")
		return
	}

	for _, pkg := range packages {
		fmt.Println("number:", pkg.Number)
		fmt.Println("company:", pkg.Company.Name())
	}
	// Output:
	// number: 1ZA0Y3176800365123
	// company: UPS
}

func ExamplePackage_Track() {
	packages := pkgtrack.Find("your tracking number is 1ZA0Y3176800365349")
	for i, pkg := range packages {
		fmt.Println("Package", i+1)
		fmt.Println("number:", pkg.Number)
		fmt.Println("company:", pkg.Company.Name())

		status, err := pkg.Track()
		if err != nil {
			fmt.Println("error:", err)
		}
		lastEvent := status.Events[len(status.Events)-1]
		fmt.Println("last status:", lastEvent.Type)
	}
	// Output:
	// Package 1
	// number: 1ZA0Y3176800365349
	// company: UPS
	// last status: Package was delivered
}
