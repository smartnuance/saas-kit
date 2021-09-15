package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/volatiletech/null/v8"
)

type MySuite struct{}

type Mock struct {
	recorder []string
}

func (mock Mock) record(s string) {
	mock.recorder = append(mock.recorder, s)
}

func (s MySuite) TestSignup(assert, require *td.T) {
	// given
	ctrl := gomock.NewController(require.TB)
	mock := NewMockDBAPI(ctrl)

	mock.
		EXPECT().
		BeginTx(gomock.Any()).
		Return(nil, nil)

	user := &m.User{
		ID:    1,
		Name:  null.StringFrom("Yanis"),
		Email: "yanis@example.com",
	}
	mock.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Eq("Yanis"), gomock.Eq("yanis@example.com"), gomock.Not(gomock.Eq("test"))).
		Return(user, nil)

	mock.EXPECT().
		CreateProfile(gomock.Any(), gomock.Any(), gomock.Eq(int64(2)), gomock.Eq(user), gomock.Eq("teacher")).
		Return(&m.Profile{
			ID:         2,
			UserID:     1,
			InstanceID: 2,
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
	userID, err := service.signup(ctx, 2, SignupBody{Name: "Yanis", Email: "yanis@example.com", Password: "test"}, "teacher")
	assert.CmpNoError(err)
	assert.CmpLax(userID, 1)
}

func TestMySuite(t *testing.T) {
	tdsuite.Run(t, MySuite{})
}
