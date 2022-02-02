package serve

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/calmitchell617/sqlpipe/internal/validator"
	"github.com/calmitchell617/sqlpipe/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readDateTime(qs url.Values, key string, defaultValue time.Time) time.Time {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	timeVal, err := time.Parse("2006-01-02T03:04", s)
	if err != nil {
		return defaultValue
	}

	return timeVal
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}
	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverErrorResponse(w, r, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	td.IsAdmin = app.isAdmin(r)

	return td
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type PaginationData struct {
	NeedsPagination bool
	Pages           []Page
	Offset          int
	PageName        string
}

type Page struct {
	PageNum  int
	IsActive bool
	Link     string
}

func getPaginationData(
	currentPage int,
	totalObjects int,
	pageSize int,
	currentPageName string,
) PaginationData {
	pages := []Page{}

	if totalObjects <= pageSize {
		return PaginationData{false, pages, 0, currentPageName}
	}

	numPages := int(math.Ceil(float64(totalObjects) / float64(pageSize)))

	offset := (pageSize * currentPage) - pageSize

	if currentPage <= 3 {
		for i := 1; i <= pkg.Min(5, numPages); i++ {
			isCurrent := currentPage == i
			link := fmt.Sprintf("/ui/%s/?page=%v&page_size=%v", currentPageName, i, pageSize)
			pages = append(pages, Page{i, isCurrent, link})
		}
		return PaginationData{true, pages, offset, currentPageName}
	} else if currentPage >= numPages-2 {
		for i := pkg.Max(1, numPages-4); i <= numPages; i++ {
			isCurrent := currentPage == i
			link := fmt.Sprintf("/ui/%s/?page=%v&page_size=%v", currentPageName, i, pageSize)
			pages = append(pages, Page{i, isCurrent, link})
		}
		return PaginationData{true, pages, offset, currentPageName}
	} else {
		for i := currentPage - 2; i <= currentPage+2; i++ {
			isCurrent := currentPage == i
			link := fmt.Sprintf("/ui/%s/?page=%v&page_size=%v", currentPageName, i, pageSize)
			pages = append(pages, Page{i, isCurrent, link})
		}
		return PaginationData{true, pages, offset, currentPageName}
	}
}

func (app *application) homeReRoute(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/transfers", http.StatusSeeOther)
}
