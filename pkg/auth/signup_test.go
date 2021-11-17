package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"github.com/rs/xid"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/volatiletech/null/v8"
)

func TestMySuite(t *testing.T) {
	tdsuite.Run(t, &MySuite{})
}

type MySuite struct{}

func (s *MySuite) Test_signup(assert, require *td.T) {
	// given
	ctrl := gomock.NewController(require.TB)
	mock := NewMockDBAPI(ctrl)

	userID := xid.New().String()
	instanceID := xid.New().String()

	mock.
		EXPECT().
		BeginTx(gomock.Any()).
		Return(nil, nil)

	user := &m.User{
		ID:    userID,
		Name:  null.StringFrom("Yanis"),
		Email: "yanis@example.com",
	}
	mock.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Eq("Yanis"), gomock.Eq("yanis@example.com"), gomock.Not(gomock.Eq("test"))).
		Return(user, nil)

	mock.EXPECT().
		CreateProfile(gomock.Any(), gomock.Any(), gomock.Eq(instanceID), gomock.Eq(user), gomock.Eq("teacher")).
		Return(&m.Profile{
			ID:         xid.New().String(),
			UserID:     user.ID,
			InstanceID: instanceID,
			Role:       null.StringFrom("teacher"),
		}, nil)

	mock.
		EXPECT().
		Commit(gomock.Any()).
		Return(nil)

	service := Service{
		DBAPI: mock,
	}
	ctx := &gin.Context{}

	// when
	userID, err := service.signup(ctx, instanceID, SignupBody{Name: "Yanis", Email: "yanis@example.com", Password: "test"}, "teacher")

	// then
	assert.CmpNoError(err)
	assert.CmpLax(userID, user.ID)
}
