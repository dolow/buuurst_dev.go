package collector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dolow/buuurst_dev.go/config"
	"github.com/dolow/buuurst_dev.go/log"
)

const (
	putUrl = "https://lambda-public.buuurst.dev/put-request-log"
)

type ResponseParam interface {
	StatusCode() int
	Header() http.Header
}

type Collector struct {
	Enable       bool     `json:"enable"`
	ProjectID    string   `json:"project_id"`
	ServiceKey   string   `json:"service_key"`
	CustomHeader []string `json:"custom_header"`
	IgnorePaths  []string `json:"appignore_paths"`

	TimeStamp *time.Time
	Method    string         `json:"method"`
	Path      string         `json:"path"`
	Status    int            `json:"status"`
	Query     url.Values     `json:"query"`
	Cookie    []*http.Cookie `json:"cookie"`
	RequestID string         `json:"request_id"`

	HttpHeaders http.Header
	CgiHeaders  http.Header
	Headers     []string
	Body        []byte
}

func New() *Collector {
	instance := &Collector{
		ProjectID:  config.Current.ProjectID,
		ServiceKey: config.Current.ServiceKey,
	}

	if config.Current.CustomHeader == nil {
		instance.CustomHeader = []string{}
	} else {
		instance.CustomHeader = config.Current.CustomHeader
	}

	if config.Current.IgnorePaths == nil {
		instance.IgnorePaths = []string{}
	} else {
		instance.IgnorePaths = config.Current.IgnorePaths
	}

	return instance
}

func (c *Collector) Collect(p ResponseParam, w http.ResponseWriter, r *http.Request, body []byte) error {
	if !c.Enable {
		return nil
	}

	if !c.ShouldSendLog() {
		return nil
	}

	c.SetRequestLog(r, body)
	c.SetResponseLog(p)

	if err := c.SendLog(); err != nil {
		return err
	}

	return nil
}

func (c *Collector) ShouldSendLog() bool {
	if len(c.IgnorePaths) == 0 {
		return true
	}

	for i := range c.IgnorePaths {
		if c.IgnorePaths[i] == c.Path {
			return false
		}
	}

	return true
}

func (c *Collector) SetRequestLog(r *http.Request, reqBody []byte) {
	now := time.Now()

	c.TimeStamp = &now
	c.Method = r.Method
	c.Path = r.URL.Path
	c.Query = r.URL.Query()
	c.Cookie = r.Cookies() // Rack::Utils.parse_cookies(env)

	for key, values := range r.Header {
		v := strings.Join(values, " ")
		c.Headers = append(c.Headers, v)
		if key == "HTTP_X_REQUEST_ID" {
			c.RequestID = v
		}
	}

	// TODO: custom header

	c.Body = reqBody
}
func (c *Collector) SetResponseLog(p ResponseParam) {
	if c.RequestID == "" {
		headers := p.Header()
		if values, exists := headers["X-Request-Id"]; exists {
			if len(values) > 0 {
				c.RequestID = strings.Join(values, " ")
			}
		}
	}
	c.Status = p.StatusCode()
}

func (c *Collector) SendLog() error {
	param, err := c.CreateParam()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", putUrl, bytes.NewBuffer(param))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	if _, err := client.Do(req); err != nil {
		return err
	}

	return nil
}

func (c *Collector) CreateParam() ([]byte, error) {
	return json.Marshal(&log.Param{
		ProjectID:   c.ProjectID,
		RequestedAt: c.TimeStamp.Format(time.RFC3339), // TODO: format
		Method:      c.Method,
		Path:        c.Path,
		Query:       c.Query,
		Cookie:      c.Cookie,
		RequestID:   c.RequestID,
		Status:      c.Status,
		Headers:     c.Headers,
		ServiceKey:  c.ServiceKey,
		Body:        c.Body,
	})
}
