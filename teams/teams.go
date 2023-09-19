package teams

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/spurtcms/spurtcms-core/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type TeamAuth struct {
	Authority *auth.Authority
}

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
func (a *TeamAuth) ListUser(limit, offset int, filter Filters) (tbluser []TblUser, totoaluser int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblUser{}, 0, checkerr
	}

	check, err := a.Authority.IsGranted("Users", auth.Read)

	if err != nil {

		return []TblUser{}, 0, err
	}

	if check {

		var Total_users int64

		var users []TblUser

		query := a.Authority.DB.Table("tbl_users").Select("tbl_users.id,tbl_users.uuid,tbl_users.role_id,tbl_users.first_name,tbl_users.last_name,tbl_users.email,tbl_users.password,tbl_users.username,tbl_users.mobile_no,tbl_users.profile_image,tbl_users.profile_image_path,tbl_users.created_on,tbl_users.created_by,tbl_users.modified_on,tbl_users.modified_by,tbl_users.is_active,tbl_users.is_deleted,tbl_users.deleted_on,tbl_users.deleted_by,tbl_users.data_access,tbl_roles.name as role_name").
			Joins("inner join tbl_roles on tbl_users.role_id = tbl_roles.id").Where("tbl_users.is_deleted=?", 0)

		if filter.Keyword != "" {

			query = query.Where("(LOWER(TRIM(tbl_users.first_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
				Or("LOWER(TRIM(tbl_users.last_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
				Or("LOWER(TRIM(tbl_roles.name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
				Or("LOWER(TRIM(tbl_users.username)) ILIKE LOWER(TRIM(?)))", "%"+filter.Keyword+"%")

		}

		if limit != 0 {

			query.Offset(offset).Limit(limit).Order("id desc").Find(&users)

			return users, 0, nil

		}

		query.Find(&users).Count(&Total_users)

		return []TblUser{}, Total_users, nil

	}

	return []TblUser{}, 0, errors.New("not authorized")
}

/*User Creation*/
func (a *TeamAuth) CreateUser(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if c.PostFormValue("mem_role") == "" || c.PostFormValue("mem_fname") == "" || c.PostFormValue("mem_lname") == "" || c.PostFormValue("mem_email") == "" || c.PostFormValue("mem_usrname") == "" || c.PostFormValue("mem_mob") == "" || c.PostFormValue("mem_pass") == "" {

		return errors.New("given some values is empty")
	}

	check, err := a.Authority.IsGranted("User", auth.Create)

	if err != nil {

		return err
	}

	if check {

		password := c.PostFormValue("mem_pass")

		uvuid := (uuid.New()).String()

		hash_pass := hashingPassword(password)

		var user TblUser

		user.Uuid = uvuid

		user.RoleId, _ = strconv.Atoi(c.PostFormValue("mem_role"))

		user.FirstName = c.PostFormValue("mem_fname")

		user.LastName = c.PostFormValue("mem_lname")

		user.Email = c.PostFormValue("mem_email")

		user.Username = c.PostFormValue("mem_usrname")

		user.Password = hash_pass

		user.MobileNo = c.PostFormValue("mem_mob")

		user.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))

		user.DataAccess, _ = strconv.Atoi(c.PostFormValue("mem_data_access"))

		user.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		user.CreatedBy = userid

		if err := a.Authority.DB.Create(&user).Error; err != nil {

			return err

		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// Update User
func (a *TeamAuth) UpdateUser(c *http.Request) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if c.PostFormValue("mem_role") == "" || c.PostFormValue("mem_fname") == "" || c.PostFormValue("mem_lname") == "" || c.PostFormValue("mem_email") == "" || c.PostFormValue("mem_usrname") == "" || c.PostFormValue("mem_mob") == "" {

		return errors.New("given some values is empty")
	}

	user_id, _ := strconv.Atoi(c.PostFormValue("mem_id"))

	password := c.PostFormValue("mem_pass")

	var user TblUser

	if password != "" {

		hash_pass := hashingPassword(password)

		user.Password = hash_pass
	}

	user.Id = user_id

	user.RoleId, _ = strconv.Atoi(c.PostFormValue("mem_role"))

	user.FirstName = c.PostFormValue("mem_fname")

	user.LastName = c.PostFormValue("mem_lname")

	user.Email = c.PostFormValue("mem_email")

	user.Username = c.PostFormValue("mem_usrname")

	user.MobileNo = c.PostFormValue("mem_mob")

	user.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	user.ModifiedBy = user_id

	if c.PostFormValue("mem_activestat") != "" {

		user.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))
	}

	if c.PostFormValue("mem_data_access") != "" {

		user.DataAccess, _ = strconv.Atoi(c.PostFormValue("mem_data_access"))
	}

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

	return nil
}

// Delete User
func (a *TeamAuth) DeleteUser(id int) error {

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

		if err := a.Authority.DB.Model(&user).Where("id=?", id).Updates(TblUser{IsDeleted: user.IsDeleted, DeletedOn: user.DeletedOn, DeletedBy: user.DeletedBy}).Error; err != nil {

			return err

		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// check email
func (a *TeamAuth) CheckEmail(Email string, userid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	if userid == 0 {
		if err := a.Authority.DB.Table("tbl_users").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and is_deleted = 0 ", Email).First(&user).Error; err != nil {

			return false, err
		}
	} else {
		if err := a.Authority.DB.Table("tbl_users").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and id not in(?) and is_deleted= 0 ", Email, userid).First(&user).Error; err != nil {

			return false, err
		}
	}

	return true, nil
}

// check mobile
func (a *TeamAuth) CheckNumber(mobile string, userid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	if userid == 0 {
		if err := a.Authority.DB.Table("tbl_users").Where("mobile_no = ? and is_deleted=0", mobile).First(&user).Error; err != nil {

			return false, err
		}
	} else {
		if err := a.Authority.DB.Table("tbl_users").Where("mobile_no = ? and id not in (?) and is_deleted=0", mobile, userid).First(&user).Error; err != nil {

			return false, err
		}

	}

	return true, nil
}

// check username
func (a *TeamAuth) CheckUsername(username string, userid int) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var user TblUser

	if userid == 0 {
		if err := a.Authority.DB.Table("tbl_users").Where("username = ? and is_deleted=0", username).First(&user).Error; err != nil {

			return false, err
		}
	} else {
		if err := a.Authority.DB.Table("tbl_users").Where("username = ? and id not in (?) and is_deleted=0", username, userid).First(&user).Error; err != nil {

			return false, err
		}

	}

	return true, nil
}
