package authspec

import (
	"net/http"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/integration"
)

// ValidSessionKey is used in tests to verify a login response
const ValidSessionKey = "adm:7ff10abb653dead4186089acbd2b7891"

// LoginGoodSpecfunc ...
func LoginGoodSpecfunc(spec *integration.Spec) {
	adminGood := api.Admin{Username: "admin", Password: "kthtest"}
	request := integration.NewRequest(http.MethodPost, "/v0/admin/login", adminGood)
	response := integration.NewResponse(http.StatusOK, AdminLoginGood).
		WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
	spec.AssertAPICall(request, response)
}
