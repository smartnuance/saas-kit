package roles

import (
	"reflect"
	"testing"

	"github.com/smartnuance/saas-kit/pkg/lib"
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
			inheritedRoles: inheritedRoles,
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
				t.Errorf("got = \n%s;\nwant \n%s", lib.PP(c), lib.PP(tt.expClosure))
			}
			if !reflect.DeepEqual(s, tt.expSwitchRoles) {
				t.Errorf("got = \n%s;\nwant \n%s", lib.PP(s), lib.PP(tt.expSwitchRoles))
			}
		})
	}
}
