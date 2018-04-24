[![Build Status](https://circleci.com/gh/stelligent/config-lint.svg?style=shield)](https://circleci.com/gh/stelligent/config-lint)

# config-lint

A command line tool to validate configurations using rules specified in a YAML file.
The data being validated can come from template files, such as a Terraform file.
There is also an example of a Linter that runs agains data returned from an AWS API call.

# Installation 
You can use [Homebrew](https://brew.sh/) to install the latest version:

```
brew tap stelligent/tap
brew install config-lint
```

Alternatively, you can install manually from the [releases](https://github.com/stelligent/config-lint/releases).

# Build Command Line tool

```
make all
```

# Run

The program has a set of built-in rules for scanning the following types of files:

* Terraform

The program can also read files from a separate YAML file, and can scan these types of files:

* Terraform
* Kubernetes
* LintRules
* YAML

And also the scanning of information from AWS Descibe API calls for:

* Security Groups
* IAM Users


## Example invocations

### Validate Terraform files with built-in rules

```
config-lint -terraform example-files/config/*
```

### Validate Terraform files with custom rules

```
config-lint -rules examples-files/rules/terraform.yml example-files/config/*
```

### Validate Kubernetes files

```
config-lint -rules example-files/rules/kubernetes.yml example-files/config/*
```

### Validate LintRules files

This type of linting allows the tool to lint its own rules.

```
config-lint -rules example-files/rules/lint-rules.yml example-files/rules/*
```

### Validate Existing Security Groups

```
config-lint -rules example-files/rules/security-groups.yml
```

### Validate Existing IAM Users

```
config-lint -rules example-files/rules/iam-users.yml
```

# Rules

A YAML file that specifies what kinds of files to process, and what validations to perform, [documented here](docs/rules.md).

# Operations

The rules contain a list of expressions that use operations that are [documented here](docs/operations.md).

## Examples

See [here](docs/example-rules.md) for examples of custom rules.

# Output

The program outputs a JSON string with the results. The JSON object has the following attributes:

* FilesScanned - a list of the filenames evaluated
* Violations - an object whose keys are the severity of any violations detected. The value for each key is an array with an entry for every violation of that severity.

## Using --query to limit the output

You can limit the output by specifying a JMESPath expression for the --query command line option. For example, if you just wanted to see the ResourceId attribute for failed checks, you can do the following:

```
./config-lint --rules example-files/rules/terraform.yml --query 'Violations.FAILURE[].ResourceId' example-files/config/*
```

# Exit Code

If at least one rule with a severity of FAILURE was triggered the exit code will be 1, otherwise it will be 0.


# Profiles

The -profile command line option takes a filename which contains a set of default values for various command line options.
If there is a file in the working directory called `config-lint.yml`, it will be loaded automatically.
All values in the profile are optional, and are overriden by anything specified on the command line.
An example profile:

```
rules:
  - example-files/rules/generic.yml

files:
  - example-files/config/*.config

ids:
  - RULE_1
  - RULE_2

tags:
  - s3
```

# Developing new rules using --search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a --search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:

```
./config-lint --rules example-files/rules/terraform.yml --search 'ami' example-files/config/terraform.tf
```

If you specify --search, the rules files is only used to determine the type of configuration files.
The files will *not* be scanned for violations.

# Releasing
To release a new version, run `make bumpversion` to increment the patch version and push a tag to GitHub to start the release process.

