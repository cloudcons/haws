package stack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/dragosboca/haws/pkg/logger"

	"github.com/tidwall/pretty"
)

// CloudFormationAPI defines the subset of methods used from the AWS CloudFormation client
// This allows us to mock AWS calls in tests.
type CloudFormationAPI interface {
	CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error)
	DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error)
	DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error)
	ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error)
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

type Stack struct {
	Template
	cloudFormationClient CloudFormationAPI
	Outputs             map[string]string
}

func NewStack(template Template) *Stack {
	return &Stack{
		Template: template,
		Outputs:  make(map[string]string),
	}
}

// Run creates or updates the stack
func (st *Stack) Run(ctx context.Context) error {
	var cfg aws.Config
	var err error

	if st.cloudFormationClient == nil {
		if st.GetRegion() != "" {
			cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(st.GetRegion()))
		} else {
			cfg, err = config.LoadDefaultConfig(ctx)
		}
		if err != nil {
			return fmt.Errorf("unable to load SDK config: %w", err)
		}
		st.cloudFormationClient = cloudformation.NewFromConfig(cfg)
	}

	templateBody, err := st.templateJson()
	if err != nil {
		return err
	}

	csName, csType, err := st.initialChangeSet(ctx, templateBody)
	if err != nil {
		return err
	}

	ok, err := st.waitForChangeSet(ctx, csName)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err := st.executeChangeSet(ctx, csName, csType); err != nil {
		return err
	}

	return nil
}

// DryRun prints the template to stdout
func (st *Stack) DryRun(ctx context.Context) error {
	templateBody, err := st.templateJson()
	if err != nil {
		return err
	}

	for k, v := range st.GetDryRunOutputs() {
		st.Outputs[k] = v
	}

	// Use formatted output with colors from pretty package
	coloredOutput := string(pretty.Color([]byte(templateBody), nil))
	logger.Info("Template JSON:\n%s", coloredOutput)

	return nil
}

// GetOutputs gets the outputs of the stack
func (st *Stack) GetOutputs(ctx context.Context) error {
	// Initialize the client if it hasn't been initialized already (e.g., in DryRun)
	if st.cloudFormationClient == nil {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(st.GetRegion()))
		if err != nil {
			return fmt.Errorf("unable to load SDK config: %w", err)
		}
		st.cloudFormationClient = cloudformation.NewFromConfig(cfg)
	}
	
	response, err := st.cloudFormationClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: st.GetStackName(),
	})
	if err != nil {
		return err
	}

	if len(response.Stacks) > 1 {
		return fmt.Errorf("multiple results for the same stack name %s", *st.GetStackName())
	}

	for _, a := range response.Stacks[0].Outputs {
		st.Outputs[*a.OutputKey] = *a.OutputValue
	}

	return nil
}

// Remove duplicate helpers: stackExist, templateJson, initialChangeSet, waitForChangeSet, executeChangeSet
// These are already defined in changeset.go and should not be redefined here.
