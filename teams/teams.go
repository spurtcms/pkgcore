package teams

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spurtcms/spurtcms-core/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type TeamAuth struct {
	Authority *auth.Authorization
}

type Team struct{}

var TM Team

func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(TblUser{})

	db.Exec(`CREATE INDEX IF NOT EXISTS username_unique
		ON public.tbl_users USING btree
		(username COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default
		WHERE is_deleted = 0;`)

	db.Exec(`insert into tbl_users('id','role_id','first_name','email','username','password','mobile_no','is_active') values(1,1,'spurtcms','spurtcms@gmail.com','spurtcms','$2a$14$r67QLbDoS0yVUbOwbzHfOOY/8eDnI5ya/Vux5j6A6LN9BCJT37ZpW','9876543210',1)`)
}

func hashingPassword(pass string) string {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		panic(err)

	}

	return string(passbyte)
}

/*List*/
func (a TeamAuth) ListUser(limit, offset int, filter Filters) (tbluser []TblUser, totoaluser int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblUser{}, 0, checkerr
	}

	check, err := a.Authority.IsGranted("Users", auth.Read)

	if err != nil {

		return []TblUser{}, 0, err
	}

	if check {

		var users []TblUser

		flg := false

		UserList, _ := TM.GetUsersList(&users, offset, limit, filter, flg, a.Authority.DB)

		var userscoount []TblUser

		_, usercount := TM.GetUsersList(&userscoount, 0, 0, filter, flg, a.Authority.DB)

		return UserList, usercount, nil

	}

	return []TblUser{}, 0, errors.New("not authorized")
}

/*User Creation*/
func (a TeamAuth) CreateUser(teamcreate TeamCreate) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if teamcreate.RoleId == 0 || teamcreate.FirstName == "" || teamcreate.Email == "" || teamcreate.Username == "" || teamcreate.MobileNo == "" || teamcreate.Password == "" {

		return errors.New("given some values is empty")
	}

	check, err := a.Authority.IsGranted("User", auth.Create)

	if err != nil {

		return err
	}

	if check {

		password := teamcreate.Password

		uvuid := (uuid.New()).String()

		hash_pass := hashingPassword(password)

		var user TblUser

		user.Uuid = uvuid

		user.RoleId = teamcreate.RoleId

		user.FirstName = teamcreate.FirstName

		user.LastName = teamcreate.LastName

		user.Email = teamcreate.Email

		user.Username = teamcreate.Username

		user.Password = hash_pass

		user.MobileNo = teamcreate.MobileNo

		user.IsActive = teamcreate.IsActive

		user.DataAccess = teamcreate.DataAccess

		user.ProfileImage = teamcreate.ProfileImage

		user.ProfileImagePath = teamcreate.ProfileImagePath

		user.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		user.CreatedBy = userid

		err := TM.CreateUser(&user, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// Update User
func (a TeamAuth) UpdateUser(teamcreate TeamCreate, userid int) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if teamcreate.RoleId == 0 || teamcreate.FirstName == "" || teamcreate.Email == "" || teamcreate.Username == "" || teamcreate.MobileNo == "" {

		return errors.New("given some values is empty")
	}

	check, err := a.Authority.IsGranted("User", auth.Create)

	if err != nil {

		return err
	}

	if check {

		user_id := userid

		password := teamcreate.Password

		var user TblUser

		if password != "" {

			hash_pass := hashingPassword(password)

			user.Password = hash_pass
		}

		user.Id = user_id

		user.RoleId = teamcreate.RoleId

		user.FirstName = teamcreate.FirstName

		user.LastName = teamcreate.LastName

		user.Email = teamcreate.Email

		user.Username = teamcreate.Username

		user.MobileNo = teamcreate.MobileNo

		user.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		user.ModifiedBy = user_id

		user.IsActive = teamcreate.IsActive

		user.DataAccess = teamcreate.DataAccess

		user.ProfileImage = teamcreate.ProfileImage

		user.ProfileImagePath = teamcreate.ProfileImagePath

		query := a.Authority.DB.Table("tbl_users").Where("id=?", user.Id)

		if user.Password == "" {

			if user.Password == "" {

				query = query.Omit("password", "profile_image", "profile_image_path").UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "data_access": user.DataAccess})

			}

			if err := query.Error; err != nil {

				return err
			}

		} else {

			if err := query.UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "profile_image": user.ProfileImage, "profile_image_path": user.ProfileImagePath, "data_access": user.DataAccess, "password": user.Password}).Error; err != nil {

				return err
			}

		}

	}
	return nil
}

// Delete User
func (a TeamAuth) DeleteUser(id int) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("User", auth.Delete)

	if err != nil {

		return err
	}

	if check {

		var user TblUser

		user.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		user.DeletedBy = userid

		user.IsDeleted = 1

		user.Id = id

		err := TM.DeleteUser(&user, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// check email
func (a TeamAuth) CheckEmail(Email string, userid int) (users TblUser, checl bool, errr error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return TblUser{}, false, checkerr
	}

	var user TblUser

	err := TM.CheckEmail(&user, Email, userid, a.Authority.DB)

	if err != nil {

		return TblUser{}, false, err
	}

	return user, true, nil
}

// check mobile
func (a TeamAuth) CheckNumber(mobile string, userid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	err := TM.CheckNumber(&user, mobile, userid, a.Authority.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

// check username
func (a TeamAuth) CheckUsername(username string, userid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	err := TM.CheckUsername(&user, username, userid, a.Authority.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

/**/
func (a TeamAuth) GetUserDetails(userid int) (user TblUser, err error) {

	var users TblUser

	usrerr := TM.GetUserDetailsTeam(&users, userid, a.Authority.DB)

	if usrerr != nil {

		log.Println(usrerr)

		return TblUser{}, usrerr
	}

	return users, nil
}

func (a TeamAuth) ChangeYourPassword(password string) (success bool, err error) {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var tbluser TblUser

	tbluser.Id = userid

	hash_pass := hashingPassword(password)

	tbluser.Password = hash_pass

	tbluser.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	tbluser.ModifiedBy = 1

	cerr := TM.ChangePasswordById(&tbluser, a.Authority.DB)

	if cerr != nil {

		return false, cerr
	}

	return true, nil
}

// Check role already used or not
func (a TeamAuth) CheckRoleUsed(roleid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var tbluser TblUser

	err := TM.CheckRoleUsed(&tbluser, roleid, a.Authority.DB)

	if err != nil {

		return false, err
	}

	return true, nil

}

// myprofile update//
func (a TeamAuth) UpdateMyUser(teamcreate TeamCreate) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if teamcreate.FirstName == "" || teamcreate.Email == "" || teamcreate.Username == "" || teamcreate.MobileNo == "" {

		return errors.New("given some values is empty")
	}

	password := teamcreate.Password

	var user TblUser

	if password != "" {

		hash_pass := hashingPassword(password)

		user.Password = hash_pass
	}

	user.Id = userid

	// user.RoleId = teamcreate.RoleId

	user.FirstName = teamcreate.FirstName

	user.LastName = teamcreate.LastName

	user.Email = teamcreate.Email

	user.Username = teamcreate.Username

	user.MobileNo = teamcreate.MobileNo

	user.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	user.ModifiedBy = userid

	user.IsActive = teamcreate.IsActive

	user.DataAccess = teamcreate.DataAccess

	user.ProfileImage = teamcreate.ProfileImage

	user.ProfileImagePath = teamcreate.ProfileImagePath

	query := a.Authority.DB.Table("tbl_users").Where("id=?", user.Id)

	if user.Password == "" {

		if user.Password == "" {

			query = query.Omit("password", "profile_image", "profile_image_path").UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "data_access": user.DataAccess})

		}

		if err := query.Error; err != nil {

			return err
		}

	} else {

		if err := query.UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "profile_image": user.ProfileImage, "profile_image_path": user.ProfileImagePath, "data_access": user.DataAccess, "password": user.Password}).Error; err != nil {

			return err
		}

	}

	return nil
}

/*check new password with old password*/
/*return false it will not match to old password*/
/*return true it will match to old password*/
func (a TeamAuth) CheckPasswordwithOld(password string) (bool, error) {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	err := TM.GetUserDetailsTeam(&user, userid, a.Authority.DB)

	if err != nil {

		return false, err
	}

	passerr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passerr == bcrypt.ErrMismatchedHashAndPassword {

		return false, nil

	}

	return true, nil
}
