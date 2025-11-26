package util

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type ActionLogger struct {
	file *os.File
	mu   sync.Mutex
}

type ActionLog struct {
	Timestamp      string `json:"timestamp"`
	RemoteAddress  string `json:"remote_address,omitempty"`
	OpCode         int32  `json:"op_code"`
	Command        string `json:"command"`
	Collection     string `json:"collection"`
	IsMaster       bool   `json:"is_master"`
	Unacknowledged bool   `json:"unacknowledged"`
	RequestSize    int    `json:"request_size"`
	ResponseSize   int    `json:"response_size,omitempty"`
	WireMessage    string `json:"wire_message"`
}

func NewActionLogger(path string) (*ActionLogger, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &ActionLogger{file: f}, nil
}

func (a *ActionLogger) Log(entry ActionLog) {
	if a == nil {
		return
	}
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	a.mu.Lock()
	defer a.mu.Unlock()
	_ = json.NewEncoder(a.file).Encode(entry)
}

func (a *ActionLogger) Close() error {
	if a == nil {
		return nil
	}
	return a.file.Close()
}
