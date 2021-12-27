package roles

import (
	"container/list"

	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
)

const (
	RoleHeader     = "role"
	InstanceHeader = "instance"

	UserKey     = "user"
	RoleKey     = "role"
	InstanceKey = "instance"
)

type Role string

const (
	RoleSuperAdmin     Role = "super admin"
	RoleInstanceAdmin  Role = "instance admin"
	RoleEventOrganizer Role = "event organizer"
	RoleTeacher        Role = "teacher"
	NoRole             Role = ""
)

var Roles = []Role{
	RoleSuperAdmin,
	RoleInstanceAdmin,
	RoleEventOrganizer,
	RoleTeacher,
	NoRole,
}

// edge builds a DAG of role inheritance with transitive permission propagation.
type edge struct {
	// Role is the role to inherit permissions from.
	Role Role
	// SwitchRequired defines if the user has to explicitly switch to inherited role to receive its permissions.
	SwitchRequired bool
}

// dag represents a DAG.
type dag map[Role][]edge

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

type closure map[Role]map[Role]struct{}

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

func initRoles(inheritanceDAG map[Role][]edge) (inheritanceClosure closure, switchRoles closure) {
	inheritanceClosure = closure{}
	switchRoles = closure{}

	// build closures in role inheritance graph
	for _, role := range Roles {
		inheritanceClosure[role] = map[Role]struct{}{
			// All roles implicitly inherit from NoRole.
			NoRole: {},
		}
		switchRoles[role] = map[Role]struct{}{
			role: {},
		}

		// the roles encountered over implicit inheritance
		todo := list.New()
		todo.PushBack(role)

		// the roles encountered over inheritance with explicit switch required
		switchableTodo := list.New()

		// track which roles has been reached in inheritance graph during traversal
		done := map[Role]bool{}

		// Breadth-first traversal to collect closure of inherited roles
		for p_ := todo.Front(); p_ != nil; p_ = p_.Next() {
			p := p_.Value.(Role)
			done[p] = true

			inheritanceClosure[role][p] = struct{}{}
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
			p := p_.Value.(Role)
			if !done[p] {
				switchRoles[role][p] = struct{}{}
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

func valid(role Role) bool {
	_, ok := inheritanceClosure[role]
	return ok
}

// rolesSpecEntry describes a user's target roles for switching and each of those roles inherited roles.
type rolesSpecEntry struct {
	Role      Role   `json:"role"`
	Inherited []Role `json:"inherited"`
}

// rolesSpec are the switching target roles, ordered from more permissive to more specific roles.
type rolesSpec []rolesSpecEntry

// RolesSpec returns a RolesSpec for the given user role.
func RolesSpec(userRole Role) (spec rolesSpec) {
	if !valid(userRole) {
		return
	}

	for _, r := range Roles {
		if _, ok := switchRoles[userRole][r]; ok {
			spec = append(spec, rolesSpecEntry{
				Role:      r,
				Inherited: InheritedRoles(r),
			})
		}
	}
	return
}

// CanSwitchTo checks if the user's role can switch to a targetRole acquiring those role's permissions.
// Switching is allowed when there is an implicit path from userRole to role
// or userrole directly, explicitly inherits targetRole.
func CanSwitchTo(userRole Role, targetRole Role) bool {
	_, okImplicit := inheritanceClosure[userRole][targetRole]
	_, okExplicit := switchRoles[userRole][targetRole]
	return okImplicit || okExplicit
}

// SwitchTo attempts to switch to a temporary targetRole.
// The user's role defined in context is checked against the rules defining if switching is allowed.
// The temporary role is set on the context under the "role" key, overwriting the original role.
func SwitchTo(ctx *gin.Context, targetRole Role) error {
	role, err := FromContext(ctx)
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
	userID, err := User(ctx)
	if err != nil {
		return false
	}

	return userID == targetUserID
}

// CanActIn checks if the user can act in the desired targetRole without switching to that role.
func CanActIn(ctx *gin.Context, targetRole Role) bool {
	role, err := FromContext(ctx)
	if err != nil {
		return false
	}

	_, ok := inheritanceClosure[role][targetRole]
	return ok
}

// InheritedRoles returns an unordered list of roles a given user role can act in, without switching.
// The returned list is a copy and can be safely modified.
func InheritedRoles(userRole Role) []Role {
	keys := make([]Role, 0, len(inheritanceClosure[userRole]))
	for k := range inheritanceClosure[userRole] {
		keys = append(keys, k)
	}

	return keys
}

// CanActFor checks if the user can act for the desired instance.
func CanActFor(ctx *gin.Context, instanceID string) bool {
	userInstance, err := Instance(ctx)
	if err != nil {
		return false
	}

	if userInstance == instanceID {
		return true
	}

	userRole, err := FromContext(ctx)
	if err != nil {
		return false
	}
	return userRole == RoleSuperAdmin
}

// User retrieves the user from context.
// There is no default user. When no user is registerd in context, this results in ErrMissingUser.
func User(ctx *gin.Context) (string, error) {
	userID_, ok := ctx.Get(UserKey) // should exist
	if !ok {
		return "", ErrMissingUser
	}
	return userID_.(string), nil
}

// FromContext retrieves the role from context.
// The default role is NoRole. An invalid role results in ErrInvalidRole.
func FromContext(ctx *gin.Context) (Role, error) {
	role_, ok := ctx.Get(RoleKey) // corresponds to NoRole if empty
	if !ok {
		role_ = NoRole
	}
	role := role_.(Role)
	if !valid(role) {
		return "", ErrInvalidRole
	}
	return role, nil
}

// Instance retrieves the instance to act for from context.
// There is no default instance. An invalid instance results in ErrMissingInstance.
func Instance(ctx *gin.Context) (string, error) {
	instanceID_, ok := ctx.Get(InstanceKey) // should exist
	if !ok {
		return "", ErrMissingInstance
	}
	return instanceID_.(string), nil
}

var (
	ErrMissingUser      = errors.New("missing user")
	ErrInvalidRole      = errors.New("invalid role provided")
	ErrMissingInstance  = errors.New("missing instance")
	ErrSwitchNotAllowed = errors.New("role switch not allowed")
	ErrUnauthorized     = errors.New("role insufficient to act on desired instance")
)
