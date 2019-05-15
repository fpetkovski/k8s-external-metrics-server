package beanstalkd_client

import "testing"

func TestGetJobsCount(t *testing.T) {
	client := NewClient("192.168.99.100", "30002")
	jobsCount := client.GetJobsCount("default")
	if jobsCount != 0 {
		t.Error("Jobs count should be 0")
	}
}