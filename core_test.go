package mongol

import "testing"

func TestStartApp(t *testing.T) {
	StartApp("./integration/config_test.yml")
}
