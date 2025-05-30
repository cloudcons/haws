package haws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/dragosboca/haws/pkg/components"
	"github.com/dragosboca/haws/pkg/logger"
	"github.com/dragosboca/haws/pkg/stack"
)

type Haws struct {
	dryRun bool
	stacks map[string]*stack.Stack
}

func New(dryRun bool, prefix string, region string, zone_id string, bucketPath string, record string) Haws {
	domain, err := getZoneDomain(zone_id)
	if err != nil {
		logger.Fatal("Failed to get zone domain: %v", err)
	}

	h := Haws{
		dryRun: dryRun,
		stacks: make(map[string]*stack.Stack),
	}

	h.stacks["certificate"] = stack.NewStack(components.NewCertificate(&components.CertificateInput{
		Prefix: prefix,
		Region: region,
		Domain: domain,
		ZoneId: zone_id,
	}))

	h.stacks["bucket"] = stack.NewStack(components.NewBucket(&components.BucketInput{
		Prefix: prefix,
		Region: region,
		Domain: domain,
	}))

	h.stacks["cloudfront"] = stack.NewStack(components.NewCdn(&components.CdnInput{
		Prefix:         prefix,
		Path:           bucketPath,
		Region:         region,
		Domain:         domain,
		Record:         record,
		CertificateArn: h.stacks["certificate"].GetExportName("Arn"), // FIXME: this is not working because one cannot reference a resource from another region
		BucketDomain:   h.stacks["bucket"].GetExportName("Domain"),
		BucketOAI:      h.stacks["bucket"].GetExportName("Oai"),
		ZoneId:         zone_id,
	}))

	h.stacks["user"] = stack.NewStack(components.NewIamUser(&components.UserInput{
		Prefix:        prefix,
		Path:          bucketPath,
		Region:        region,
		Domain:        domain,
		Record:        record,
		BucketName:    h.stacks["bucket"].GetExportName("Name"),
		CloudfrontArn: h.stacks["cloudfront"].GetExportName("Arn"),
	}))
	return h
}

func (h *Haws) Deploy(ctx context.Context) error {
	stacks := []string{"certificate", "bucket", "cloudfront", "user"}
	for _, stack := range stacks {
		if stack == "cloudfront" { // CloudFormation cross-region limitation workaround
			if err := h.GetStackOutput(ctx, "certificate"); err != nil {
				return err
			}

			certificateArn, err := h.GetOutputByName("certificate", "certificateArn")
			if err != nil {
				return err
			}

			if err = h.SetStackParameterValue("certificate", "certificateArn", certificateArn); err != nil {
				return err
			}
		}
		if err := h.DeployStack(ctx, stack); err != nil {
			return err
		}
	}
	return nil
}

func (h *Haws) DeployStack(ctx context.Context, name string) error {
	if h.dryRun {
		logger.Info("DryRunning %s", name)
		return h.stacks[name].DryRun(ctx)
	} else {
		logger.Info("Running %s", name)
		return h.stacks[name].Run(ctx)
	}
}

func (h *Haws) GetStackOutput(ctx context.Context, name string) error {
	return h.stacks[name].GetOutputs(ctx)
}

func getZoneDomain(zoneId string) (string, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %w", err)
	}
	
	svc := route53.NewFromConfig(cfg)
	
	result, err := svc.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: &zoneId,
	})
	if err != nil {
		return "", err
	}
	// trim trailing dot if any
	domain := strings.TrimSuffix(*result.HostedZone.Name, ".")

	return domain, nil
}

func (h *Haws) SetStackParameterValue(stack string, parameter string, value string) error {
	if st, ok := h.stacks[stack]; ok {
		return st.SetParameterValue(parameter, value)
	}
	return fmt.Errorf("stack %s not found", stack)
}

func (h *Haws) GetOutputByName(stack string, output string) (string, error) {
	if st, ok := h.stacks[stack]; ok {
		return st.Outputs[st.GetExportName(output)], nil
	}
	return "", fmt.Errorf("stack %s not found", stack)

}
