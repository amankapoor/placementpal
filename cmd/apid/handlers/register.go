package handlers

import (
	"context"
	"net/http"

	"github.com/amankapoor/goth"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/amankapoor/placementpal/internal/user"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

type registerForm struct {
	Batch           string     `valid:"required"`
	Degree          string     `valid:"required"`
	College         string     `valid:"required"`
	CollegeLocation string     `valid:"required"`
	GU              *goth.User `valid:"required"`
}

func RegisterGet(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {

	// becuase of template, direct access to register page yields error

	s := registerForm{
		GU: ShareGothUser,
	}

	web.Respond(ctx, w, "register.html", s.GU.FirstName, 200)
	return nil
}

func (db *Database) RegisterPost(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {

	err := r.ParseForm()
	if err != nil {
		return errors.Wrap(err, "unable to parse register form")
	}

	s := registerForm{
		Batch:           r.FormValue("batch"),
		Degree:          r.FormValue("degree"),
		College:         r.FormValue("college"),
		CollegeLocation: r.FormValue("collegeLocation"),
		GU:              ShareGothUser,
	}

	//fmt.Println("REGISTERING...", s.Batch, s.Degree)

	res, err := govalidator.ValidateStruct(s)
	if err != nil || res == false {
		return errors.Wrap(err, "invalid inputs by user")
	}
	//fmt.Println(res)

	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer reqDB.MGOClose()

	cu := user.CreateUser{
		Batch:           s.Batch,
		Degree:          s.Degree,
		College:         s.College,
		CollegeLocation: s.CollegeLocation,
		Email:           s.GU.Email,
		FirstName:       s.GU.FirstName,
		LastName:        s.GU.LastName,
		GoogleID:        s.GU.UserID,
		AvatarURL:       s.GU.AvatarURL,
		UserType:        user.Student,
	}

	//fmt.Printf("CREATING THIS USER %+v", cu)

	// create new db entry ands issue session
	u, err := user.Create(ctx, reqDB, cu)
	if err == nil {
		s := sessionStruct{
			Batch:     u.Batch,
			Degree:    u.Degree,
			FirstName: u.FirstName,
			GoogleID:  u.GoogleID,
		}
		issueSession(s, w, r)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		return errors.Wrap(err, "<<unable to create new user>>")

	}
	return nil

	//fmt.Printf("CREATED USER IS: %+v", nu)

}
