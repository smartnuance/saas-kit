package roles

import (
	"container/list"

	"github.com/gin-gonic/gin"
)

const (
	RoleSuperAdmin     = "super admin"
	RoleInstanceAdmin  = "instance admin"
	RoleEventOrganizer = "event organizer"
	RoleTeacher        = "teacher"
	NoRole             = ""
)

var Roles = []string{
	RoleSuperAdmin,
	RoleInstanceAdmin,
	RoleEventOrganizer,
	RoleTeacher,
	NoRole,
}

// inheritedRole builds a DAG of role inheritance with transitive permission propagation.
type inheritedRole struct {
	// Role is the role to inherit permissions from.
	Role string
	// ExplicitSwitch defines if the user has to explicitly switch to inherited role to receive its permissions.
	ExplicitSwitch bool
}

// inheritedRoles describes the role inheritance DAG.
// All roles implicitly inherit from RoleAnonymous without stating it here. RoleAnonymous makes the inheritance DAG rooted.
var inheritedRoles = map[string][]inheritedRole{
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

// implicitRolesClosure lists the transitive closure of each role's implicitly inherited roles.
// The map's structure is
//   current role -> ancestor role -> (if ancestor is in closure)
var implicitRolesClosure = ClosureMap{}

// switchRoles lists each role's ancestor roles allowed to switch to.
// The map's structure is
//   current role -> ancestor role -> (if current role can switch to ancestor role)
var switchRoles = ClosureMap{}

func init() {
	implicitRolesClosure, switchRoles = initRoles(inheritedRoles)
}

func initRoles(inheritedRoles map[string][]inheritedRole) (implicitRolesClosure ClosureMap, switchRoles ClosureMap) {
	implicitRolesClosure = ClosureMap{}
	switchRoles = ClosureMap{}

	// build closures in role inheritance graph
	for _, role := range Roles {
		implicitRolesClosure[role] = map[string]bool{
			// All roles implicitly inherit from NoRole.
			NoRole: true,
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

// CanSwitchTo checks if the user's role can switch to a given role and acquiring those role's permissions.
func CanSwitchTo(userRole string, role string) bool {
	_, ok := switchRoles[userRole][role]
	return ok
}

// CanActIn checks if the user can act in the desired role implicitly.
func CanActIn(ctx *gin.Context, role string) bool {
	userRole, _, ok := FromContext(ctx)
	if !ok {
		return false
	}

	_, ok = implicitRolesClosure[userRole][role]
	return ok
}

// CanActFor checks if the user can act for the desired instance.
func CanActFor(ctx *gin.Context, instanceID int) bool {
	_, userInstance, ok := FromContext(ctx)
	if !ok {
		return false
	}

	return userInstance == instanceID
}

// FromContext retrieves user's role from context's claims.
// Returns ok equals true when the instance was found in context (since there is no default instance).
// The default role on the other hand ist NoRole.
func FromContext(ctx *gin.Context) (role string, instanceID int, ok bool) {
	role = ctx.GetString("role") // corresponds to NoRole if empty
	var instance interface{}
	instance, ok = ctx.Get("instance") // should never be empty
	if !ok {
		return
	}
	instanceID = instance.(int)
	return
}
