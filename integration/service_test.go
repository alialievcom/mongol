package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AliAlievMos/mongol"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func init() {
	err := os.Setenv("MONGO_URI", "mongodb://root:password@127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	go func() {
		mongol.StartApp("./config_test.yml")
	}()
	time.Sleep(1 * time.Second)
}

func TestRoles(t *testing.T) {
	type in struct {
		log, pass string
	}
	var cases = []struct {
		name       string
		in         in
		needStatus int
		collection string
	}{
		{
			name:       "with star",
			in:         in{"admin", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "dictionaries",
		},
		{
			name:       "with delete",
			in:         in{"user_can_delete", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "dictionaries",
		},
		{
			name:       "cant delete",
			in:         in{"user_can_get", "adminpass"},
			needStatus: http.StatusUnauthorized,
			collection: "dictionaries",
		},
		{
			name:       "nil role",
			in:         in{"user_with_nil", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "dictionaries",
		},
		{
			name:       "empty slice",
			in:         in{"user_with_empty", "adminpass"},
			needStatus: http.StatusUnauthorized,
			collection: "dictionaries",
		},

		{
			name:       "with star",
			in:         in{"admin", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "publication",
		},
		{
			name:       "with delete",
			in:         in{"user_can_delete", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "publication",
		},
		{
			name:       "cant delete",
			in:         in{"user_can_get", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "publication",
		},
		{
			name:       "nil role",
			in:         in{"user_with_nil", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "publication",
		},
		{
			name:       "empty slice",
			in:         in{"user_with_empty", "adminpass"},
			needStatus: http.StatusBadRequest,
			collection: "publication",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := loginReq(tc.in.log, tc.in.pass)
			if err != nil {
				t.Error(err)
				return
			}
			client := &http.Client{}
			reqDelete, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/authenticated/%s/1", tc.collection), nil)
			if err != nil {
				t.Error(err)
				return
			}
			reqDelete.Header.Set("Authorization", token)
			resp, err := client.Do(reqDelete)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.StatusCode != tc.needStatus {
				t.Errorf("not expected status:%d", resp.StatusCode)
				return
			}
		})

	}
}

func loginReq(log, pass string) (token string, err error) {
	resp, err := http.Post("http://localhost:8080/login", "application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"login":"%s", "password": "%s"}`, log, pass))))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("not success status:%d", resp.StatusCode)
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var login map[string]any
	err = json.Unmarshal(bodyBytes, &login)
	if err != nil {
		return
	}
	token, ok := login["access"].(string)
	if !ok {
		err = fmt.Errorf("cant parse accsees: %v", login)
		return
	}
	return token, err
}
