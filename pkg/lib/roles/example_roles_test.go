package roles_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
)

// Demonstrates the parsing of the current instance and user's role and a successful authorization check.
func Example_successfulCheck() {
	// Create dummy context
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Set("user", "user-guid")
	ctx.Set("role", "teacher")
	ctx.Set("instance", "instance-guid")

	// Check permission to revoke token for potentially different user
	if !(roles.CanActFor(ctx, "instance-guid") && roles.CanActIn(ctx, roles.RoleTeacher)) {
		fmt.Println("unauthorized")
		return
	}

	// Do something a teacher can do
	fmt.Println("authorized")

	// Output: authorized
}
