package httputils

import "testing"

func TestUserPasswordFromAutzHdr(t *testing.T) {
	user, password, err := UserPasswordFromAutzHdr("Basic YWRtaW46eHBlcnRzY2Fu")
	if user != "admin" || password != "xpertscan" {
		t.Fatalf("Failed to decode HTTP autz header user:%s password:%s err:%s", user, password, err)
	}
}
