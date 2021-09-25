package roles_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
)

// Demonstrates the parsing of the current instance and user's role and a successful authorization check.
func Example_switchRole() {
	// Create dummy context
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Set("user", "user-guid")
	ctx.Set("role", "super admin")
	ctx.Set("instance", "instance-guid")

	// Check permission to revoke token for potentially different user
	if !roles.CanActFor(ctx, "instance-guid") {
		fmt.Println("instance unauthorized")
		return
	}

	if !roles.CanActIn(ctx, roles.RoleTeacher) {
		fmt.Println("unauthorized to act as teacher")
	}

	err := roles.SwitchRole(ctx, roles.RoleInstanceAdmin)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("switched to instance admin")

	if !roles.CanActIn(ctx, roles.RoleTeacher) {
		fmt.Println("still unauthorized to act as teacher")
		return
	}

	// Do something a teacher can do
	fmt.Println("authorized")

	// Output:
	// unauthorized to act as teacher
	// switched to instance admin
	// authorized
}
