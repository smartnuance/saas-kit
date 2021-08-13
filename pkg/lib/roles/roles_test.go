package roles

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_initRoles(t *testing.T) {
	tests := []struct {
		name           string
		inheritedRoles map[string][]inheritedRole
		expClosure     ClosureMap
		expSwitchRoles ClosureMap
	}{
		{
			name:           "default",
			inheritedRoles: InheritedRoles,
			expClosure: ClosureMap{
				"event organizer": {
					"event organizer": true,
					"teacher":         true,
					"anonymous":       true,
				},
				"instance admin": {
					"event organizer": true,
					"instance admin":  true,
					"teacher":         true,
					"anonymous":       true,
				},
				"super admin": {
					"super admin": true,
					"anonymous":   true,
				},
				"teacher": {
					"teacher":   true,
					"anonymous": true,
				},
				"anonymous": {
					"anonymous": true,
				},
			},
			expSwitchRoles: ClosureMap{
				"event organizer": {},
				"instance admin":  {},
				"super admin": {
					"instance admin": true,
				},
				"teacher":   {},
				"anonymous": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, s := initRoles(tt.inheritedRoles)
			if !reflect.DeepEqual(c, tt.expClosure) {
				t.Errorf("got = \n%s;\nwant \n%s", pp(c), pp(tt.expClosure))
			}
			if !reflect.DeepEqual(s, tt.expSwitchRoles) {
				t.Errorf("got = \n%s;\nwant \n%s", pp(s), pp(tt.expSwitchRoles))
			}
		})
	}
}

func pp(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return err.Error()
}
