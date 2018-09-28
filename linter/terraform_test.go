package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"os"
	"testing"
)

func TestTerraformLinter(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/terraform_instance.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformLinter to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("TestTerraformLinter scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 1 {
		t.Errorf("TestTerraformLinter scanned %d files, expecting 1", len(report.FilesScanned))
	}
	assertViolationsCount("TestTerraformLinter ", 0, report.Violations, t)
}

func TestTerraformVariables(t *testing.T) {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/uses_variables.tf")
	if err != nil {
		t.Error("Expecting TestTerraformLinter.Load to not return an error")
	}
	resources, err := loader.PostLoad(loaded)
	if err != nil {
		t.Error("Expecting TestTerraformLinter.PostLoad to not return an error")
	}
	if len(resources) != 1 {
		t.Errorf("Expecting to load 1 resources, not %d", len(loaded.Resources))
	}
	properties := resources[0].Properties.(map[string]interface{})
	if properties["ami"] != "ami-f2d3638a" {
		t.Errorf("Unexpected value for variable: %s", properties["ami"])
	}
	// this test covers string, map, and slice cases
	tags := properties["tags"].([]interface{})
	tag := tags[0].(map[string]interface{})
	project := tag["project"].(string)
	if project != "demo" {
		t.Errorf("Expected project tag to be 'demo', got: %s", project)
	}
	comment := tag["comment"].(string)
	if comment != "bar" {
		t.Errorf("Expected project tag to be 'bar', got: %s", comment)
	}
	environment := tag["environment"].(string)
	if environment != "test" {
		t.Errorf("Expected environment tag to be 'test', got: '%s'", environment)
	}
}

func TestTerraformVariablesFromEnvironment(t *testing.T) {
	os.Setenv("TF_VAR_instance_type", "c4.large")
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/uses_variables.tf")
	if err != nil {
		t.Error("Expecting TestTerraformLinter.Load to not return an error")
	}
	resources, err := loader.PostLoad(loaded)
	if err != nil {
		t.Error("Expecting TestTerraformLinter.PostLoad to not return an error")
	}
	if len(resources) != 1 {
		t.Errorf("Expecting to load 1 resources, not %d", len(loaded.Resources))
	}
	properties := resources[0].Properties.(map[string]interface{})
	if properties["instance_type"] != "c4.large" {
		t.Errorf("Unexpected value for variable: %s", properties["instance_type"])
	}
	os.Setenv("TF_VAR_instance_type", "")
}

func TestTerraformFileFunction(t *testing.T) {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/reference_file.tf")
	if err != nil {
		t.Error("Expecting TestTerraformFileFunction.Load to not return an error")
	}
	resources, err := loader.PostLoad(loaded)
	if err != nil {
		t.Error("Expecting TestTerraformFileFunction.PostLoad to not return an error")
	}
	if len(resources) != 1 {
		t.Errorf("Expecting to load 1 resources, not %d", len(loaded.Resources))
	}
	properties := resources[0].Properties.(map[string]interface{})
	if properties["bucket"] != "example" {
		t.Errorf("Unexpected value for file: %s", properties["bucket"])
	}
}

func TestTerraformVariablesInDifferentFile(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{
		"./testdata/resources/defines_variables.tf",
		"./testdata/resources/reference_variables.tf",
	}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformVariablesInDifferentFile to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("TestTerraformVariablesInDifferentFile scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 2 {
		t.Errorf("TestTerraformVariablesInDifferentFile scanned %d files, expecting 2", len(report.FilesScanned))
	}
	assertViolationsCount("TestTerraformVariablesInDifferentFile ", 0, report.Violations, t)
}

type TestingValueSource struct{}

func (s TestingValueSource) GetValue(a assertion.Expression) (string, error) {
	if a.ValueFrom.URL != "" {
		return "TEST", nil
	}
	return a.Value, nil
}

func TestTerraformDataLoader(t *testing.T) {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/terraform_data.tf")
	if err != nil {
		t.Error("Expecting TestTerraformDataLoader to not return an error")
	}
	if len(loaded.Resources) != 1 {
		t.Errorf("TestTerraformDataLoader scanned %d resources, expecting 1", len(loaded.Resources))
	}
}

type terraformLinterTestCase struct {
	ConfigurationFilename   string
	RulesFilename           string
	ExpectedViolationCount  int
	ExpectedViolationRuleID string
}

func TestTerraformLinterCases(t *testing.T) {
	testCases := map[string]terraformLinterTestCase{
		"ParseError": {
			"./testdata/resources/terraform_syntax_error.tf",
			"./testdata/rules/terraform_provider.yml",
			1,
			"FILE_LOAD",
		},
		"Provider": {
			"./testdata/resources/terraform_provider.tf",
			"./testdata/rules/terraform_provider.yml",
			1,
			"AWS_PROVIDER",
		},
		"DataObject": {
			"./testdata/resources/terraform_data.tf",
			"./testdata/rules/terraform_data.yml",
			1,
			"DATA_NOT_CONTAINS",
		},
		"PoliciesWithVariables": {
			"./testdata/resources/policy_with_variables.tf",
			"./testdata/rules/policy_variable.yml",
			0,
			"",
		},
		"HereDocWithExpression": {
			"./testdata/resources/policy_with_expression.tf",
			"./testdata/rules/policy_variable.yml",
			0,
			"",
		},
		"Policies": {
			"./testdata/resources/terraform_policy.tf",
			"./testdata/rules/terraform_policy.yml",
			1,
			"TEST_POLICY",
		},
		"PolicyInvalidJSON": {
			"./testdata/resources/terraform_policy_invalid_json.tf",
			"./testdata/rules/terraform_policy.yml",
			0,
			"",
		},
		"PolicyEmpty": {
			"./testdata/resources/terraform_policy_empty.tf",
			"./testdata/rules/terraform_policy.yml",
			0,
			"",
		},
		"Module": {
			"./testdata/resources/terraform_module.tf",
			"./testdata/rules/terraform_module.yml",
			1,
			"MODULE_DESCRIPTION",
		},
	}
	for name, tc := range testCases {
		options := Options{
			Tags:    []string{},
			RuleIDs: []string{},
		}
		filenames := []string{tc.ConfigurationFilename}
		linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
		ruleSet := loadRulesForTest(tc.RulesFilename, t)
		report, err := linter.Validate(ruleSet, options)
		if err != nil {
			t.Errorf("Expecting %s to return without an error: %s", name, err.Error())
		}
		if len(report.FilesScanned) != 1 {
			t.Errorf("TestTerraformLinterCases scanned %d files, expecting 1", len(report.FilesScanned))
		}
		if len(report.Violations) != tc.ExpectedViolationCount {
			t.Errorf("%s returned %d violations, expecting %d", name, len(report.Violations), tc.ExpectedViolationCount)
			t.Errorf("Violations: %v", report.Violations)
		}
		if tc.ExpectedViolationRuleID != "" {
			assertViolationByRuleID(name, tc.ExpectedViolationRuleID, report.Violations, t)
		}
	}
}
