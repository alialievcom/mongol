package sites_core

import "testing"

func TestStartApp(t *testing.T) {
	StartApp("./config/config_test.yml")
}
