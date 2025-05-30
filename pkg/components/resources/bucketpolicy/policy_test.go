package bucketpolicy

import "testing"

func TestNewAndAddStatement(t *testing.T) {
	d := New("id1")
	if d == nil {
		t.Fatal("New returned nil")
	}
	if d.Version == "" {
		t.Error("Version should not be empty")
	}
	if d.Id != "id1" {
		t.Errorf("Expected Id 'id1', got '%s'", d.Id)
	}
	if len(d.Statement) != 0 {
		t.Errorf("Expected no statements initially, got %d", len(d.Statement))
	}

	s := Statement{
		Effect: "Allow",
		Principal: Principal{"AWS": "arn:aws:iam::123456789012:user/test"},
		Action: []string{"s3:GetObject"},
		Resource: []string{"arn:aws:s3:::bucket/*"},
	}
	d.AddStatement("sid1", s)
	if len(d.Statement) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(d.Statement))
	}
	if d.Statement[0].Sid != "sid1" {
		t.Errorf("Expected Sid 'sid1', got '%s'", d.Statement[0].Sid)
	}
}
