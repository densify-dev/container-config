package config

import (
	"fmt"
	rhttp "github.com/hashicorp/go-retryablehttp"
	"net/http"
	"strings"
	"time"
)

// policies
const (
	DefaultPolicy     = "default"
	ExponentialPolicy = "exponential"
	JitterPolicy      = "jitter"
)

var policies = map[string]rhttp.Backoff{
	Empty:             rhttp.DefaultBackoff,
	DefaultPolicy:     rhttp.DefaultBackoff,
	ExponentialPolicy: rhttp.DefaultBackoff,
	JitterPolicy:      rhttp.LinearJitterBackoff,
}

type RetryConfig struct {
	WaitMin     time.Duration `yaml:"wait_min"`
	WaitMax     time.Duration `yaml:"wait_max"`
	MaxAttempts int           `yaml:"max_attempts"`
	Policy      string        `yaml:"policy,omitempty"`
	backoff     rhttp.Backoff `yaml:"-"`
	isValid     bool          `yaml:"-"`
}

// Validate must be called once, after rc has been constructed / unmarshalled
func (rc *RetryConfig) Validate() (err error) {
	if rc != nil {
		if err = validDuration(rc.WaitMin); err == nil {
			if err = validDuration(rc.WaitMin); err == nil {
				if err = validPositive(rc.MaxAttempts); err == nil {
					if rc.backoff = policies[strings.ToLower(rc.Policy)]; rc.backoff == nil {
						err = fmt.Errorf("invalid backoff policy %s", rc.Policy)
					}
				}
			}
		}
		rc.isValid = err == nil
	}
	return
}

// NewClient should be called only after Validate has been called, to make sure
// that rc is a valid RetryConfig
func (rc *RetryConfig) NewClient(rt http.RoundTripper) (*http.Client, error) {
	c := rhttp.NewClient()
	if rc != nil {
		if !rc.isValid {
			return nil, fmt.Errorf("retry configuration is not valid")
		}
		c.RetryWaitMin = rc.WaitMin
		c.RetryWaitMax = rc.WaitMax
		c.RetryMax = rc.MaxAttempts
		c.Backoff = rc.backoff
	}
	c.HTTPClient = &http.Client{Transport: rt}
	return c.StandardClient(), nil
}

func validDuration(d time.Duration) (err error) {
	if d <= 0 {
		err = fmt.Errorf("duration %v must be positive", d)
	}
	return
}

func validPositive(n int) (err error) {
	if n <= 0 {
		err = fmt.Errorf("number %d must be positive", n)
	}
	return
}
