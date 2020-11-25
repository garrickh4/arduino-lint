// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package checkdata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/arduino/arduino-check/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
	"github.com/arduino/arduino-check/result/feedback"
	"github.com/arduino/arduino-cli/arduino/libraries"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/client9/misspell"
	"github.com/ory/jsonschema/v3"
	"github.com/sirupsen/logrus"
)

// Initialize gathers the library check data for the specified project.
func InitializeForLibrary(project project.Type, schemasPath *paths.Path) {
	var err error

	libraryProperties, libraryPropertiesLoadError = libraryproperties.Properties(project.Path)
	if libraryPropertiesLoadError != nil {
		logrus.Errorf("Error loading library.properties from %s: %s", project.Path, libraryPropertiesLoadError)
		// TODO: can I even do this?
		libraryPropertiesSchemaValidationResult = nil
	} else {
		libraryPropertiesSchemaValidationResult = libraryproperties.Validate(libraryProperties, schemasPath)
	}

	loadedLibrary, err = libraries.Load(project.Path, libraries.User)
	if err != nil {
		logrus.Errorf("Error loading library from %s: %s", project.Path, err)
		loadedLibrary = nil
		sourceHeaders = nil
	} else {
		sourceHeaders, err = loadedLibrary.SourceHeaders()
		if err != nil {
			panic(err)
		}
	}

	if libraryManagerIndex == nil { // Only download the Library Manager index once
		url := "http://downloads.arduino.cc/libraries/library_index.json"
		httpResponse, err := http.Get(url)
		if err != nil {
			feedback.Errorf("%s Unable to download Library Manager index from %s", err, url)
			os.Exit(1)
		}
		defer httpResponse.Body.Close()

		bytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &libraryManagerIndex)
		if err != nil {
			panic(err)
		}
	}

	if misspelledWordsReplacer == nil { // The replacer only needs to be compiled once per run.
		misspelledWordsReplacer = misspell.New()
		misspelledWordsReplacer.Compile()
	}
}

var libraryPropertiesLoadError error

// LibraryPropertiesLoadError returns the error output from loading the library.properties metadata file.
func LibraryPropertiesLoadError() error {
	return libraryPropertiesLoadError
}

var libraryProperties *properties.Map

// LibraryProperties returns the data from the library.properties metadata file.
func LibraryProperties() *properties.Map {
	return libraryProperties
}

var libraryPropertiesSchemaValidationResult map[compliancelevel.Type]*jsonschema.ValidationError

// LibraryPropertiesSchemaValidationResult returns the result of validating library.properties against the JSON schema.
func LibraryPropertiesSchemaValidationResult() map[compliancelevel.Type]*jsonschema.ValidationError {
	return libraryPropertiesSchemaValidationResult
}

var loadedLibrary *libraries.Library

// LoadedLibrary returns the library object generated by Arduino CLI.
func LoadedLibrary() *libraries.Library {
	return loadedLibrary
}

var sourceHeaders []string

// SourceHeaders returns the list of library source header filenames discovered by Arduino CLI.
func SourceHeaders() []string {
	return sourceHeaders
}

var libraryManagerIndex map[string]interface{}

// LibraryManagerIndex returns the Library Manager index data.
func LibraryManagerIndex() map[string]interface{} {
	return libraryManagerIndex
}

var misspelledWordsReplacer *misspell.Replacer

// MisspelledWordsReplacer returns the misspelled words replacer used for spell check.
func MisspelledWordsReplacer() *misspell.Replacer {
	return misspelledWordsReplacer
}
