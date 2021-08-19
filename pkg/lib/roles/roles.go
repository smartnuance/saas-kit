package roles

import (
	"github.com/pkg/errors"

	"container/list"

	"github.com/gin-gonic/gin"
)

const (
	RoleSuperAdmin     = "super admin"
	RoleInstanceAdmin  = "instance admin"
	RoleEventOrganizer = "event organizer"
	RoleTeacher        = "teacher"
	RoleAnonymous      = "anonymous"
)

var Roles = []string{
	RoleSuperAdmin,
	RoleInstanceAdmin,
	RoleEventOrganizer,
	RoleTeacher,
	RoleAnonymous,
}

// inheritedRole builds a DAG of role inheritance with transitive permission propagation.
type inheritedRole struct {
	// Role is the role to inherit permissions from.
	Role string
	// ExplicitSwitch defines if the user has to explicitly switch to inherited role to receive its permissions.
	ExplicitSwitch bool
}

// InheritedRoles describes the role inheritance DAG.
// All roles implicitly inherit from RoleAnonymous without stating it here. RoleAnonymous makes the inheritance DAG rooted.
var InheritedRoles = map[string][]inheritedRole{
	RoleSuperAdmin: {
		inheritedRole{
			Role:           RoleInstanceAdmin,
			ExplicitSwitch: true,
		},
	},
	RoleInstanceAdmin: {
		inheritedRole{
			Role:           RoleEventOrganizer,
			ExplicitSwitch: false,
		},
	},
	RoleEventOrganizer: {
		inheritedRole{
			Role:           RoleTeacher,
			ExplicitSwitch: false,
		},
	},
}

type ClosureMap map[string]map[string]bool

// ImplicitRolesClosure lists the transitive closure of each role's implicitly inherited roles.
// The map's structure is
//   current role -> ancestor role -> (if ancestor is in closure)
var ImplicitRolesClosure = ClosureMap{}

// SwitchRoles lists each role's ancestor roles allowed to switch to.
// The map's structure is
//   current role -> ancestor role -> (if current role can switch to ancestor role)
var SwitchRoles = ClosureMap{}

func init() {
	ImplicitRolesClosure, SwitchRoles = initRoles(InheritedRoles)
}

func initRoles(inheritedRoles map[string][]inheritedRole) (implicitRolesClosure ClosureMap, switchRoles ClosureMap) {
	implicitRolesClosure = ClosureMap{}
	switchRoles = ClosureMap{}

	// build closures in role inheritance graph
	for _, role := range Roles {
		implicitRolesClosure[role] = map[string]bool{
			// All roles implicitly inherit from RoleAnonymous.
			RoleAnonymous: true,
		}
		switchRoles[role] = map[string]bool{}

		// the roles encountered over implicit inheritance
		todoImp := list.New()
		todoImp.PushBack(role)

		// the roles encountered over inheritance with explicit switch required
		todoExp := list.New()

		// track which roles has been reached in inheritance graph during traversal
		done := map[string]bool{}

		// Breadth-first traversal to collect closure of implicitly inherited roles
		for todoImp.Len() > 0 {
			p_ := todoImp.Front()
			todoImp.Remove(p_)
			p := p_.Value.(string)
			done[p] = true

			implicitRolesClosure[role][p] = true
			for _, inheritedRole := range inheritedRoles[p] {
				if !done[inheritedRole.Role] {
					if inheritedRole.ExplicitSwitch {
						todoExp.PushBack(inheritedRole.Role)
					} else {
						todoImp.PushBack(inheritedRole.Role)
					}
				}
			}
		}

		// Find ancestor roles only reachable over an explicit inheritance.
		// When an ancestor role P was already reached via a path of implicit inheritance, any explicit inheritance of P has no effect.
		for p_ := todoExp.Front(); p_ != nil; p_ = p_.Next() {
			p := p_.Value.(string)
			if !done[p] {
				switchRoles[role][p] = true
			}
		}
	}
	return
}

// CanActAs checks if the user's role can act in the given role implicitly.
func CanActAs(userRole, role string) bool {
	_, ok := ImplicitRolesClosure[userRole][role]
	return ok
}

// CanSwitchTo checks if the user's role can switch to a given role and acquiring those role's permissions.
func CanSwitchTo(userRole, role string) bool {
	_, ok := SwitchRoles[userRole][role]
	return ok
}

// ActAs checks if the context's user can act in the given targetRole implicitly.
// Otherwise returns ErrForbiddenRole or ErrRoleSwitchRequired.
func ActAs(c *gin.Context, targetRole string) error {
	role, err := Role(c)
	if err != nil {
		return err
	}

	_, ok := ImplicitRolesClosure[role]
	if !ok {
		if SwitchRoles[role][role] {
			return ErrRoleSwitchRequired
		} else {
			return ErrForbiddenRole
		}
	}

	return nil
}

// Role checks if the context's user can act in the given role implicitly.
// Otherwise returns ErrNotAllowed or ErrImplicitImpersonationNotAllowed.
func Role(c *gin.Context) (string, error) {
	role := c.GetString("role")
	if role == "" {
		return "", ErrMissingRole
	}

	return role, nil
}

var (
	ErrMissingRole        = errors.New("missing role in context")
	ErrForbiddenRole      = errors.New("forbidden to act as role")
	ErrRoleSwitchRequired = errors.New("switch to role required")
)
