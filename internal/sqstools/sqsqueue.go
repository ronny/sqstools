package sqstools

import "regexp"

var (
	pattern = regexp.MustCompile(`^https://sqs.(?P<region>[a-z\-0-9]+).amazonaws.com/(?P<accountID>\d+)/(?P<queueName>.*)$`)
)

type SQSQueue struct {
	url     string
	matches map[string]string
}

func NewSQSQueue(url string) SQSQueue {
	return SQSQueue{
		url:     url,
		matches: findStringSubmatchMap(url, pattern),
	}
}

func (u *SQSQueue) Region() string {
	return u.matches["region"]
}

func (u *SQSQueue) AccountID() string {
	return u.matches["accountID"]
}

func (u *SQSQueue) QueueName() string {
	return u.matches["queueName"]
}

func (u *SQSQueue) URL() string {
	return u.url
}

func findStringSubmatchMap(s string, r *regexp.Regexp) map[string]string {
	groups := make(map[string]string)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return groups
	}

	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		groups[name] = match[i]
	}
	return groups
}
