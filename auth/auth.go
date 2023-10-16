package auth

import (
	"errors"
	"fmt"
	"log"
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

type Role struct {
	Auth Authority
}

type Authstruct struct{}
var AS Authstruct

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

	db.Exec(`insert into tbl_roles('id','name','description','is_active','created_by','created_on') values(1,'admin','Has the full administration power',1,1,'2023-07-25 05:50:14')`)

	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS tbl_role_permisison_unique
    ON public.tbl_role_permissions USING btree
    (role_id ASC NULLS LAST, permission_id ASC NULLS LAST)
    TABLESPACE pg_default;`)
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

	usrid := Claims["user_id"]

	rolid := Claims["role_id"]

	return int(usrid.(float64)), int(rolid.(float64)), nil
}

// Check UserName Password
func Checklogin(c *http.Request, db *gorm.DB, secretkey string) (string, error) {

	username := c.PostFormValue("username")

	password := c.PostFormValue("pass")

	var user TblUser

	if err := db.Table("tbl_users").Where("username = ?", username).First(&user).Error; err != nil {

		return "", err

	}

	passerr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passerr != nil || passerr == bcrypt.ErrMismatchedHashAndPassword {

		return "", errors.New("invalid password")

	}

	token, err := CreateToken(user.Id, user.RoleId, secretkey)

	if err != nil {

		return "", err
	}

	return token, nil
}

// create role
func (a Role) RoleList(limit int, offset int, filter Filter) (roles []TblRole, rolecount int64, err error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return []TblRole{}, 0, checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if err != nil {

			return []TblRole{}, 0, err
		}

		var role []TblRole

		AS.GetAllRoles(&role, limit, offset, filter, a.Auth.DB)

		var roleco []TblRole

		rolecounts, _ := AS.GetAllRoles(&roleco, limit, offset, filter, a.Auth.DB)

		return role, rolecounts, nil
	}

	return []TblRole{}, 0, errors.New("not authorized")
}

// create role
func (a Role) CreateRole(rolec RoleCreation) error {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if rolec.Name == "" || rolec.Description == "" {

			return errors.New("empty value")
		}

		var role TblRole

		role.Name = rolec.Name

		role.Description = rolec.Description

		role.Slug = strings.ToLower(role.Name)

		role.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		role.CreatedBy = userid

		err := AS.RoleCreate(&role, a.Auth.DB)

		if err != nil {

			return err
		}

	}

	return errors.New("not authorized")
}

// update role
func (a Role) UpdateRole(rolec RoleCreation, roleid int) (err error) {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if rolec.Name == "" || rolec.Description == "" {

			return errors.New("empty value")
		}

		var role TblRole

		role.Id = roleid

		role.Name = rolec.Name

		role.Description = rolec.Description

		role.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		role.ModifiedBy = userid

		err1 := AS.RoleUpdate(&role, a.Auth.DB)

		if err1 != nil {

			return err1
		}

	}

	return errors.New("not authorized")
}

// delete role
func (a Role) DeleteRole(roleid int) (err error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if roleid <= 0 {

			return errors.New("invalid role id cannot delete")
		}

		var role TblRole

		err1 := AS.RoleDelete(&role, roleid, a.Auth.DB)

		if err != nil {

			return err1
		}

	}
	return errors.New("not authorized")
}

// create permission
func (a Authority) CreatePermission(Perm MultiPermissin) error {

	userid, _, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.IsGranted("Permissions", CRUD)

	if err != nil {

		return err
	}

	if check {

		var checknotexist []TblRolePermission

		cnerr := AS.CheckPermissionIdNotExist(&checknotexist, Perm.RoleId, Perm.Ids, a.DB)

		if cnerr != nil {

			log.Println(cnerr)

		} else if len(checknotexist) != 0 {

			AS.DeleteRolePermissionById(&checknotexist, Perm.RoleId, a.DB)
		}

		var checkexist []TblRolePermission

		cerr := AS.CheckPermissionIdExist(&checkexist, Perm.RoleId, Perm.Ids, a.DB)

		if cerr != nil {

			log.Println(cerr)

		} else {

			var existid []int

			for _, exist := range checkexist {

				existid = append(existid, exist.PermissionId)

			}

			pid := Difference([]int{}, existid)

			var createrolepermission []TblRolePermission

			for _, roleperm := range pid {

				var createmod TblRolePermission

				createmod.PermissionId = roleperm

				createmod.RoleId = Perm.RoleId

				createmod.CreatedBy = userid

				createmod.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

				createrolepermission = append(createrolepermission, createmod)

			}

			if len(createrolepermission) != 0 {

				AS.CreateRolePermission(&createrolepermission, a.DB)

			}

		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// Create permissionforrole
func (a Authority) AssignPermissionforRole() error {

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

	if err := a.DB.Debug().Table("tbl_modules").Where("module_name=?", modulename).Find(&module).Error; err != nil {

		return false, err
	}

	if err1 := a.DB.Debug().Table("tbl_module_permissions").Where("display_name=?", modulename).Find(&modpermissions).Error; err1 != nil {

		return false, err1
	}

	if module.Id != 0 {

		modid = module.Id

	} else {

		modid = modpermissions.Id
	}

	var modulepermission []TblModulePermission

	if permisison == "CRUD" {

		if err := a.DB.Debug().Table("tbl_module_permissions").Where("module_id=? and (full_access_permission=1 or display_name='View' or display_name='Update' or  display_name='Create' or display_name='Delete'", modid).Find(&modulepermission).Error; err != nil {

			return false, err
		}

	} else {

		if err := a.DB.Debug().Table("tbl_module_permissions").Where("module_id=? and display_name=?", modid, permisison).Find(&modulepermission).Error; err != nil {

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

// Set Difference: A - B
func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
