package roles

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
)

func TestMySuite(t *testing.T) {
	tdsuite.Run(t, MySuite{})
}

type MySuite struct{}

func (s MySuite) Test_initRoles(assert, require *td.T) {
	tests := []struct {
		name               string
		inheritanceDAG     map[string][]edge
		inheritanceClosure closure
		switchableRoles    closure
	}{
		{
			name:           "default",
			inheritanceDAG: inheritanceDAG,
			inheritanceClosure: closure{
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
			switchableRoles: closure{
				"event organizer": {},
				"instance admin":  {},
				"super admin": {
					"instance admin":  true,
					"event organizer": true,
					"teacher":         true,
				},
				"teacher": {},
				"":        {},
			},
		},
	}
	for _, test := range tests {
		assert.Run(test.name, func(t *td.T) {
			c, s := initRoles(test.inheritanceDAG)

			t.CmpDeeply(c, test.inheritanceClosure)
			t.CmpDeeply(s, test.switchableRoles)
		})
	}
}
