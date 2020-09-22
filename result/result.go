// Package result records check results and provides reports and summary text on those results.
package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checklevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result/feedback"
	"github.com/arduino/go-paths-helper"
)

type reportType struct {
	Configuration toolConfigurationReportType `json:"configuration"`
	Projects      []projectReportType         `json:"projects"`
	Summary       summaryReportType           `json:"summary"`
}

type toolConfigurationReportType struct {
	Paths       []*paths.Path `json:"paths"`
	ProjectType string        `json:"projectType"`
	Recursive   bool          `json:"recursive"`
}

type projectReportType struct {
	Path          *paths.Path                    `json:"path"`
	ProjectType   string                         `json:"projectType"`
	Configuration projectConfigurationReportType `json:"configuration"`
	Checks        []checkReportType              `json:"checks"`
	Summary       summaryReportType              `json:"summary"`
}

type projectConfigurationReportType struct {
	Permissive           bool `json:"permissive"`
	LibraryManagerSubmit bool `json:"libraryManagerSubmit"`
	LibraryManagerUpdate bool `json:"libraryManagerUpdate"`
	Official             bool `json:"official"`
}

type checkReportType struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	ID          string `json:"ID"`
	Brief       string `json:"brief"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

type summaryReportType struct {
	Pass         bool `json:"pass"`
	WarningCount int  `json:"warningCount"`
	ErrorCount   int  `json:"errorCount"`
}

var report reportType

// Initialize adds the tool configuration data to the report.
func Initialize() {
	report.Configuration = toolConfigurationReportType{
		Paths:       []*paths.Path{configuration.TargetPath()},
		ProjectType: configuration.SuperprojectTypeFilter().String(),
		Recursive:   configuration.Recursive(),
	}
}

// Record records the result of a check and returns a text summary for it.
func Record(checkedProject project.Type, checkConfiguration checkconfigurations.Type, checkResult checkresult.Type, checkOutput string) string {
	checkMessage := message(checkConfiguration.MessageTemplate, checkOutput)

	checkLevel, err := checklevel.CheckLevel(checkConfiguration)
	if err != nil {
		feedback.Errorf("Error while determining check level: %v", err)
		os.Exit(1)
	}

	summaryText := fmt.Sprintf("%v\n", checkResult.String())

	if checkResult == checkresult.NotRun {
		// TODO: make the check functions output an explanation for why they didn't run
		summaryText += fmt.Sprintf("%s: %s\n", checklevel.Notice.String(), checkOutput)
	} else if checkResult != checkresult.Pass {
		summaryText += fmt.Sprintf("%s: %s\n", checkLevel.String(), checkMessage)
	}

	checkReport := checkReportType{
		Category:    checkConfiguration.Category,
		Subcategory: checkConfiguration.Subcategory,
		ID:          checkConfiguration.ID,
		Brief:       checkConfiguration.Brief,
		Description: checkConfiguration.Description,
		Result:      checkResult.String(),
		Level:       checkLevel.String(),
		Message:     checkMessage,
	}

	reportExists, projectReportIndex := getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		// There is no existing report for this project.
		report.Projects = append(
			report.Projects,
			projectReportType{
				Path:        checkedProject.Path,
				ProjectType: checkedProject.ProjectType.String(),
				Configuration: projectConfigurationReportType{
					Permissive:           configuration.CheckModes(checkedProject.ProjectType)[checkmode.Permissive],
					LibraryManagerSubmit: configuration.CheckModes(checkedProject.ProjectType)[checkmode.Permissive],
					LibraryManagerUpdate: configuration.CheckModes(checkedProject.ProjectType)[checkmode.LibraryManagerIndexed],
					Official:             configuration.CheckModes(checkedProject.ProjectType)[checkmode.Official],
				},
				Checks: []checkReportType{checkReport},
			},
		)
	} else {
		// There's already a report for this project, just add the checks report to it
		report.Projects[projectReportIndex].Checks = append(report.Projects[projectReportIndex].Checks, checkReport)
	}

	return summaryText
}

// AddProjectSummaryReport summarizes the results of all checks on the given project and adds it to the report.
func AddProjectSummaryReport(checkedProject project.Type) {
	reportExists, projectReportIndex := getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary", checkedProject.Path))
	}

	pass := true
	warningCount := 0
	errorCount := 0
	for _, checkReport := range report.Projects[projectReportIndex].Checks {
		if checkReport.Result == checkresult.Fail.String() {
			if checkReport.Level == checklevel.Warning.String() {
				warningCount += 1
			} else if checkReport.Level == checklevel.Error.String() {
				errorCount += 1
				pass = false
			}
		}
	}

	report.Projects[projectReportIndex].Summary = summaryReportType{
		Pass:         pass,
		WarningCount: warningCount,
		ErrorCount:   errorCount,
	}
}

// ProjectSummaryText returns a text summary of the check results for the given project.
func ProjectSummaryText(checkedProject project.Type) string {
	reportExists, projectReportIndex := getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary text", checkedProject.Path))
	}

	projectSummaryReport := report.Projects[projectReportIndex].Summary
	return fmt.Sprintf("\nFinished checking project. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v\n\n", projectSummaryReport.WarningCount, projectSummaryReport.ErrorCount, projectSummaryReport.Pass)
}

// AddSummaryReport summarizes the check results for all projects and adds it to the report.
func AddSummaryReport() {
	pass := true
	warningCount := 0
	errorCount := 0
	for _, projectReport := range report.Projects {
		if !projectReport.Summary.Pass {
			pass = false
		}
		warningCount += projectReport.Summary.WarningCount
		errorCount += projectReport.Summary.ErrorCount
	}

	report.Summary = summaryReportType{
		Pass:         pass,
		WarningCount: warningCount,
		ErrorCount:   errorCount,
	}
}

// SummaryText returns a text summary of the cumulative check results.
func SummaryText() string {
	return fmt.Sprintf("Finished checking projects. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v\n", report.Summary.WarningCount, report.Summary.ErrorCount, report.Summary.Pass)
}

// Report returns a JSON formatted report of checks on all projects.
func JSONReport() string {
	return string(jsonReportRaw())
}

func jsonReportRaw() []byte {
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Error while formatting checks report: %v", err))
	}

	return reportJSON
}

// WriteReport writes a report for all projects to the specified file.
func WriteReport() {
	// Write report file
	err := configuration.ReportFilePath().WriteFile(jsonReportRaw())
	if err != nil {
		feedback.Errorf("Error while writing report: %v", err)
		os.Exit(1)
	}
}

// Passed returns whether the checks passed cumulatively.
func Passed() bool {
	return report.Summary.Pass
}

func getProjectReportIndex(projectPath *paths.Path) (bool, int) {
	var index int
	var projectReport projectReportType
	for index, projectReport = range report.Projects {
		if projectReport.Path == projectPath {
			return true, index
		}
	}

	// There is no element in the report for this project.
	return false, index + 1
}

// message fills the message template provided by the check configuration with the check output.
// TODO: make checkOutput a struct to allow for more advanced message templating
func message(templateText string, checkOutput string) string {
	messageTemplate := template.Must(template.New("messageTemplate").Parse(templateText))

	messageBuffer := new(bytes.Buffer)
	messageTemplate.Execute(messageBuffer, checkOutput)

	return messageBuffer.String()
}
