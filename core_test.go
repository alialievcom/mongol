package mongol

import "testing"

func TestStartApp(t *testing.T) {
	StartApp("./config/config_test.yml")
}
