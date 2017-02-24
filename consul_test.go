package consul

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
)

func ExampleLookup() {
	service := "authentication.service.consul"
	r, err := Lookup(service)

	if err != nil {
		fmt.Printf("Error Looking up consul service %q - %s\n", service, err.Error())
	}

	//services-dev-03.node.dc1.consul:8081
	fmt.Println(*r)
}

func TestLookup(t *testing.T) {
	dummyLookup := func(service, proto, name string) (cname string, addrs []*net.SRV, err error) {
		result := net.SRV{
			Target: "node1.example.com.consul.",
			Port:   8089,
		}
		results := []*net.SRV{&result}
		return "", results, nil
	}

	service := "example.com.consul"
	r, err := lookupWihLookupSRV(service, dummyLookup)
	assertNil(t, err, "Error should have been nil")
	assertEqual(t, "node1.example.com.consul:8089", *r)
}

func TestLookup_HandleProtacol(t *testing.T) {
	dummyLookup := func(service, proto, name string) (cname string, addrs []*net.SRV, err error) {
		if name != "example.com.consul" {
			return "", nil, nil
		}
		result := net.SRV{
			Target: "node1.example.com.consul.",
			Port:   8089,
		}
		results := []*net.SRV{&result}
		return "", results, nil
	}

	service := "http://example.com.consul"
	r, err := lookupWihLookupSRV(service, dummyLookup)
	assertNil(t, err, "Error should have been nil")
	assertEqual(t, "http://node1.example.com.consul:8089", *r)
}

func TestLookup_NotConsul(t *testing.T) {
	service := "example.com"
	r, err := Lookup(service)
	assertNil(t, err, "Error should have been nil")
	assertNil(t, r, "Result should have been nil")
}

func TestLookup_InvalidName(t *testing.T) {
	service := ".consul"
	r, err := Lookup(service)
	assertEqual(t, "Error performing SRV DNS Lookup for .consul - lookup : invalid domain name", err.Error())
	assertNil(t, r, "Result should have been nil")
}

func TestLookup_LookupError(t *testing.T) {
	errorLookup := func(service, proto, name string) (cname string, addrs []*net.SRV, err error) {
		return "", nil, errors.New("Boom")
	}

	service := "example.com.consul"
	r, err := lookupWihLookupSRV(service, errorLookup)

	assertEqual(t, "Error performing SRV DNS Lookup for example.com.consul - Boom", err.Error())
	assertNil(t, r, "Result should have been nil")
}

func TestLookup_LookupNilResults(t *testing.T) {
	errorLookup := func(service, proto, name string) (cname string, addrs []*net.SRV, err error) {
		return "", nil, nil
	}

	service := "example.com.consul"
	r, err := lookupWihLookupSRV(service, errorLookup)
	assertNil(t, err, "Error should have been nil")
	assertNil(t, r, "Result should have been nil")
}

func assertNil(t *testing.T, i interface{}, message string) {
	if !isNil(i) {
		t.Fatal(message)
	}
}

func assertNotNil(t *testing.T, i interface{}, message string) {
	if isNil(i) {
		t.Fatal(message)
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

func assertEqual(t *testing.T, expected string, result string) {
	if result != expected {
		t.Fatalf("Expected %q got %q", expected, result)
	}
}
