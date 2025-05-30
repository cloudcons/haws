package stack

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cfn "github.com/awslabs/goformation/v4/cloudformation"
	s3 "github.com/awslabs/goformation/v4/cloudformation/s3"
)

// MockTemplate is a simple implementation of the Template interface for testing
type MockTemplate struct {
	stackName string
	region    string
}

func (m *MockTemplate) GetStackName() *string {
	return aws.String(m.stackName)
}

func (m *MockTemplate) GetRegion() string {
	return m.region
}

func (m *MockTemplate) GetExportName(name string) string {
	return m.stackName + "-" + name
}

func (m *MockTemplate) Build() *cfn.Template {
	return cfn.NewTemplate()
}

func (m *MockTemplate) GetParameters() []types.Parameter {
	return []types.Parameter{}
}

func (m *MockTemplate) GetDryRunOutputs() map[string]string {
	return map[string]string{
		m.GetExportName("Arn"): "arn:aws:acm:region:account:certificate/12345",
	}
}

func (m *MockTemplate) SetParameterValue(name, value string) error {
	return nil
}

func TestNewStack(t *testing.T) {
	// Create a mock template
	mockTemplate := &MockTemplate{
		stackName: "test-prefix-certificate",
		region: "us-east-1",
	}
	
	// Create a new stack using the mock template
	stack := NewStack(mockTemplate)
	
	// Test that the stack is properly initialized
	if stack == nil {
		t.Fatal("Stack should not be nil")
	}
	
	if stack.Template == nil {
		t.Fatal("Stack.Template should not be nil")
	}
	
	if stack.ctx == nil {
		t.Fatal("Stack.ctx should not be nil")
	}
	
	if stack.Outputs == nil {
		t.Fatal("Stack.Outputs should not be nil")
	}
	
	// Test that template values are properly set
	if stack.GetRegion() != "us-east-1" {
		t.Errorf("Expected region to be 'us-east-1', got '%s'", stack.GetRegion())
	}
	
	expectedStackName := "test-prefix-certificate"
	if *stack.GetStackName() != expectedStackName {
		t.Errorf("Expected stack name to be '%s', got '%s'", expectedStackName, *stack.GetStackName())
	}
}

func TestGetExportName(t *testing.T) {
	// Create a mock template
	mockTemplate := &MockTemplate{
		stackName: "test-prefix-certificate",
		region: "us-east-1",
	}
	
	// Create a new stack using the mock template
	stack := NewStack(mockTemplate)
	
	// Test GetExportName
	exportName := stack.GetExportName("Arn")
	expectedExport := "test-prefix-certificate-Arn"
	
	if exportName != expectedExport {
		t.Errorf("Expected export name to be '%s', got '%s'", expectedExport, exportName)
	}
}

func TestSetParameterValue(t *testing.T) {
	mock := &MockTemplate{stackName: "foo", region: "us-east-1"}
	err := mock.SetParameterValue("param", "value")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetParameters(t *testing.T) {
	mock := &MockTemplate{stackName: "foo", region: "us-east-1"}
	params := mock.GetParameters()
	if params == nil {
		t.Error("Expected parameters slice, got nil")
	}
}

func TestGetDryRunOutputs(t *testing.T) {
	mock := &MockTemplate{stackName: "foo", region: "us-east-1"}
	outputs := mock.GetDryRunOutputs()
	if outputs == nil || outputs["foo-Arn"] == "" {
		t.Error("Expected dry run outputs to contain foo-Arn")
	}
}

func TestBuild(t *testing.T) {
	mock := &MockTemplate{stackName: "foo", region: "us-east-1"}
	tmpl := mock.Build()
	if tmpl == nil {
		t.Error("Expected non-nil template")
	}
}

func TestStack_Run_DryRun_GetOutputs(t *testing.T) {
	t.Skip("Skipping test that requires real AWS calls")
}

func TestStack_Run_Error(t *testing.T) {
	t.Skip("Skipping test that requires real AWS calls")
}

func TestStack_DryRun_Error(t *testing.T) {
	t.Skip("Skipping test that requires real AWS calls")
}

func TestStack_GetOutputs_Error(t *testing.T) {
	t.Skip("Skipping test that requires real AWS calls")
}

func TestTemplate_AddParameter_AddResource_AddOutput(t *testing.T) {
	tmpl := NewTemplate("us-east-1")
	param := cfn.Parameter{Description: "desc", Default: "default"}
	tmpl.AddParameter("Param1", param, "default")
	bucket := &s3.Bucket{}
	tmpl.AddResource("Res1", bucket)
	output := cfn.Output{Description: "desc", Value: "value"}
	tmpl.AddOutput("Out1", output, "dryrunval")
	params := tmpl.GetParameters()
	if len(params) == 0 {
		t.Error("Expected parameters to be added")
	}
	region := tmpl.GetRegion()
	if region != "us-east-1" {
		t.Errorf("Expected region 'us-east-1', got '%s'", region)
	}
	err := tmpl.SetParameterValue("Param1", "newval")
	if err != nil {
		t.Errorf("Expected no error setting parameter value, got %v", err)
	}
}

// Mock for successful DescribeStacks
// Moved outside the test function for Go syntax

type mockCFNSuccess struct{}
func (m *mockCFNSuccess) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return &cloudformation.DescribeStacksOutput{
		Stacks: []types.Stack{
			{
				Outputs: []types.Output{
					{OutputKey: aws.String("mock-key"), OutputValue: aws.String("mock-value")},
				},
			},
		},
	}, nil
}

// Satisfy CloudFormationAPI interface for test mocks
func (m *mockCFNSuccess) CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNSuccess) DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNSuccess) DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNSuccess) ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error) {
	return nil, nil
}

// Mock for error DescribeStacks

type mockCFNError struct{}
func (m *mockCFNError) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return nil, fmt.Errorf("mock error")
}

// Satisfy CloudFormationAPI interface for test mocks
func (m *mockCFNError) CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNError) DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNError) DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error) {
	return nil, nil
}
func (m *mockCFNError) ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error) {
	return nil, nil
}

func TestStack_GetOutputs_Mock(t *testing.T) {
	mockTemplate := &MockTemplate{stackName: "mock-stack", region: "us-east-1"}
	stk := NewStack(mockTemplate)
	stk.cloudFormationClient = &mockCFNSuccess{}

	err := stk.GetOutputs(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if v, ok := stk.Outputs["mock-key"]; !ok || v != "mock-value" {
		t.Errorf("Expected output 'mock-key' to be 'mock-value', got '%v'", stk.Outputs)
	}
}

func TestStack_GetOutputs_MockError(t *testing.T) {
	mockTemplate := &MockTemplate{stackName: "mock-stack", region: "us-east-1"}
	stk := NewStack(mockTemplate)
	stk.cloudFormationClient = &mockCFNError{}

	err := stk.GetOutputs(context.Background())
	if err == nil || err.Error() != "mock error" {
		t.Errorf("Expected 'mock error', got %v", err)
	}
}

// Mock for Stack.Run and Stack.DryRun
// Only DescribeStacks is used in GetOutputs, but for Run/DryRun, you may want to simulate CreateChangeSet, ExecuteChangeSet, etc.
type mockCFNRun struct {
	createChangeSetErr   error
	describeChangeSetErr error
	describeChangeSetCnt int
	executeChangeSetErr  error
	describeStacksErr    error
	describeStacksCnt    int
}

func (m *mockCFNRun) CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error) {
	return &cloudformation.CreateChangeSetOutput{}, m.createChangeSetErr
}
func (m *mockCFNRun) DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error) {
	m.describeChangeSetCnt++
	if m.describeChangeSetCnt < 2 {
		return &cloudformation.DescribeChangeSetOutput{
			Status: types.ChangeSetStatusCreateInProgress,
		}, m.describeChangeSetErr
	}
	return &cloudformation.DescribeChangeSetOutput{
		Status: types.ChangeSetStatusCreateComplete,
	}, m.describeChangeSetErr
}
func (m *mockCFNRun) DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error) {
	return &cloudformation.DeleteChangeSetOutput{}, nil
}
func (m *mockCFNRun) ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error) {
	return &cloudformation.ExecuteChangeSetOutput{}, m.executeChangeSetErr
}
func (m *mockCFNRun) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	m.describeStacksCnt++
	if m.describeStacksCnt < 2 {
		return &cloudformation.DescribeStacksOutput{
			Stacks: []types.Stack{{StackStatus: types.StackStatusCreateInProgress}},
		}, m.describeStacksErr
	}
	return &cloudformation.DescribeStacksOutput{
		Stacks: []types.Stack{{StackStatus: types.StackStatusCreateComplete}},
	}, m.describeStacksErr
}

// Test for Stack.Run with all mocks succeeding
func TestStack_Run_Mock(t *testing.T) {
	mockTemplate := &MockTemplate{stackName: "mock-stack", region: "us-east-1"}
	stk := NewStack(mockTemplate)
	stk.cloudFormationClient = &mockCFNRun{}
	err := stk.Run(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Test for Stack.Run with CreateChangeSet error
func TestStack_Run_Mock_CreateChangeSetError(t *testing.T) {
	mockTemplate := &MockTemplate{stackName: "mock-stack", region: "us-east-1"}
	stk := NewStack(mockTemplate)
	stk.cloudFormationClient = &mockCFNRun{createChangeSetErr: fmt.Errorf("create changeset error")}
	err := stk.Run(context.Background())
	if err == nil || err.Error() != "create changeset error" {
		t.Errorf("Expected 'create changeset error', got %v", err)
	}
}

// Test for Stack.Run with ExecuteChangeSet error
func TestStack_Run_Mock_ExecuteChangeSetError(t *testing.T) {
	mockTemplate := &MockTemplate{stackName: "mock-stack", region: "us-east-1"}
	stk := NewStack(mockTemplate)
	stk.cloudFormationClient = &mockCFNRun{executeChangeSetErr: fmt.Errorf("execute changeset error")}
	err := stk.Run(context.Background())
	if err == nil || err.Error() != "execute changeset error" {
		t.Errorf("Expected 'execute changeset error', got %v", err)
	}
}
