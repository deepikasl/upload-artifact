package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jfrog/jfrog-pipelines-tasks-sdk-go/tasks"
)

type UploadArtifact struct {
	inputs       Inputs
	resource     *tasks.Resource
	resourcePath string
	runVariables map[string]string
}

type Inputs struct {
	inputResource                 string
	buildInfoResource             string
	autoPublishBuildInfo          string
	forceXrayScan                 string
	failOnScan                    string
	sourcePath                    string
	targetPath                    string
	properties                    string
	regExp                        string
	flat                          string
	module                        string
	deb                           string
	recursive                     string
	dryRun                        string
	symlinks                      string
	explode                       string
	exclusions                    string
	includeDirs                   string
	syncDeletes                   string
}

var (
	readInput    = tasks.GetInput
	readResource = tasks.GetResource
	execute      = exec.Command
)

const MaxNumberOfRetries = 3

func (m *UploadArtifact) readInputs() {
	// Fetch inputs
	i := Inputs{}
	i.inputResource = readInput("inputResource")
	i.buildInfoResource = readInput("buildInfoResource")
	i.autoPublishBuildInfo = readInput("autoPublishBuildInfo")
	i.forceXrayScan = readInput("forceXrayScan")
	i.failOnScan = readInput("failOnScan")
	i.sourcePath = readInput("sourcePath")
	i.targetPath = readInput("targetPath")
	i.properties = readInput("properties")
	i.regExp = readInput("regExp")
	i.flat = readInput("flat")
	i.module = readInput("module")
	i.deb = readInput("deb")
	i.recursive = readInput("recursive")
	i.dryRun = readInput("dryRun")
	i.symlinks = readInput("symlinks")
	i.explode = readInput("explode")
	i.exclusions = readInput("exclusions")
	i.includeDirs = readInput("includeDirs")
	i.syncDeletes = readInput("syncDeletes")
	m.inputs = i

	tasks.Debug(fmt.Sprintf("Received inputs are [%+v]", i))
}

func (m *UploadArtifact) runPreRequisites() {
	err := m.inputs.validateInputs()
	if err != nil {
		haltExecution(err.Error())
	}
	m.verifyJFrogCLIInstallation()
	m.setResource()
}

// validateInputs validates params throw tasks error if any of the params is empty
func (i *Inputs) validateInputs() error {
	inputs := []string{i.inputResource, i.targetPath, }
	for _, input := range inputs {
		if input == "" || len(input) == 0 {
			tasks.Error("One of the mandatory input " + input + " is missing")
			return errors.New("missing mandatory inputs")
		}
	}

	return nil
}

func haltExecution(errMessage string) {
	tasks.Error(errMessage)
	os.Exit(1)
}

// verifyJFrogCLIInstallation verifies jfrog cli installation
func (m *UploadArtifact) verifyJFrogCLIInstallation() {
	cmd := execute("jf", "--version")
	if output, err := cmd.Output(); err != nil {
		haltExecution("Failed to verify jfrog cli installation, make sure jfrog cli v2 is installed: " + err.Error())
	} else {
		tasks.Info(string(output))
	}
}

func (m *UploadArtifact) setResource() {
	resource, err := readResource(m.inputs.inputResource)
	if err != nil {
		haltExecution("Failed to fetch resource using name: " + m.inputs.inputResource)
	}
	m.resource = &resource
	m.resourcePath = m.resource.ResourcePath
}

func (g *UploadArtifact) handleExecution(name string, options ...string) (string, error) {
	var output []byte
	var err error
	cmdExecLocation := g.resource.ResourcePath + g.inputs.executionLocation
	// simulating retry_command functionality here
	cmd := execute(name, options...)
	cmd.Dir = cmdExecLocation
	if len(g.resourcePath) > 0 {
		tasks.Debug("Received resource path is ", g.resourcePath)
		cmd.Dir = g.resourcePath
	}
	if output, err = cmd.CombinedOutput(); err != nil {
		tasks.Debug("Combined Output: => ", string(output))
		tasks.Error("Failed to run " + cmd.String())
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			fmt.Println(exitCode)
		}
		return string(output), err
	} else {
		tasks.Debug(string(output))
	}
	return string(output), nil
}

func (g *UploadArtifact) run() error {
  stepTmpDir := os.Getenv("step_tmp_dir")
  var uploadArtifactPath = append(stepTmpDir, "/ArtifactUpload")
	_, err := g.handleExecution("mkdir", "-p", uploadArtifactPath)
	if err != nil {
		return err
	}
	_, err := g.handleExecution("pushd", uploadArtifactPath)
	if err != nil {
    return err
  }
  _, err := g.handleExecution("cp", "-r",  g.resourcePath)
	if err != nil {
    return err
  }
  if len(g.inputs.sourcePath) == 0 {
    if g.inputs.regExp == "true" {
      g.inputs.sourcePath=".*"
    } else {
      g.inputs.sourcePath="*"
    }
  }
  if len(g.inputs.targetPath) == 0 {
    haltExecution("Failed to create path for outputLocation" + g.inputs.outputLocation)
  }
  var parameters=""

  if len(g.inputs.module) > 0 {
    parameters = append(parameters, "--module", g.inputs.module)
  }

  var uploadProperties=""
  if len(g.inputs.properties) > 0 {
    uploadProperties= g.inputs.properties
    uploadProperties = append(uploadProperties, ";")
  }
  uploadProperties = append(uploadProperties, "pipelines_step_name=", os.Getenv(step_name), ";")
  uploadProperties = append(uploadProperties, "pipelines_run_number=", os.Getenv(run_number), ";")
  uploadProperties = append(uploadProperties, "pipelines_step_id=", os.Getenv(step_id), ";")
  uploadProperties = append(uploadProperties, "pipelines_pipeline_name=", os.Getenv(pipeline_name), ";")
  uploadProperties = append(uploadProperties, "pipelines_step_type=", os.Getenv(step_type), ";")
  uploadProperties = append(uploadProperties, "pipelines_step_platform=", os.Getenv(step_platform), ";")

  parameters = append(parameters, "--props='", uploadProperties, "'")

  if len(g.inputs.deb) > 0 {
    parameters = append(parameters, "--deb", g.inputs.deb)
  }

  if len(g.inputs.flat) > 0 {
    parameters = append(parameters, "--flat", g.inputs.flat)
  }

  if len(g.inputs.recursive) > 0 {
    parameters = append(parameters, "--recursive", g.inputs.recursive)
  }

  if len(g.inputs.regExp) > 0 {
    parameters = append(parameters, "--regExp", g.inputs.regExp)
  }

  if len(g.inputs.dryRun) > 0 {
    parameters = append(parameters, "--dry-run", g.inputs.dryRun)
  }

  if len(g.inputs.symlinks) > 0 {
    parameters = append(parameters, "--symlinks", g.inputs.symlinks)
  }

  if len(g.inputs.explode) > 0 {
    parameters = append(parameters, "--explode", g.inputs.explode)
  }

  if len(g.inputs.includeDirs) > 0 {
    parameters = append(parameters, "--include-dirs", g.inputs.includeDirs)
  }

  if len(g.inputs.exclusions) > 0 {
    parameters = append(parameters, "--exclusions", g.inputs.exclusions)
  }

  if len(g.inputs.syncDeletes) > 0 {
    parameters = append(parameters, "--sync-deletes", g.inputs.syncDeletes)
  }

  var uploadCommand = ""
  uploadCommand = append(uploadCommand, "\"", g.inputs.sourcePath, "\"")
  uploadCommand = append(uploadCommand, "\"", g.inputs.targetPath, "\"")
  uploadCommand = append(uploadCommand, parameters)
  uploadCommand = append(uploadCommand, "--insecure-tls=", os.Getenv(no_verify_ssl))
  uploadCommand = append(uploadCommand, "--fail-no-op=true")
  uploadCommand = append(uploadCommand, "--detailed-summary=true")

  execute_command "jfrog rt upload  --fail-no-op=true --detailed-summary=true"

  _, err := g.handleExecution("jf", "rt", "upload", uploadCommand...)
	if err != nil {
		return err
	}

  _, err = g.handleExecution("jf", "rt", "build-collect-env", g.runVariables[stepName+"_buildName"], g.runVariables[stepName+"_buildNumber"])
	if err != nil {
		return err
	}

	if g.inputs.forceXrayScan == "true" {
    if len(g.inputs.failOnScan) > 0 {
      g.inputs.failOnScan = "true"
    }
    scanCommand = ""
    scanCommand = append(scanCommand, "--insecure-tls=", os.Getenv(no_verify_ssl))
    scanCommand = append(scanCommand, "--fail=", g.inputs.failOnScan)
    _, err = g.handleExecution("jf", "rt", "build-scan", scanCommand, g.runVariables[stepName+"_buildName"], g.runVariables[stepName+"_buildNumber"])
    if err != nil {
      return err
    }
  }

	err = tasks.AddCacheFiles("output", []string{g.inputs.outputLocation})
	if err != nil {
		haltExecution("Failed to cache outputFileLocation " + g.inputs.outputLocation)
		return err
	}

	err = g.addStepVariables()
	if err != nil {
		return err
	}

	stepName := tasks.GetStep().Name
	

	return nil
}

// addStepVariables uses tasks go sdk to set run variables
func (m *UploadArtifact) addStepVariables() error {
	runVariables := prepareRunVariables()
	m.runVariables = runVariables
	// Add run variables present in map
	for k, v := range runVariables {
		err := tasks.AddRunVariable(k, v)
		if err != nil {
			tasks.Warn("Failed to add run variable: ", k)
		}
	}
	return nil
}

func prepareRunVariables() map[string]string {
	// Create map to store all required step variables
	runVariables := make(map[string]string)
	stepName := tasks.GetStep().Name
	if len(stepName) > 0 {
		runVariables[stepName+"_payloadType"] = "mvn"
		jfCLIBuildNum := os.Getenv("JFROG_CLI_BUILD_NUMBER")
		runVariables[stepName+"_buildNumber"] = jfCLIBuildNum
		jfCLIBuildName := os.Getenv("JFROG_CLI_BUILD_NAME")
		runVariables[stepName+"_buildName"] = jfCLIBuildName
		runVariables[stepName+"_isPromoted"] = "false"
		sourceLocation := os.Getenv("sourceLocation")
		runVariables[stepName+"_sourceLocation"] = sourceLocation
	}
	return runVariables
}

func main() {
	tasks.Info("Preparing upload-artifact task...")
	m := new(UploadArtifact)
	m.readInputs()
	m.runPreRequisites()
	err := m.run()
	if err != nil {
		haltExecution(fmt.Sprintf("%+v", err))
	}
}