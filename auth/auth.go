package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type Authority struct {
	DB     *gorm.DB
	Token  string
	Secret string
}

type Option struct {
	DB     *gorm.DB
	Token  string
	Secret string
}

func MigrationTable(db *gorm.DB) {
	db.AutoMigrate(
		TblModule{},
		TblModulePermission{},
		TblRolePermission{},
		TblRole{},
	)

	db.Exec(`insert into tbl_roles('id','name','description','is_active','created_by','created_on') values(1,'admin','Has the full administration power',1,'2023-07-25 05:50:14')`);

	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS tbl_role_permisison_unique
    ON public.tbl_role_permissions USING btree
    (role_id ASC NULLS LAST, permission_id ASC NULLS LAST)
    TABLESPACE pg_default;`);
}

type Permission struct {
	ModuleName string
	Action     []string //create,edit,update,delete

}

type MultiPermissin struct {
	Permissions []Permission
}

type Action string

const (
	Create Action = "Create"

	Read Action = "View"

	Update Action = "Update"

	Delete Action = "Delete"

	CRUD Action = "CRUD"
)

// CreateToken creates a token with the given claims
func CreateToken(userid, roleid int, secretkey string) (string, error) {

	atClaims := jwt.MapClaims{}

	atClaims["user_id"] = userid

	atClaims["role_id"] = roleid

	atClaims["expiry_time"] = time.Now().Add(2 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	return token.SignedString([]byte(secretkey))
}

// verify token
func VerifyToken(token string, secret string) (userid, roleid int, err error) {
	Claims := jwt.MapClaims{}

	tkn, err := jwt.ParseWithClaims(token, Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(err)
			return 0, 0, errors.New("invalid token")
		}

		return 0, 0, errors.New(err.Error())
	}

	if !tkn.Valid {
		fmt.Println(tkn)
		return 0, 0, errors.New("invalid token")
	}

	// exptime := Claims["expiry_time"]

	// current_time := time.Now().Unix()

	// etime := int64(exptime.(float64))

	// if (current_time <= etime) && (current_time >= etime - 300) {

	// }

	usrid := Claims["user_id"].(int)

	rolid := Claims["role_id"].(int)

	return usrid, rolid, nil
}

// Check UserName Password
func (a Authority) Checklogin(c *http.Request, secretkey string) (string,error) {

	username := c.PostFormValue("username")

	password := c.PostFormValue("pass")

	var user TblUser

	if err := a.DB.Table("tbl_users").Where("username = ?", username).First(&user).Error; err != nil {

		return "",err

	}

	passerr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passerr != nil || passerr == bcrypt.ErrMismatchedHashAndPassword {

		return "",errors.New("invalid password")

	}

	token, err := CreateToken(user.Id, user.RoleId, secretkey)

	if err != nil{

		return "",err
	}

	return token,nil
}

// create role
func (a Authority) CreateRole(c *http.Request) (TblRole, error) {

	userid, _, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return TblRole{}, checkerr
	}

	var role TblRole

	role.Name = c.PostFormValue("name")

	role.Description = c.PostFormValue("description")

	role.Slug = strings.ToLower(role.Name)

	role.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	role.CreatedBy = userid

	if err := a.DB.Table("tbl_roles").Create(&role).Error; err != nil {

		return TblRole{}, err
	}

	return role, nil
}

// create permission
func (a Authority) CreatePermission(c *http.Request, Perm MultiPermissin) error {

	_, _, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.IsGranted("Permissions", CRUD)

	if err != nil {

		return err
	}

	if check {

		// for _, val := range Perm.Permissions {

		// 	var modperm TblModulePermission

		// }

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// Create permissionforrole
func (a Authority) AssignPermissionforRole(c *http.Request) error {

	userid, roleid, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.IsGranted("Permissions", CRUD)

	if err != nil {

		return err
	}

	if check {

		var createmod TblRolePermission

		createmod.PermissionId = 1

		createmod.RoleId = roleid

		createmod.CreatedBy = userid

		createmod.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.DB.Table("tbl_role_permissions").Create(&createmod).Error; err != nil {

			return err

		}
	} else {

		return errors.New("not authorized")
	}

	return nil
}

// Check User Permission
func (a Authority) IsGranted(modulename string, permisison Action) (bool, error) {

	_, roleid, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var modid int

	var module TblModule

	var modpermissions TblModulePermission

	if err := a.DB.Table("tbl_modules").Where("module_name=?", modulename).Find(&module).Error; err != nil {

		if err1 := a.DB.Table("tbl_modules_permissions").Where("display_name=?", modulename).Find(&modpermissions).Error; err != nil {

			return false, err1
		}

		return false, err
	}

	if module.Id != 0 {

		modid = module.Id

	} else {

		modid = modpermissions.Id
	}

	var modulepermission []TblModulePermission

	if permisison == "CRUD" {

		if err := a.DB.Table("tbl_module_permissions").Where("module_id=? and (full_access_permission=1 or display_name='View' or display_name='Update' or  display_name='Create' or display_name='Delete'", modid).Find(&modulepermission).Error; err != nil {

			return false, err
		}

	} else {

		if err := a.DB.Table("tbl_module_permissions").Where("module_id=? and display_name=?", modid, permisison).Find(&modulepermission).Error; err != nil {

			return false, err
		}

	}

	for _, val := range modulepermission {

		var rolecheck TblRolePermission

		if err := a.DB.Table("tbl_role_permissions").Where("permission_id=? and role_id=?", val.Id, roleid).First(&rolecheck).Error; err != nil {

			return false, err
		}

	}

	return true, nil

}
