# upload-artifact

This Pipelines Task is implementation of uploadArtifact custom step.

As per input configuration this task performs below actions

- This task can be used to upload files to Artifactory.
- FileSpec, RemoteFile, and GitRepo inputs are supported (up to one of each), and by default everything in those resources will be uploaded.
- Optional outputs are BuildInfo (with autoPublishBuildInfo set to true) and FileSpec resources. See the JFrog Artifactory CLI rt upload command for more information about the possible filters.


### What's New
- initial release

### Prerequisites
This task requires an artifactory integration where you want to upload the files

## Usage

**Basic:**

```yaml
- task: jfrog/upload-artifact@v0.0.1
  id: upload-artifact
  input:
    artifactoryIntegration: artifactory-integration
    inputResource: input-resource-name
    buildInfoResource: output-build-info-resource-name
    sourcePath: *
    targetPath: some-repository-on-artifactory
    autoPublishBuildInfo: true
    forceXrayScan: false
    failOnScan: false
```

### Input Variables

| Name                        | Required | Default                               | Description                     |
|-----------------------------|----------|---------------------------------------|---------------------------------|
| inputResource                       | true      |     | May specify a GitRepo, FileSpec, or RemoteFile resource containing the file(s) to be uploaded |
| autoPublishBuildInfo                | false     |     | If true, Build Info for the step will be published. |
| buildInfoResource                   | false     |                                 | Must specify a BuildInfo resource if autoPublishBuildInfo is set as true |
| targetPath                          | true      |     | Where to upload the files, including repository name. Required |
| sourcePath                          | false      |     | Files to upload. Default * |
| properties                          | false     |     | Semi-colon separated properties for the uploaded Artifact, e.g. "myFirstProperty=one;mySecondProperty=two". pipelines_step_name, pipelines_run_number, pipelines_step_id, pipelines_pipeline_name, pipelines_step_type, and pipelines_step_platform will also be added |
| regExp                              | false     |     | If true, sourcePath uses regular expressions instead of wildcards. |
| flat                                | false     |     | If true, the uploaded files are flattened removing the directory structure. |
| module                              | false     |     | A module name for the Build Info. |
| deb                                 | false     |     | A distribution/component/architecture for Debian packages. If a component includes a / it must be double-escaped, e.g. distribution/my\\\/component/architecture for a my/component component. |
| recursive                           | false     |     | If false, do not upload any matches in sub-directories. |
| dryRun                              | false     |     | If true, nothing is uploaded. |
| symlinks                            | false     |     | If true, symlinks are uploaded. |
| explode                             | false     |     | If true and the uploaded Artifact is an archive, the archive is expanded. |
| exclusions                          | false     |     | Semi-colon separated patterns to exclude. |
| includeDirs                         | false     |     | If true, empty directories matching the criteria are uploaded. |
| syncDeletes                         | false     |     | A path under which to delete any existing files in Artifactory. |
| forceXrayScan                       | false     |     | If true, an Xray Scan will be triggered for the step. |
| failOnScan                          | false     |     | If a scan failure should cause a step failure, default true. |
