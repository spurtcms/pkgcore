// Package Auth
package auth

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Authorization struct {
	DB     *gorm.DB
	Token  string
	Secret string
}
/*this struct holds dbconnection ,token*/
type Role struct {
	Auth Authorization
}
/*this struct holds dbconnection ,token*/
type PermissionAu struct {
	Auth Authorization
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

	db.Exec(`insert into tbl_roles('id','name','slug,'description','is_active','created_by','created_on') values(1,'Admin','admin','Has the full administration power',1,1,'2023-07-25 05:50:14')`)

	// db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS tbl_role_permisison_unique
	// ON public.tbl_role_permissions USING btree
	// (role_id ASC NULLS LAST, permission_id ASC NULLS LAST)
	// TABLESPACE pg_default;`)
}

type Action string

const ( //for permission check
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

		return 0, 0, errors.New("invalid token")
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
func Checklogin(Lc LoginCheck, db *gorm.DB, secretkey string) (string, int, error) {

	username := Lc.Username

	password := Lc.Password

	var user TblUser

	if err := db.Table("tbl_users").Where("username = ? and is_deleted=0", username).First(&user).Error; err != nil {

		return "", 0, err

	}

	passerr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passerr != nil || passerr == bcrypt.ErrMismatchedHashAndPassword {

		return "", 0, errors.New("invalid password")

	}

	token, err := CreateToken(user.Id, user.RoleId, secretkey)

	if err != nil {

		return "", 0, err
	}

	return token, user.Id, nil
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

		rolecounts, _ := AS.GetAllRoles(&roleco, 0, 0, filter, a.Auth.DB)

		return role, rolecounts, nil
	}

	return []TblRole{}, 0, errors.New("not authorized")
}

// get role by id
func (a Role) GetRoleById(roleid int) (tblrole TblRole, err error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return TblRole{}, checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		var role TblRole

		AS.GetRoleById(&role, roleid, a.Auth.DB)

		return role, nil
	}

	return TblRole{}, errors.New("not authorized")
}

// create role
func (a Role) CreateRole(rolec RoleCreation) (TblRole, error) {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return TblRole{}, checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if rolec.Name == "" {

			return TblRole{}, errors.New("empty value")
		}

		var role TblRole

		role.Name = rolec.Name

		role.Description = rolec.Description

		role.Slug = strings.ToLower(role.Name)

		role.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		role.CreatedBy = userid

		err := AS.RoleCreate(&role, a.Auth.DB)

		if err != nil {

			return TblRole{}, err
		}

		return role, nil
	}

	return TblRole{}, errors.New("not authorized")
}

// update role
func (a Role) UpdateRole(rolec RoleCreation, roleid int) (role TblRole, err error) {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return TblRole{}, checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if rolec.Name == "" {

			return TblRole{}, errors.New("empty value")
		}

		var role TblRole

		role.Id = roleid

		role.Name = rolec.Name

		role.Description = rolec.Description

		role.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		role.ModifiedBy = userid

		err1 := AS.RoleUpdate(&role, a.Auth.DB)

		if err1 != nil {

			return TblRole{}, err1
		}

		return role, nil

	}

	return TblRole{}, errors.New("not authorized")
}

// delete role
func (a Role) DeleteRole(roleid int) (bool, error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if roleid <= 0 {

			return false, errors.New("invalid role id cannot delete")
		}

		var role TblRole

		err1 := AS.RoleDelete(&role, roleid, a.Auth.DB)

		var permissions []TblRolePermission

		AS.DeleteRolePermissionById(&permissions, roleid, a.Auth.DB)

		if err1 != nil {

			return false, err1
		}

		return true, nil

	}
	return false, errors.New("not authorized")
}

// change role status 0-inactive, 1-active
func (a Role) RoleStatus(roleid int, status int) (err error) {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, _ := a.Auth.IsGranted("Roles", CRUD)

	if check {

		if roleid <= 0 {

			return errors.New("invalid role id cannot change the status")
		}

		var rolestatus TblRole

		rolestatus.ModifiedBy = userid

		rolestatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err1 := AS.RoleIsActive(&rolestatus, roleid, status, a.Auth.DB)

		if err1 != nil {

			return err1
		}

		return nil

	}
	return errors.New("not authorized")
}

/*Check Role Already Exists*/
func (a Role) CheckRoleAlreadyExists(roleid int, rolename string) (bool, error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var role TblRole

	err1 := AS.CheckRoleExists(&role, roleid, rolename, a.Auth.DB)

	if err1 != nil {

		return false, err1
	}

	if role.Id == 0 {

		return false, nil
	}

	return true, nil
}

/**/
func (a Role) GetAllRoleData() (roles []TblRole, err error) {

	var role []TblRole

	rerr := AS.GetRolesData(&role, a.Auth.DB)

	if rerr != nil {
		return []TblRole{}, rerr
	}

	return role, nil
}

// create permission
func (a PermissionAu) CreatePermission(Perm MultiPermissin) error {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Auth.IsGranted("Roles", CRUD)

	if err != nil {

		return err
	}

	if check {

		var createrolepermission []TblRolePermission

		for _, roleperm := range Perm.Ids {

			var createmod TblRolePermission

			createmod.PermissionId = roleperm

			createmod.RoleId = Perm.RoleId

			createmod.CreatedBy = userid

			createmod.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			createrolepermission = append(createrolepermission, createmod)

		}

		if len(createrolepermission) != 0 {

			AS.CreateRolePermission(&createrolepermission, a.Auth.DB)

		}

	}

	return errors.New("not authorized")

}

// update permission
func (a PermissionAu) CreateUpdatePermission(Perm MultiPermissin) error {

	userid, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Auth.IsGranted("Roles", CRUD)

	if err != nil {

		return err
	}

	if check {

		var checknotexist []TblRolePermission

		cnerr := AS.CheckPermissionIdNotExist(&checknotexist, Perm.RoleId, Perm.Ids, a.Auth.DB)

		if len(Perm.Ids) == 0 {

			AS.Deleterolepermission(&TblRolePermission{}, Perm.RoleId, a.Auth.DB)
		}

		if cnerr != nil {

			log.Println(cnerr)

		} else if len(checknotexist) != 0 {

			AS.DeleteRolePermissionById(&checknotexist, Perm.RoleId, a.Auth.DB)
		}

		var checkexist []TblRolePermission

		cerr := AS.CheckPermissionIdExist(&checkexist, Perm.RoleId, Perm.Ids, a.Auth.DB)

		if cerr != nil {

			log.Println(cerr)

		}

		var existid []int

		for _, exist := range checkexist {

			existid = append(existid, exist.PermissionId)

		}

		pid := Difference(Perm.Ids, existid)

		var createrolepermission []TblRolePermission

		for _, roleperm := range pid {

			var createmod TblRolePermission

			createmod.PermissionId = roleperm

			createmod.RoleId = Perm.RoleId

			createmod.CreatedBy = userid

			createmod.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			createrolepermission = append(createrolepermission, createmod)

		}

		if len(createrolepermission) != 0 {

			AS.CreateRolePermission(&createrolepermission, a.Auth.DB)

		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// permission List
func (a PermissionAu) PermissionListRoleId(limit, offset, roleid int, filter Filter) (Module []TblModule, count int64, err error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return []TblModule{}, 0, checkerr
	}

	check, err := a.Auth.IsGranted("Roles", CRUD)

	if err != nil {

		return []TblModule{}, 0, err
	}

	if check {

		var allmodule []TblModule

		var allmodules []TblModule

		var parentid []int //all parentid

		AS.GetAllParentModules1(&allmodule, a.Auth.DB)

		for _, val := range allmodule {

			parentid = append(parentid, val.Id)
		}

		var submod []TblModule

		AS.GetAllSubModules(&submod, parentid, a.Auth.DB)

		for _, val := range allmodule {

			if val.ModuleName == "Settings" {

				var newmod TblModule

				newmod.Id = val.Id

				newmod.Description = val.Description

				newmod.CreatedBy = val.CreatedBy

				newmod.ModuleName = val.ModuleName

				newmod.IsActive = val.IsActive

				newmod.IconPath = val.IconPath

				newmod.CreatedDate = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

				for _, sub := range submod {

					if sub.ParentId == val.Id {

						for _, getmod := range sub.TblModulePermission {

							if getmod.ModuleId == sub.Id {

								var modper TblModulePermission

								modper.Id = getmod.Id

								modper.Description = getmod.Description

								modper.DisplayName = getmod.DisplayName

								modper.ModuleName = getmod.ModuleName

								modper.RouteName = getmod.RouteName

								modper.CreatedBy = getmod.CreatedBy

								modper.Description = getmod.Description

								modper.TblRolePermission = getmod.TblRolePermission

								modper.CreatedDate = val.CreatedOn.Format("2006-01-02 15:04:05")

								modper.FullAccessPermission = getmod.FullAccessPermission

								newmod.TblModulePermission = append(newmod.TblModulePermission, modper)
							}

						}
					}

				}

				allmodules = append(allmodules, newmod)

			} else if val.ModuleName == "Spaces" {

				var newmod TblModule

				newmod.Id = val.Id

				newmod.Description = val.Description

				newmod.CreatedBy = val.CreatedBy

				newmod.ModuleName = val.ModuleName

				newmod.IsActive = val.IsActive

				newmod.IconPath = val.IconPath

				newmod.CreatedDate = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

				for _, sub := range submod {

					if sub.Id == val.Id {

						for _, getmod := range sub.TblModulePermission {

							if getmod.ModuleId == val.Id {

								var modper TblModulePermission

								modper.Id = getmod.Id

								modper.Description = getmod.Description

								modper.DisplayName = getmod.DisplayName

								modper.ModuleName = getmod.ModuleName

								modper.RouteName = getmod.RouteName

								modper.CreatedBy = getmod.CreatedBy

								modper.Description = getmod.Description

								modper.TblRolePermission = getmod.TblRolePermission

								modper.CreatedDate = val.CreatedOn.Format("2006-01-02 15:04:05")

								modper.FullAccessPermission = getmod.FullAccessPermission

								newmod.TblModulePermission = append(newmod.TblModulePermission, modper)
							}

						}
					}

				}

				allmodules = append(allmodules, newmod)

			} else {

				for _, sub := range submod {

					if sub.ParentId == val.Id {

						var newmod TblModule

						newmod.Id = sub.Id

						newmod.Description = sub.Description

						newmod.CreatedBy = sub.CreatedBy

						newmod.ModuleName = sub.ModuleName

						newmod.IsActive = sub.IsActive

						newmod.IconPath = sub.IconPath

						newmod.CreatedDate = sub.CreatedOn.Format("02 Jan 2006 03:04 PM")

						for _, getmod := range sub.TblModulePermission {

							if getmod.ModuleId == sub.Id {

								var modper TblModulePermission

								modper.Id = getmod.Id

								modper.Description = sub.Description

								modper.DisplayName = getmod.DisplayName

								modper.ModuleName = getmod.ModuleName

								modper.RouteName = getmod.RouteName

								modper.CreatedBy = getmod.CreatedBy

								modper.Description = getmod.Description

								modper.TblRolePermission = getmod.TblRolePermission

								modper.CreatedDate = val.CreatedOn.Format("2006-01-02 15:04:05")

								modper.FullAccessPermission = getmod.FullAccessPermission

								newmod.TblModulePermission = append(newmod.TblModulePermission, modper)
							}

						}

						allmodules = append(allmodules, newmod)

					}
				}

			}

		}

		var allmodul []TblModule

		Totalcount := AS.GetAllModules(&allmodul, 0, 0, roleid, filter, a.Auth.DB)

		return allmodules, Totalcount, nil

	} else {

		return []TblModule{}, 0, errors.New("not authorized")
	}

}

// permission List
func (a PermissionAu) GetPermissionDetailsById(roleid int) (rolepermissionid []int, err error) {

	_, _, checkerr := VerifyToken(a.Auth.Token, a.Auth.Secret)

	if checkerr != nil {

		return []int{}, checkerr
	}

	check, err := a.Auth.IsGranted("Roles", CRUD)

	if err != nil {

		return []int{}, err
	}

	if check {

		var permissionid []int

		var roleper []TblRolePermission

		AS.GetPermissionId(&roleper, roleid, a.Auth.DB)

		for _, val := range roleper {

			permissionid = append(permissionid, val.PermissionId)

		}

		return permissionid, nil

	} else {

		return []int{}, errors.New("not authorized")
	}

}

// Check User Permission
func (a Authorization) IsGranted(modulename string, permisison Action) (bool, error) {

	_, roleid, checkerr := VerifyToken(a.Token, a.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	if roleid != 1 { //if not an admin user

		var modid int

		var module TblModule

		var modpermissions TblModulePermission

		if err := a.DB.Table("tbl_modules").Where("module_name=?", modulename).Find(&module).Error; err != nil {

			return false, err
		}

		if err1 := a.DB.Table("tbl_module_permissions").Where("display_name=?", modulename).Find(&modpermissions).Error; err1 != nil {

			return false, err1
		}

		if module.Id != 0 {

			modid = module.Id

		} else {

			modid = modpermissions.Id
		}

		var modulepermission []TblModulePermission

		if permisison == "CRUD" {

			if err := a.DB.Table("tbl_module_permissions").Where("id=? and (full_access_permission=1 or display_name='View' or display_name='Update' or  display_name='Create' or display_name='Delete')", modid).Find(&modulepermission).Error; err != nil {

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
