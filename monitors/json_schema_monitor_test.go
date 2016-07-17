package monitors

import (
	core "github.com/gerty-monit/core"
	"os"
	"testing"
)

func TestShouldPingValidUrl(t *testing.T) {
	url := "https://api.github.com/rate_limit"
	schema := "file://" + os.Getenv("GOPATH") + "/src/github.com/gerty-monit/extras/monitors/github-limits.schema.json"
	monitor := NewJsonSchemaMonitor("Gerty Repositories Api", "this monitor checks the GitHub Repository API Schema", url, schema)
	status := monitor.Check()
	if status != core.OK {
		t.Fatalf("error while checking url %s", url)
	}
}
