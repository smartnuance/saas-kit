package roles

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
)

func TestMySuite(t *testing.T) {
	tdsuite.Run(t, &MySuite{})
}

type MySuite struct{}

func (s *MySuite) Test_initRoles(assert, require *td.T) {
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
					"event organizer": struct{}{},
					"teacher":         struct{}{},
					"":                struct{}{},
				},
				"instance admin": {
					"event organizer": struct{}{},
					"instance admin":  struct{}{},
					"teacher":         struct{}{},
					"":                struct{}{},
				},
				"super admin": {
					"super admin": struct{}{},
					"":            struct{}{},
				},
				"teacher": {
					"teacher": struct{}{},
					"":        struct{}{},
				},
				"": {
					"": struct{}{},
				},
			},
			switchableRoles: closure{
				"event organizer": {
					"event organizer": struct{}{},
				},
				"instance admin": {
					"instance admin": struct{}{},
				},
				"super admin": {
					"super admin":     struct{}{},
					"instance admin":  struct{}{},
					"event organizer": struct{}{},
					"teacher":         struct{}{},
				},
				"teacher": {
					"teacher": struct{}{},
				},
				"": {
					"": struct{}{},
				},
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
