/*
Package roles manages user's permissions depending on its roles.
Roles are organized hierarchically and are implicitely or explicitely inherited.

From implicit inheritance a user receives the union of all capabilities of
transitively inherited roles. This inheritance is transparent to both the user's
current active role and code that does checks upon a possibly inherited role to allow/prohibid
an action that requires authorization.

From explicit inheritance a user does not transparently inherit any capabilities.
The are only signalers to which role a user might switch to receive those (and
implicitely) inherited capabilities.

Some syntactic sugar for gin.Context is provided to allow simple authorization checks like

	roles.CanActIn(ctx, roles.RoleTeacher)

	roles.CanActFor(ctx, instanceID)

and permanent switching to another role

	roles.SwitchTo(ctx, roles.RoleInstanceAdmin)

*/
package roles
