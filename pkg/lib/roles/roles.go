package roles

import (
	"container/list"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	RoleKey     = "role"
	InstanceKey = "instance"
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

type closureMap map[string]map[string]bool

// implicitRolesClosure lists the transitive closure of each role's implicitly inherited roles.
// The map's structure is
//   current role -> ancestor role -> (if ancestor is in closure)
var implicitRolesClosure = closureMap{}

// switchRoles lists each role's ancestor roles allowed to switch to.
// The map's structure is
//   current role -> ancestor role -> (if current role can switch to ancestor role)
var switchRoles = closureMap{}

func init() {
	implicitRolesClosure, switchRoles = initRoles(inheritedRoles)
}

func initRoles(inheritedRoles map[string][]inheritedRole) (implicitRolesClosure closureMap, switchRoles closureMap) {
	implicitRolesClosure = closureMap{}
	switchRoles = closureMap{}

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

func valid(role string) bool {
	_, ok := implicitRolesClosure[role]
	return ok
}

// CanSwitchTo checks if the user's role can switch to a targetRole and acquiring those role's permissions.
// Switching is allowed when there is an implicit path from userRole to role
// or userrole directly, explicitly inherits targetRole.
func CanSwitchTo(userRole string, targetRole string) bool {
	_, okImplicit := implicitRolesClosure[userRole][targetRole]
	_, okExplicit := switchRoles[userRole][targetRole]
	return okImplicit || okExplicit
}

// SwitchTo attempts to switch to a temporary targetRole.
// The user's role defined in context is checked against the rules defining if switching is allowed.
// The temporary role is set on the context under the "role" key, overwriting the original role.
func SwitchTo(ctx *gin.Context, targetRole string) error {
	role, _, err := Get(ctx)
	if err != nil {
		return err
	}
	if !CanSwitchTo(role, targetRole) {
		return ErrSwitchNotAllowed
	}
	ctx.Set("role", targetRole)
	return nil
}

// CanActIn checks if the user can act in the desired targetRole implicitly.
func CanActIn(ctx *gin.Context, targetRole string) bool {
	role, _, err := Get(ctx)
	if err != nil {
		return false
	}

	_, ok := implicitRolesClosure[role][targetRole]
	return ok
}

// CanActFor checks if the user can act for the desired instance.
func CanActFor(ctx *gin.Context, instanceID string) bool {
	_, userInstance, err := Get(ctx)
	if err != nil {
		return false
	}

	return userInstance == instanceID
}

// Get retrieves the user's role and instance to act for from context.
// The default role is NoRole. An invalid role results in ErrInvalidRole.
// There is no default instance. An invalid instance results in ErrInvalidInstance.
func Get(ctx *gin.Context) (role string, instanceID string, err error) {
	role = ctx.GetString(RoleKey) // corresponds to NoRole if empty
	if !valid(role) {
		err = ErrInvalidRole
		return
	}
	instance, ok := ctx.Get(InstanceKey) // should never be empty
	if !ok {
		err = ErrInvalidInstance
		return
	}
	instanceID = instance.(string)
	return
}

// FromHeaders parses headers to retrieve user's temporary role and instance to act for,
// overwriting default role/instance from context.
// When role parameter is missing, falls back to role specified in context.
// When instance parameter is missing, falls back to instance specified in context.
// Returns an error when neither parameter nor fallback was provided for role or instance.
func FromHeaders(ctx *gin.Context) (role string, instanceID string, err error) {
	role = ctx.GetHeader(RoleKey)
	if len(role) > 0 {
		if !valid(role) {
			err = ErrInvalidRole
			return
		}
	} else {
		// if no instance is provided, fallback to role from context
		role, _, err = Get(ctx)
		if err != nil {
			return
		}
	}

	instanceID = ctx.GetHeader(InstanceKey)
	if len(instanceID) == 0 {
		// if no instance is provided, fallback to instance from context
		_, instanceID, err = Get(ctx)
		if err != nil {
			return
		}
	}
	return
}

var (
	ErrMissingRole      = errors.New("missing role")
	ErrInvalidRole      = errors.New("invalid role provided")
	ErrMissingInstance  = errors.New("missing instance")
	ErrInvalidInstance  = errors.New("invalid instance provided")
	ErrSwitchNotAllowed = errors.New("role switch not allowed")
)
