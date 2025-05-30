package haws

import (
	"fmt"
	"strings"
	"text/template"
	
	"github.com/dragosboca/haws/pkg/logger"
)

const hugoConfig = `
[[deployment.targets]]
    name = "Haws"
    URL = "s3://{{ .BucketName }}?region={{ .Region }}&prefix={{ .OriginPath }}/"
    cloudFrontDistributionID = "{{ .CloudFrontId }}"
`

type deployment struct {
	BucketName   string
	Region       string
	CloudFrontId string
	OriginPath   string
}

func (h *Haws) GenerateHugoConfig(region string, path string) {
	bucketName, err := h.GetOutputByName("bucket", "Name")
	if err != nil {
		logger.Fatal("Failed to get bucket name: %v", err)
	}

	cloudFrontId, err := h.GetOutputByName("cloudfront", "CloudFrontId")
	if err != nil {
		logger.Fatal("Failed to get CloudFront ID: %v", err)
	}

	deploymentConfig := deployment{
		BucketName:   bucketName,
		Region:       region,
		CloudFrontId: cloudFrontId,
		OriginPath:   fmt.Sprintf("%s/", strings.Trim(path, "/")),
	}
	t := template.Must(template.New("deployment").Parse(hugoConfig))
	
	logger.Info("\n\nUse the following minimal configuration for hugo deployment")
	
	// Create a buffer to collect template output
	var outputBuffer strings.Builder
	err = t.Execute(&outputBuffer, deploymentConfig)
	if err != nil {
		logger.Fatal("Error executing template: %s", err)
	}
	
	// Print the output using our logger
	logger.Info("\n%s", outputBuffer.String())
}
