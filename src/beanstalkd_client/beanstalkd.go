package beanstalkd_client

import (
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
)

type beanstalkdClient struct {
	connection *beanstalk.Conn
}

func NewClient(host string, port string) *beanstalkdClient {
	c, err := beanstalk.Dial("tcp", host + ":" + port)
	if err != nil {
		panic(err)
	}

	return &beanstalkdClient{
		connection:c,
	}
}

func getValue(stats map[string]string, item string) int64 {
	value, err := strconv.ParseInt(stats[item], 10, 64)
	if err != nil {
		panic(err)
	}

	return value
}

func (client *beanstalkdClient) GetJobsCount(tubeName string) int64 {
	tube := beanstalk.Tube{
		Conn: client.connection,
		Name: tubeName,
	}

	stats, err := tube.Stats()
	if err != nil {
		panic(err)
	}

	readyJobs := getValue(stats, "current-jobs-ready")
	reservedJobs := getValue(stats, "current-jobs-reserved")
	return readyJobs + reservedJobs
}