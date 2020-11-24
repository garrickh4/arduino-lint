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

/*
Package checkconfigurations defines the configuration of each check:
- metadata
- output template
- under which conditions it's enabled
- the level of a failure
- which function implements it
*/
package checkconfigurations

import (
	"github.com/arduino/arduino-check/check/checkfunctions"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
)

// Type is the type for check configurations.
type Type struct {
	ProjectType projecttype.Type // The project type the check applies to.
	// The following fields provide arbitrary text for the tool output associated with each check:
	Category        string
	Subcategory     string
	ID              string // Unique check identifier: <project type identifier (L|S|P|I)><category identifier><number>
	Brief           string // Short description of the check.
	Description     string // Supplemental information about the check.
	MessageTemplate string // The warning/error message template displayed when the check fails. Will be filled by check function output.
	// The following fields define under which tool configuration modes the check will run:
	DisableModes []checkmode.Type // Check is disabled when tool is in any of these modes.
	EnableModes  []checkmode.Type // Check is only enabled when tool is in one of these modes.
	// The following fields define the check level in each configuration mode:
	InfoModes     []checkmode.Type    // Failure of the check only results in an informational message.
	WarningModes  []checkmode.Type    // Failure of the check is considered a warning.
	ErrorModes    []checkmode.Type    // Failure of the check is considered an error.
	CheckFunction checkfunctions.Type // The function that implements the check.
}

// Configurations returns the slice of check configurations.
func Configurations() []Type {
	return configurations
}

// configurations is an array of structs that define the configuration of each check.
var configurations = []Type{
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "general",
		ID:              "LP001",
		Brief:           "invalid format",
		Description:     "",
		MessageTemplate: "library.properties has an invalid format: {{.}}",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesFormat,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP002",
		Brief:           "missing name field",
		Description:     "",
		MessageTemplate: "missing name field in library.properties",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldMissing,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP003",
		Brief:           "disallowed characters",
		Description:     "",
		MessageTemplate: "disallowed characters in library.properties name field. See: https://arduino.github.io/arduino-cli/latest/library-specification/#libraryproperties-file-format",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldDisallowedCharacters,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP005",
		Brief:           "duplicate name",
		Description:     "This requirement only applies to the library.properties name value. There is no requirement to change the repository or header file names.",
		MessageTemplate: "Library name {{.}} is in use by a library in the Library Manager index. Each library must have a unique name value.",
		DisableModes:    []checkmode.Type{checkmode.LibraryManagerIndexed},
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.All},
		ErrorModes:      []checkmode.Type{checkmode.LibraryManagerSubmission},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldDuplicate,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP006",
		Brief:           "not in LM index",
		Description:     "The name value is the identifier used to install the library and define dependencies, so it should not be changed.",
		MessageTemplate: "Library name {{.}} not found in the Library Manager index. Library names are not allowed to change after being added to the index. See: https://github.com/arduino/Arduino/wiki/Library-Manager-FAQ#how-can-i-change-my-librarys-name",
		DisableModes:    []checkmode.Type{checkmode.Default},
		EnableModes:     []checkmode.Type{checkmode.LibraryManagerIndexed},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldNotInIndex,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "version field",
		ID:              "LP004",
		Brief:           "missing version field",
		Description:     "",
		MessageTemplate: "missing version field in library.properties",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesVersionFieldMissing,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "version field",
		ID:              "",
		Brief:           "invalid",
		Description:     `Must be compliant with "relaxed semver".`,
		MessageTemplate: "library.properties version value {{.}} is invalid. See https://arduino.github.io/arduino-cli/latest/library-specification/#libraryproperties-file-format",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesVersionFieldNonRelaxedSemver,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "version field",
		ID:              "",
		Brief:           "non-semver",
		Description:     "",
		MessageTemplate: "library.properties version value {{.}} is not compliant with the semver specification. See https://semver.org/",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.All},
		ErrorModes:      nil,
		CheckFunction:   checkfunctions.LibraryPropertiesVersionFieldNonSemver,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "author field",
		ID:              "",
		Brief:           "missing author field",
		Description:     "",
		MessageTemplate: "missing required author field in library.properties. See https://arduino.github.io/arduino-cli/latest/library-specification/#libraryproperties-file-format",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesAuthorFieldMissing,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "author field",
		ID:              "",
		Brief:           "author < min length",
		Description:     "",
		MessageTemplate: "library.properties author value is less than minimum length",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.All},
		CheckFunction:   checkfunctions.LibraryPropertiesAuthorFieldLTMinLength,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "sentence field",
		ID:              "",
		Brief:           "sentence spell check",
		Description:     "",
		MessageTemplate: "A commonly misspelled word was found in the library.properties sentence field. Suggested correction: {{.}}",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.All},
		ErrorModes:      nil,
		CheckFunction:   checkfunctions.LibraryPropertiesSentenceFieldSpellCheck,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "paragraph field",
		ID:              "",
		Brief:           "paragraph spell check",
		Description:     "",
		MessageTemplate: "A commonly misspelled word was found in the library.properties paragraph field. Suggested correction: {{.}}",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.All},
		ErrorModes:      nil,
		CheckFunction:   checkfunctions.LibraryPropertiesParagraphFieldSpellCheck,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "depends field",
		ID:              "LP012",
		Brief:           "Dependency not in index",
		Description:     "",
		MessageTemplate: "library.properties depends field item(s) {{.}} not found in the Library Manager index.",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.All},
		ErrorModes:      nil,
		CheckFunction:   checkfunctions.LibraryPropertiesDependsFieldNotInIndex,
	},
	{
		ProjectType:     projecttype.Sketch,
		Category:        "structure",
		Subcategory:     "",
		ID:              "SS001",
		Brief:           ".pde extension",
		Description:     "The .pde extension is used by both Processing sketches and Arduino sketches. Processing sketches should either be in the \"data\" subfolder of the sketch or in the \"extras\" folder of the library. Arduino sketches should use the modern .ino extension",
		MessageTemplate: "{{.}} uses deprecated .pde file extension. Use .ino for Arduino sketches",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.All},
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.Permissive},
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.PdeSketchExtension,
	},
}