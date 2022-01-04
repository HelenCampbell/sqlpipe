package serve

import (
	"expvar"
	"net/http"

	"github.com/calmitchell617/sqlpipe/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	commonMiddleware := alice.New(app.metrics, app.recoverPanic, app.logRequest, app.rateLimit)

	apiStandardMiddleware := alice.New(app.requireAuthApi)
	apiRequireAdmin := apiStandardMiddleware.Append(app.requireAdmin)

	uiStandardMiddleware := alice.New(secureHeaders, app.session.Enable, noSurf, app.authenticateUi)

	router := httprouter.New()

	router.Handler(http.MethodPost, "/api/v1/users", apiRequireAdmin.ThenFunc(app.createUserApiHandler))
	router.Handler(http.MethodGet, "/api/v1/users", apiRequireAdmin.ThenFunc(app.listUsersApiHandler))
	router.Handler(http.MethodGet, "/api/v1/users/:id", apiRequireAdmin.ThenFunc(app.showUserApiHandler))
	router.Handler(http.MethodPut, "/api/v1/users", apiRequireAdmin.ThenFunc(app.updateUserApiHandler))
	router.Handler(http.MethodDelete, "/api/v1/users/:id", apiRequireAdmin.ThenFunc(app.deleteUserApiHandler))

	router.Handler(http.MethodGet, "/ui/users", uiStandardMiddleware.ThenFunc(app.listUsersUiHandler))

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)
	router.Handler(http.MethodGet, "/api/v1/debug/vars", expvar.Handler())

	router.NotFound = http.FileServer(http.FS(ui.Files))
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	return commonMiddleware.Then(router)
}
