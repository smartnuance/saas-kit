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
		expClosure     closureMap
		expSwitchRoles closureMap
	}{
		{
			name:           "default",
			inheritedRoles: inheritedRoles,
			expClosure: closureMap{
				"event organizer": {
					"event organizer": true,
					"teacher":         true,
					"":                true,
				},
				"instance admin": {
					"event organizer": true,
					"instance admin":  true,
					"teacher":         true,
					"":                true,
				},
				"super admin": {
					"super admin": true,
					"":            true,
				},
				"teacher": {
					"teacher": true,
					"":        true,
				},
				"": {
					"": true,
				},
			},
			expSwitchRoles: closureMap{
				"event organizer": {},
				"instance admin":  {},
				"super admin": {
					"instance admin": true,
				},
				"teacher": {},
				"":        {},
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
