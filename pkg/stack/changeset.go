package stack

import (
	"fmt"
	"strings"
	"time"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/dragosboca/haws/pkg/logger"
	"github.com/goombaio/namegenerator"
)

const EmptyChangeSet = "The submitted information didn't contain changes. Submit different information to create a change set."

func (st *Stack) stackExist(ctx context.Context) bool {
	_, err := st.cloudFormationClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: st.GetStackName(),
	})
	return err == nil
}

// templateJson returns the template as a JSON string
// return: string - the template as a JSON string
// return: error - the error if any
func (st *Stack) templateJson() (string, error) {
	template := st.Build()
	templateBody, err := template.JSON()
	if err != nil {
		logger.Error("Create template error: %s", err)
		return "", err
	}
	return string(templateBody), nil
}

// initialChangeSet creates the initial changeset
// param: templateBody - the template body
// return: csName - the name of the changeset
// return: csType - the type of the changeset (CREATE or UPDATE)
// return: error - the error if any
func (st *Stack) initialChangeSet(ctx context.Context, templateBody string) (string, string, error) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	csName := nameGenerator.Generate()

	csType := "CREATE"
	if st.stackExist(ctx) {
		csType = "UPDATE"
		logger.Info("Updating stack: %s with changeset: %s", *st.GetStackName(), csName)
	} else {
		logger.Info("Creating stack: %s with changeset: %s", *st.GetStackName(), csName)
	}

	// Convert CloudFormation v1 Parameters to v2 Parameters
	v2Params := make([]types.Parameter, 0, len(st.GetParameters()))
	for _, param := range st.GetParameters() {
		v2Params = append(v2Params, types.Parameter{
			ParameterKey:   param.ParameterKey,
			ParameterValue: param.ParameterValue,
		})
	}
	
	_, err := st.cloudFormationClient.CreateChangeSet(ctx, &cloudformation.CreateChangeSetInput{
		ClientToken:   &csName,
		ChangeSetName: &csName,
		ChangeSetType: types.ChangeSetType(csType),
		Parameters:    v2Params,
		StackName:     st.GetStackName(),
		TemplateBody:  &templateBody,
	})
	if err != nil {
		return "", "", err
	}
	return csName, csType, nil
}

// waitForChangeSet returns true if the changeset is empty and should be deleted
// param: csName - the name of the changeset
// return: bool - true if the changeset is empty and should be deleted
// return: error - the error if any
func (st *Stack) waitForChangeSet(ctx context.Context, csName string) (bool, error) {
	logger.Info("Waiting for the changeset %s creation to complete", csName)
	
	// AWS SDK v2 doesn't have built-in waiters, so we'll implement a simple polling mechanism
	maxAttempts := 30
	delay := time.Second * 2
	
	for i := 0; i < maxAttempts; i++ {
		desc, err := st.cloudFormationClient.DescribeChangeSet(ctx, &cloudformation.DescribeChangeSetInput{
			ChangeSetName: &csName,
			StackName:     st.GetStackName(),
		})
		
		if err != nil {
			return false, err
		}
		
		// Check if the change set is ready
		if desc.Status == types.ChangeSetStatusCreateComplete {
			return false, nil
		}
		
		// Check if the change set failed because it's empty
		if desc.Status == types.ChangeSetStatusFailed && *desc.StatusReason == EmptyChangeSet {
			logger.Info("Deleting empty changeset %s", csName)
			_, err := st.cloudFormationClient.DeleteChangeSet(ctx, &cloudformation.DeleteChangeSetInput{
				ChangeSetName: &csName,
				StackName:     st.GetStackName(),
			})
			if err != nil {
				return false, err
			}
			return true, nil
		} else if desc.Status == types.ChangeSetStatusFailed {
			// Failed for some other reason
			return false, fmt.Errorf("change set creation failed: %s", *desc.StatusReason)
		}
		
		// Wait before checking again
		time.Sleep(delay)
	}
	
	return false, fmt.Errorf("timed out waiting for change set creation to complete")
}

// executeChangeSet executes the changeset
// param: csName - the name of the changeset
// param: csType - the type of the changeset (CREATE or UPDATE)
// return: error - the error if any
func (st *Stack) executeChangeSet(ctx context.Context, csName string, csType string) error {
	logger.Info("Executing change set: %s on stack %s", csName, *st.GetStackName())
	_, err := st.cloudFormationClient.ExecuteChangeSet(ctx, &cloudformation.ExecuteChangeSetInput{
		ChangeSetName:      &csName,
		ClientRequestToken: &csName,
		StackName:          st.GetStackName(),
	})
	if err != nil {
		return err
	}

	logger.Info("Waiting for the changeset %s execution to complete", csName)
	
	// Implementation of waiter using polling since SDK v2 doesn't have built-in waiters
	maxAttempts := 120 // Cloud formation stacks can take a while
	delay := time.Second * 5
	
	targetStatus := ""
	if csType == "CREATE" {
		targetStatus = string(types.StackStatusCreateComplete)
	} else {
		targetStatus = string(types.StackStatusUpdateComplete)
	}
	
	for i := 0; i < maxAttempts; i++ {
		resp, err := st.cloudFormationClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: st.GetStackName(),
		})
		
		if err != nil {
			return err
		}
		
		if len(resp.Stacks) == 0 {
			return fmt.Errorf("stack not found")
		}
		
		stackStatus := string(resp.Stacks[0].StackStatus)
		
		if stackStatus == targetStatus {
			return nil
		}
		
		// Check for failure states
		if strings.HasSuffix(stackStatus, "FAILED") || 
		   strings.HasSuffix(stackStatus, "ROLLBACK_COMPLETE") {
			return fmt.Errorf("stack operation failed: %s", stackStatus)
		}
		
		time.Sleep(delay)
	}
	
	return fmt.Errorf("timed out waiting for stack operation to complete")
}
