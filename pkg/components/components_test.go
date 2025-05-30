package components

import "testing"

func TestNewBucketAndExports(t *testing.T) {
	b := NewBucket(&BucketInput{Prefix: "test", Region: "us-east-1", Domain: "example.com"})
	if b == nil {
		t.Fatal("NewBucket returned nil")
	}
	if b.GetExportName("Arn") == "" {
		t.Error("GetExportName should not return empty string")
	}
	if b.GetStackName() == nil || *b.GetStackName() == "" {
		t.Error("GetStackName should not return nil or empty string")
	}
}

func TestNewCertificateAndExports(t *testing.T) {
	c := NewCertificate(&CertificateInput{Prefix: "test", Region: "us-east-1", Domain: "example.com", ZoneId: "zone"})
	if c == nil {
		t.Fatal("NewCertificate returned nil")
	}
	if c.GetExportName("Arn") == "" {
		t.Error("GetExportName should not return empty string")
	}
	if c.GetStackName() == nil || *c.GetStackName() == "" {
		t.Error("GetStackName should not return nil or empty string")
	}
}

func TestNewCdnAndExports(t *testing.T) {
	cdn := NewCdn(&CdnInput{Prefix: "test", Path: "/", Region: "us-east-1", Domain: "example.com", Record: "www"})
	if cdn == nil {
		t.Fatal("NewCdn returned nil")
	}
	if cdn.GetExportName("Arn") == "" {
		t.Error("GetExportName should not return empty string")
	}
	if cdn.GetStackName() == nil || *cdn.GetStackName() == "" {
		t.Error("GetStackName should not return nil or empty string")
	}
}

func TestNewIamUserAndExports(t *testing.T) {
	u := NewIamUser(&UserInput{Prefix: "test", Path: "/", Region: "us-east-1", Domain: "example.com", Record: "www", BucketName: "bucket", CloudfrontArn: "arn"})
	if u == nil {
		t.Fatal("NewIamUser returned nil")
	}
	if u.GetExportName("Arn") == "" {
		t.Error("GetExportName should not return empty string")
	}
	if u.GetStackName() == nil || *u.GetStackName() == "" {
		t.Error("GetStackName should not return nil or empty string")
	}
}

func TestBucketExportNameVariants(t *testing.T) {
	b := NewBucket(&BucketInput{Prefix: "x", Region: "us-east-1", Domain: "d.com"})
	if b.GetExportName("") == "" {
		t.Error("GetExportName with empty string should not return empty string")
	}
	if b.GetExportName("Oai") == "" {
		t.Error("GetExportName with Oai should not return empty string")
	}
}

func TestCertificateExportNameVariants(t *testing.T) {
	c := NewCertificate(&CertificateInput{Prefix: "x", Region: "us-east-1", Domain: "d.com", ZoneId: "z"})
	if c.GetExportName("") == "" {
		t.Error("GetExportName with empty string should not return empty string")
	}
}

func TestCdnExportNameVariants(t *testing.T) {
	cdn := NewCdn(&CdnInput{Prefix: "x", Path: "/", Region: "us-east-1", Domain: "d.com", Record: "r"})
	if cdn.GetExportName("") == "" {
		t.Error("GetExportName with empty string should not return empty string")
	}
}

func TestUserExportNameVariants(t *testing.T) {
	u := NewIamUser(&UserInput{Prefix: "x", Path: "/", Region: "us-east-1", Domain: "d.com", Record: "r", BucketName: "b", CloudfrontArn: "a"})
	if u.GetExportName("") == "" {
		t.Error("GetExportName with empty string should not return empty string")
	}
}
