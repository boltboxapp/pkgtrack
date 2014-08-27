pkgtrack
========

[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](http://godoc.org/github.com/xconstruct/pkgtrack)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT)

A Go package to extract shipment tracking numbers of various logistics companies and fetch their delivery status from the web.
of various logistics companies.

Currently supported are:
* DHL
* UPS

Install
-------

```
go get "github.com/xconstruct/pkgtrack"
```

Example
-------

```go
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
```
