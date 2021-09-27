package roles

import (
	"container/list"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	UserKey     = "user"
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

// edge builds a DAG of role inheritance with transitive permission propagation.
type edge struct {
	// Role is the role to inherit permissions from.
	Role string
	// SwitchRequired defines if the user has to explicitly switch to inherited role to receive its permissions.
	SwitchRequired bool
}

// dag represents a DAG.
type dag map[string][]edge

// inheritanceDAG describes the role inheritance DAG.
// All roles implicitly inherit from NoRole without stating it here. NoRole makes the inheritance DAG rooted.
var inheritanceDAG = dag{
	RoleSuperAdmin: {
		edge{
			Role:           RoleInstanceAdmin,
			SwitchRequired: true,
		},
	},
	RoleInstanceAdmin: {
		edge{
			Role:           RoleEventOrganizer,
			SwitchRequired: false,
		},
	},
	RoleEventOrganizer: {
		edge{
			Role:           RoleTeacher,
			SwitchRequired: false,
		},
	},
}

type closure map[string]map[string]bool

// inheritanceClosure lists the transitive closure of each role's inherited roles.
// The map's structure is
//   current role -> inherited role -> (true if inherited is in closure)
var inheritanceClosure = closure{}

// switchRoles lists each role's inherited roles allowed to switch to.
// The map's structure is
//   current role -> inherited role -> (true if current role can switch to inherited role)
var switchRoles = closure{}

func init() {
	inheritanceClosure, switchRoles = initRoles(inheritanceDAG)
}

func initRoles(inheritanceDAG map[string][]edge) (inheritanceClosure closure, switchRoles closure) {
	inheritanceClosure = closure{}
	switchRoles = closure{}

	// build closures in role inheritance graph
	for _, role := range Roles {
		inheritanceClosure[role] = map[string]bool{
			// All roles implicitly inherit from NoRole.
			NoRole: true,
		}
		switchRoles[role] = map[string]bool{}

		// the roles encountered over implicit inheritance
		todo := list.New()
		todo.PushBack(role)

		// the roles encountered over inheritance with explicit switch required
		switchableTodo := list.New()

		// track which roles has been reached in inheritance graph during traversal
		done := map[string]bool{}

		// Breadth-first traversal to collect closure of inherited roles
		for p_ := todo.Front(); p_ != nil; p_ = p_.Next() {
			p := p_.Value.(string)
			done[p] = true

			inheritanceClosure[role][p] = true
			for _, e := range inheritanceDAG[p] {
				if !done[e.Role] {
					if e.SwitchRequired {
						switchableTodo.PushBack(e.Role)
					} else {
						todo.PushBack(e.Role)
					}
				}
			}
		}

		// Find inherited roles only reachable over an explicit inheritance.
		// When an inherited role P was already reached via a path of implicit inheritance, any inheritance of P with switch required has no effect.
		for p_ := switchableTodo.Front(); p_ != nil; p_ = p_.Next() {
			p := p_.Value.(string)
			if !done[p] {
				switchRoles[role][p] = true
				for _, inheritedRole := range inheritanceDAG[p] {
					if !done[inheritedRole.Role] {
						switchableTodo.PushBack(inheritedRole.Role)
					}
				}
			}
		}
	}
	return
}

func valid(role string) bool {
	_, ok := inheritanceClosure[role]
	return ok
}

// CanSwitchTo checks if the user's role can switch to a targetRole and acquiring those role's permissions.
// Switching is allowed when there is an implicit path from userRole to role
// or userrole directly, explicitly inherits targetRole.
func CanSwitchTo(userRole string, targetRole string) bool {
	_, okImplicit := inheritanceClosure[userRole][targetRole]
	_, okExplicit := switchRoles[userRole][targetRole]
	return okImplicit || okExplicit
}

// SwitchTo attempts to switch to a temporary targetRole.
// The user's role defined in context is checked against the rules defining if switching is allowed.
// The temporary role is set on the context under the "role" key, overwriting the original role.
func SwitchTo(ctx *gin.Context, targetRole string) error {
	_, role, _, err := From(ctx)
	if err != nil {
		return err
	}
	if targetRole == role {
		return nil
	}
	if !CanSwitchTo(role, targetRole) {
		return ErrSwitchNotAllowed
	}
	ctx.Set(RoleKey, targetRole)
	return nil
}

// CanActAs checks if the user can act as a desired user.
func CanActAs(ctx *gin.Context, targetUserID string) bool {
	userID, _, _, err := From(ctx)
	if err != nil {
		return false
	}

	return userID == targetUserID
}

// CanActIn checks if the user can act in the desired targetRole without switching to that role.
func CanActIn(ctx *gin.Context, targetRole string) bool {
	_, role, _, err := From(ctx)
	if err != nil {
		return false
	}

	_, ok := inheritanceClosure[role][targetRole]
	return ok
}

// CanActFor checks if the user can act for the desired instance.
func CanActFor(ctx *gin.Context, instanceID string) bool {
	_, _, userInstance, err := From(ctx)
	if err != nil {
		return false
	}

	return userInstance == instanceID
}

// From retrieves the user, his role and the instance to act for from context.
// There is no default user. When no user is registerd in context, this results in ErrMissingUser.
// The default role is NoRole. An invalid role results in ErrInvalidRole.
// There is no default instance. An invalid instance results in ErrMissingInstance.
func From(ctx *gin.Context) (userID, role, instanceID string, err error) {
	userID_, ok := ctx.Get(UserKey) // should exist
	if !ok {
		err = ErrMissingUser
		return
	}
	userID = userID_.(string)
	role = ctx.GetString(RoleKey) // corresponds to NoRole if empty
	if !valid(role) {
		err = ErrInvalidRole
		return
	}
	instanceID_, ok := ctx.Get(InstanceKey) // should exist
	if !ok {
		err = ErrMissingInstance
		return
	}
	instanceID = instanceID_.(string)
	return
}

// ApplyHeaders parses headers to retrieve user's temporary role and instance to act for,
// overwriting default role/instance from context. Switches to specified role.
// When role parameter is missing, falls back to role specified in context set by role in JWT.
// When instance parameter is missing, falls back to instance specified in context.
// Returns an error when neither parameter nor fallback was provided for role or instance.
func ApplyHeaders(ctx *gin.Context) (role, instanceID string, err error) {
	role = ctx.GetHeader(RoleKey)
	if len(role) > 0 {
		if !valid(role) {
			err = ErrInvalidRole
			return
		}
	} else {
		// if no instance is provided, fallback to role from context
		_, role, _, err = From(ctx)
		if err != nil {
			return
		}
	}
	err = SwitchTo(ctx, role)
	if err != nil {
		return
	}

	instanceID = ctx.GetHeader(InstanceKey)
	if len(instanceID) == 0 {
		// if no instance is provided, fallback to instance from context
		_, _, instanceID, err = From(ctx)
		if err != nil {
			return
		}
	}
	return
}

var (
	ErrMissingUser      = errors.New("missing user")
	ErrInvalidRole      = errors.New("invalid role provided")
	ErrMissingInstance  = errors.New("missing instance")
	ErrSwitchNotAllowed = errors.New("role switch not allowed")
)
