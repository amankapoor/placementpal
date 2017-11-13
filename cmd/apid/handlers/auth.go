package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/amankapoor/goth"
	"github.com/amankapoor/goth/gothic"
	"github.com/amankapoor/placementpal/internal/platform/scs"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/amankapoor/placementpal/internal/user"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

var sessionManager *scs.Manager

var ShareGothUser *goth.User

// SV stands for session value
const SVGoogleID string = "googleID"
const SVFirstName string = "firstName"
const SVBatch string = "batch"
const SVDegree string = "degree"

type sessionStruct struct {
	GoogleID  string
	FirstName string
	Batch     string
	Degree    string
}

func init() {
	sessionManager = scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")
	sessionManager.Lifetime(time.Hour) // Set the maximum session lifetime to 1 hour.
	sessionManager.Persist(true)       // Persist the session after a user has closed their browser.
	sessionManager.Secure(false)       // Set the Secure flag on the session cookie.
	sessionManager.Path("/")
	sessionManager.HttpOnly(true)
	sessionManager.Name("S")

}

func (db *Database) LoginWithGoogle(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {
	// try to get the user without re-authenticating
	if gothuser, err := gothic.CompleteUserAuth(w, r, ps); err == nil {
		u, dbErr := db.userInDB(ctx, gothuser.UserID)
		if dbErr != nil {
			if dbErr == mgo.ErrNotFound {
				// issue a session on register page
				ShareGothUser = &gothuser
				http.Redirect(w, r, "/register", 302)
				return nil
			}

			return errors.Wrap(err, "<<some other db error occured>>")
		}

		s := sessionStruct{
			GoogleID:  u.GoogleID,
			FirstName: u.FirstName,
			Batch:     u.Batch,
			Degree:    u.Degree,
		}
		issueSession(s, w, r)
		http.Redirect(w, r, "/dashboard", 302)
	} else {
		gothic.BeginAuthHandler(w, r, ps)
	}
	return nil
}

func (db *Database) LoginWithGoogleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {

	if usr, err := gothic.CompleteUserAuth(w, r, ps); err == nil {
		// now get the user from gothUser.ID from db as the primary key
		// from the database and then pass user's degree and batch to issue session
		// then we will use data from session to give user's personalised results
		// on the dashboard
		//fmt.Println("GOTHUSER IS: ", gothuser)
		// reqDB, err := db.MasterDB.MGOCopy()
		// if err != nil {
		// 	return errors.Wrapf(web.ErrDBNotConfigured, "")
		// }
		// defer reqDB.MGOClose()
		// u, err := user.Retrieve(ctx, reqDB, usr.UserID)
		// //fmt.Println("RETRIEVED USER IS: %+v", &u)
		// if err != nil {
		// 	if err == mgo.ErrNotFound {
		// 		// now render a new page and ask user for degree and batch
		// 		ShareGothUser = &usr
		// 		http.Redirect(w, r, "/register", http.StatusSeeOther)
		// 		return nil
		// 	}

		// 	return errors.Wrap(err, "<<Unable to retrieve user from db while completing user auth>>")

		// }
		//issue session here with user's degree and batch and id details
		// and redirect to dashboard
		// issuing session from db data
		// s := sessionStruct{
		// 	GoogleID:  DBUser.User.GoogleID,
		// 	FirstName: DBUser.User.FirstName,
		// 	Degree:    DBUser.User.Degree,
		// 	Batch:     DBUser.User.Batch,
		// }

		// issueSession(s, w, r)

		u, dbErr := db.userInDB(ctx, usr.UserID)
		if dbErr != nil {
			if dbErr == mgo.ErrNotFound {
				// issue a session on register page
				ShareGothUser = &usr
				http.Redirect(w, r, "/register", 302)
				return nil
			}

			return errors.Wrap(err, "<<some other db error occured>>")
		}

		s := sessionStruct{
			GoogleID:  u.GoogleID,
			FirstName: u.FirstName,
			Batch:     u.Batch,
			Degree:    u.Degree,
		}
		issueSession(s, w, r)
		http.Redirect(w, r, "/dashboard", 302)
	} else {
		return errors.Wrap(err, "<<unable to complete user auth callback>>")
	}

	// func (db *Database) LoginWithGoogleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {

	// 	usr, err := gothic.CompleteUserAuth(w, r, ps)
	// 	if err != nil {
	// 		//fmt.Fprintln(w, err)
	// 		return errors.Wrap(err, "<<unable to complete user auth callback>>")
	// 	}

	// 	reqDB, err := db.MasterDB.MGOCopy()
	// 	if err != nil {
	// 		return errors.Wrapf(web.ErrDBNotConfigured, "")
	// 	}
	// 	defer reqDB.MGOClose()
	// 	u, err := user.Retrieve(ctx, reqDB, usr.UserID)
	// 	if err != nil {
	// 		if err == mgo.ErrNotFound {
	// 			// now render a new page and ask user for degree and batch
	// 			ShareGothUser = &usr
	// 			http.Redirect(w, r, "/register", http.StatusSeeOther)
	// 			return nil
	// 		}

	// 		return errors.Wrap(err, "<<Unable to retrieve user from db while completing user auth>>")

	// 	}

	// REPEAT SAME OF ABOVE FIND AND CREATE NEW STRUCT
	//insert session here
	return nil
}

func LogoutHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {
	gothic.Logout(w, r, ps)
	session := sessionManager.Load(r)
	err := session.Destroy(w)
	if err != nil {
		return errors.Wrap(err, "<<something bad happened, user should clear his cache>>")
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
	return nil
}

// func cacheInvalidator(next web.Handler) web.Handler {

// 	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {
// 		session := sessionManager.Load(r)
// 		err := session.Destroy(w)
// 		if err != nil {
// 			errors.Wrap(err, "<<unable to invalidate cache>>")
// 		}
// 		next(ctx, w, r, ps)
// 		return nil
// 	}
// 	fmt.Println("cache succesfully invalidated")
// 	return fn

// }

// issueSession issues a cookie session after successful Google login
func issueSession(s sessionStruct, w http.ResponseWriter, r *http.Request) error {

	//fmt.Printf("Printing the user retrieved during ISSUE SESSION: %v", user)

	// 2. Implement a success handler to issue some form of session
	session := sessionManager.Load(r)
	err := session.PutString(w, SVGoogleID, s.GoogleID)
	if err != nil {
		return errors.Wrap(err, "<<unable to PutString of google id while issuing session>>")
	}
	err = session.PutString(w, SVFirstName, s.FirstName)
	if err != nil {
		return errors.Wrap(err, "<<unable to PutString of firstname while issuing session>>")
	}

	err = session.PutString(w, SVBatch, s.Batch)
	if err != nil {
		return errors.Wrap(err, "<<unable to PutString of batch while issuing session>>")
	}

	err = session.PutString(w, SVDegree, s.Degree)
	if err != nil {
		return errors.Wrap(err, "<<unable to PutString of degree while issuing session>>")
	}

	return nil
}

// requireLogin redirects unauthenticated users to the login route.
func requireLogin(next web.Handler) web.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p web.Params) error {
		if !isAuthenticated(r) {
			// change this to wherever you have login page
			http.Redirect(w, r, "/", 302)
			return errors.Wrap(nil, "user is not authenticated so redirecting")
		}
		next(ctx, w, r, p)
		return nil
	}
	return fn
}

// onlybyadmin restricts access to admin only
func requireAdmin(next web.Handler) web.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p web.Params) error {
		if !isAuthenticated(r) {
			// change this to wherever you have login page
			http.Redirect(w, r, "/", 302)
			return errors.Wrap(nil, "user is not authenticated so redirecting")
		}
		if !isAdmin(r) {
			http.Error(w, "This page does not exist.", 404)
			return errors.Wrap(nil, "user is trying unauthorised access to admin zones")
		}
		next(ctx, w, r, p)
		return nil
	}
	return fn
}

// func (db *Database) ConditionalRegister(next web.Handler) web.Handler {
// 	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p web.Params) error {
// 		var err error
// 		_, err = db.userInDB(ctx, gid)
// 		if err == mgo.ErrNotFound {
// 			http.Redirect(w, r, "/register", 302)
// 			return errors.Wrap(nil, "user is not in database so redirecting")
// 		}
// 		next(ctx, w, r, p)
// 		return nil
// 	}
// 	return fn
// }

func (db *Database) userInDB(ctx context.Context, googleID string) (*user.User, error) {
	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return nil, err
	}
	defer reqDB.MGOClose()
	u, err := user.Retrieve(ctx, reqDB, googleID)
	//fmt.Println("RETRIEVED USER IS: %+v", &u)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, err
		}

		return nil, errors.Wrap(err, "<<Unable to retrieve user from db while completing user auth>>")
	}

	if u.GoogleID == "" {
		return nil, errors.Wrap(nil, "<<retrieved user from db but returned google id is of zero length>>")
	}

	return u, nil
}

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	session := sessionManager.Load(req)
	gID, err := session.GetString(SVGoogleID)
	if err != nil || gID == "" {
		return false
	}
	return true
}

// isAdmin returns true is the signed user is an admin
func isAdmin(req *http.Request) bool {
	session := sessionManager.Load(req)
	gID, err := session.GetString(SVGoogleID)
	if err != nil || gID == "" {
		return false
	}
	if gID != "116042618152348131594" {
		return false
	}
	return true
}

//////////////
// const (
// 	sessionName    = "ppal"
// 	sessionSecret  = "example cookie signing secret"
// 	sessionUserKey = "googleID"
// )

// // sessionStore encodes and decodes session data stored in signed cookies
// var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// var oauth2Config = &oauth2.Config{
// 	ClientID:     "1030362162437-qjd6mi9ld2gd0v2kiml6orvd590vpfjm.apps.googleusercontent.com",
// 	ClientSecret: "rhEICw3dqRBgGDDcRXDirESm",
// 	RedirectURL:  "http://localhost:8080/google/callback",
// 	Endpoint:     googleOAuth2.Endpoint,
// 	Scopes:       []string{"profile", "email"},
// }

// // New returns a new ServeMux with app routes.
// func New() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	//	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
// 	//mux.HandleFunc("/logout", logoutHandler)

// 	// state param cookies require HTTPS by default; disable for localhost development
// 	//stateConfig := gologin.DefaultCookieConfig
// 	//mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
// 	//mux.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)))
// 	return mux
// }

// // issueSession issues a cookie session after successful Google login
// func issueSession() web.Handler {
// 	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request, _ web.Params) error {

// 		googleUser, err := google.UserFromContext(ctx)
// 		fmt.Printf("Printing the user retrieved: %v", googleUser)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return nil
// 		}
// 		// 2. Implement a success handler to issue some form of session
// 		session := sessionStore.New(sessionName)
// 		session.Values[sessionUserKey] = googleUser.Id
// 		session.Save(w)
// 		http.Redirect(w, req, "/profile", http.StatusFound)
// 		return nil
// 	}
// 	return fn
// }

// // welcomeHandler shows a welcome message and login button.
// func WelcomeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {

// 	fmt.Fprintf(w, string("welcome page"))
// 	return nil
// }

// // profileHandler shows protected user content.
// func ProfileHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
// 	fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
// 	return nil
// }

// // // logoutHandler destroys the session on POSTs and redirects to home.
// // func LogoutHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {

// // 	if r.Method == "POST" {
// // 		sessionStore.Destroy(w, sessionName)
// // 	}
// // 	http.Redirect(w, r, "/", http.StatusFound)
// // 	return nil
// // }

// // requireLogin redirects unauthenticated users to the login route.
// func requireLogin(next web.Handler) web.Handler {
// 	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p web.Params) error {
// 		if !isAuthenticated(r) {
// 			http.Redirect(w, r, "/", http.StatusFound)
// 			return errors.Wrap(nil, "user is not authentication so redirecting")
// 		}
// 		next(ctx, w, r, p)
// 		return nil
// 	}
// 	return fn
// }

// // isAuthenticated returns true if the user has a signed session cookie.
// func isAuthenticated(req *http.Request) bool {
// 	if _, err := sessionStore.Get(req, sessionName); err == nil {
// 		return true
// 	}
// 	return false
// }

// // main creates and starts a Server listening.
// func main() {
// 	const address = "localhost:8080"

// 	log.Printf("Starting Server listening on %s\n", address)
// 	err := http.ListenAndServe(address, New())
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
