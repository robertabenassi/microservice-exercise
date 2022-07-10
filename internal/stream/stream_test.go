package stream

import (
	"fmt"
	"microservice-exercise/internal/data_model"
	"os"
	"testing"
)

// TestStreamReadingOK is green when everything goes fine
func TestStreamReadingOK(t *testing.T) {
	// expected results in a map format
	var testPorts = map[string]data_model.Port{
		"AEAJM": {
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
			Alias:   []string{},
			Regions: []string{},
			Coordinates: []float32{
				55.5136433,
				25.4052165,
			},
			Province: "Ajman",
			Timezone: "Asia/Dubai",
			Unlocs: []string{
				"AEAJM",
			},
			Code: "52000",
		},
		"ZABFN": {
			Name:    "Bloemfontein",
			City:    "Bloemfontein",
			Country: "South Africa",
			Alias:   []string{},
			Regions: []string{},
			Coordinates: []float32{
				26.1595761,
				-29.085214,
			},
			Province: "Free State",
			Timezone: "Africa/Johannesburg",
			Unlocs: []string{
				"AEAUH",
			},
			Code: "52001",
		},
	}

	stream := NewStream()

	// in order to test we collect errors into a channel.
	errs := make(chan error, 1)

	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errs <- data.Error
			}

			// it is a pity not to be able to check the whole struct to be equal
			// I decided to add this simple test just to show the intention
			// The point is how to check if two data struct are equals in terms of value/content
			// I probably got this by a basic reflection, or may be new update will come in go generics (I am not so updated about Go versions!)
			if testPorts[data.ID].Name != data.Port.Name {
				errs <- fmt.Errorf("Expected name of the port was %s, got %s", testPorts[data.ID].Name, data.Port.Name)
			}
			errs <- nil // got no error, so put nil in the channel
		}
	}()
	pwd, _ := os.Getwd()
	stream.Load(pwd + "/testdata/testdata.json")

	// wait for errors in the test
	err := <-errs
	if err != nil {
		t.Fatal(err)
	}
}

// TestStreamReadingKO should be green when
func TestStreamReadingKO(t *testing.T) {
	// expected results in a map format
	var testPorts = map[string]data_model.Port{
		"AEAJM": {
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
			Alias:   []string{},
			Regions: []string{},
			Coordinates: []float32{
				55.5136433,
				25.4052165,
			},
			Province: "Ajman",
			Timezone: "Asia/Dubai",
			Unlocs: []string{
				"AEAJM",
			},
			Code: "52000",
		},
		"AAAA": {
			Name:    "Bloemfontein",
			City:    "Bloemfontein",
			Country: "South Africa",
			Alias:   []string{},
			Regions: []string{},
			Coordinates: []float32{
				26.1595761,
				-29.085214,
			},
			Province: "Free State",
			Timezone: "Africa/Johannesburg",
			Unlocs: []string{
				"AAAA",
			},
			Code: "52001",
		},
	}

	stream := NewStream()

	// in order to test we collect errors into a channel.
	errs := make(chan error, 1)

	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				errs <- data.Error
			}
			if testPorts[data.ID].Name != data.Port.Name {
				errs <- fmt.Errorf("Expected name of the port was %s, got %s", testPorts[data.ID].Name, data.Port.Name)
			}
			errs <- nil // got no error, so put nil in the channel
		}
	}()
	pwd, _ := os.Getwd()
	stream.Load(pwd + "/testdata/testdata.json")

	// wait for errors in the test
	err := <-errs
	if err == nil {
		t.Fatal("An error was expected, but got no error")
	}
}
