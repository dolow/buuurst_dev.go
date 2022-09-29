package log

type Param struct {
	ProjectID   string            `json:"project_id"`
	RequestedAt string            `json:"requested_at"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Query       map[string]string `json:"query"`
	Cookie      map[string]string `json:"cookie"`
	RequestID   string            `json:"request_id"`
	Status      int               `json:"status"`
	Headers     []string          `json:"header"`
	ServiceKey  string            `json:"service_key"`
	Body        string            `json:"body"`
}
