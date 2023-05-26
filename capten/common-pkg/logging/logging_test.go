package logging

import (
	"testing"
	"time"
)

func TestDebugf(t *testing.T) {
	l := NewLogger()
	l.Debugf("hi day: %v", time.Now().Day())
}

func TestInfof(t *testing.T) {
	l := NewLogger()
	l.Infof("hi day: %v", time.Now().Day())
}

func TestErrorf(t *testing.T) {
	l := NewLogger()
	l.Errorf("hi day: %v", time.Now().Day())
}

func TestAudit(t *testing.T) {
	l := NewLogger()
	l.Audit("audit", "dummy", "success", "system", "hi day: %v", time.Now().Day())
}

func TestInfo(t *testing.T) {
	l := NewLogger()
	l.Infof("hi day: %v", time.Now().Day())
}
