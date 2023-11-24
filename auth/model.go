package auth

import (
	"time"

	"gorm.io/gorm"
)

type TblModule struct {
	Id                  int `gorm:"primaryKey;auto_increment"`
	ModuleName          string
	IsActive            int
	CreatedBy           int
	CreatedOn           time.Time
	CreatedDate         string `gorm:"-:migration;<-:false"`
	DefaultModule       int
	ParentId            int
	IconPath            string
	TblModulePermission []TblModulePermission `gorm:"-:migration;<-:false; foreignKey:ModuleId"`
	Description         string
}

type TblModulePermission struct {
	Id                   int `gorm:"primaryKey;auto_increment"`
	RouteName            string
	DisplayName          string
	Description          string
	ModuleId             int
	CreatedBy            int
	CreatedOn            time.Time
	CreatedDate          string    `gorm:"-"`
	ModifiedBy           int       `gorm:"DEFAULT:NULL"`
	ModifiedOn           time.Time `gorm:"DEFAULT:NULL"`
	ModuleName           string    `gorm:"-:migration;<-:false"`
	FullAccessPermission int
	ParentId             int
	AssignPermission     int
	BreadcrumbName       string
	TblRolePermission    []TblRolePermission `gorm:"-:migration;<-:false; foreignKey:PermissionId"`
}

type TblRolePermission struct {
	Id           int `gorm:"primaryKey;auto_increment"`
	RoleId       int
	PermissionId int
	CreatedBy    int
	CreatedOn    time.Time
	CreatedDate  string `gorm:"-:migration;<-:false"`
}

type TblUser struct {
	Id                   int `gorm:"primaryKey;auto_increment"`
	Uuid                 string
	FirstName            string
	LastName             string
	RoleId               int
	Email                string
	Username             string
	Password             string
	MobileNo             string
	IsActive             int
	ProfileImage         string
	ProfileImagePath     string
	DataAccess           int
	CreatedOn            time.Time
	CreatedBy            int
	ModifiedOn           time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy           int       `gorm:"DEFAULT:NULL"`
	LastLogin            int
	IsDeleted            int
	DeletedOn            time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy            int       `gorm:"DEFAULT:NULL"`
	ModuleName           string    `gorm:"-"`
	RouteName            string    `gorm:"-:migration;<-:false"`
	DisplayName          string    `gorm:"-:migration;<-:false"`
	Description          string    `gorm:"-"`
	ModuleId             int       `gorm:"-:migration;<-:false"`
	PermissionId         int       `gorm:"-"`
	FullAccessPermission int       `gorm:"-:migration;<-:false"`
	Roles                []TblRole `gorm:"-"`
}

type TblRole struct {
	Id          int `gorm:"primaryKey;auto_increment"`
	Name        string
	Description string
	Slug        string
	IsActive    int
	IsDeleted   int
	CreatedOn   time.Time
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	CreatedDate string    `gorm:"-:migration;<-:false"`
	User        []TblUser `gorm:"-"`
}

type TblRoleUser struct {
	Id           int `gorm:"primaryKey;auto_increment"`
	RoleId       int
	UserId       int
	CreatedBy    int
	CreatedOn    time.Time
	ModifiedBy   int       `gorm:"DEFAULT:NULL"`
	ModifiedOn   time.Time `gorm:"DEFAULT:NULL"`
	ModuleName   string    `gorm:"-"`
	RouteName    string    `gorm:"-:migration;<-:false"`
	DisplayName  string    `gorm:"-:migration;<-"`
	Description  string    `gorm:"-"`
	ModuleId     int       `gorm:"-:migration;<-"`
	PermissionId int       `gorm:"-"`
}

type Filter struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

type Permission struct {
	ModuleName string
	Action     []string //create,View,update,delete

}

type MultiPermissin struct {
	RoleId      int
	Ids         []int
	// Permissions []Permission
}

type RoleCreation struct {
	Name        string
	Description string
}

type LoginCheck struct {
	Username string
	Password string
}

/*get all roles*/
func (as Authstruct) GetAllRoles(role *[]TblRole, limit, offset int, filter Filter, DB *gorm.DB) (rolecount int64, err error) {

	query := DB.Table("tbl_roles").Where("is_deleted = 0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&role)

	} else {

		query.Find(&role).Count(&rolecount)

		return rolecount, nil
	}

	return 0, nil
}

func (as Authstruct) GetRolesData(roles *[]TblRole, DB *gorm.DB) error {

	if err := DB.Where("is_deleted=? and is_active=1", 0).Find(&roles).Error; err != nil {

		return err

	}

	return nil
}

// Roels Insert
func (as Authstruct) RoleCreate(role *TblRole, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Create(role).Error; err != nil {

		return err
	}

	return nil
}

// Delete the role data
func (as Authstruct) RoleDelete(role *TblRole, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Where("id = ?", id).Update("is_deleted", 1).Error; err != nil {

		return err

	}

	return nil
}

/**/
func (as Authstruct) RoleUpdate(role *TblRole, DB *gorm.DB) error {

	if err := DB.Model(&role).Where("id=?", role.Id).Updates(TblRole{Name: role.Name, Description: role.Description, Slug: role.Slug, IsActive: role.IsActive, IsDeleted: role.IsDeleted, ModifiedOn: role.ModifiedOn, ModifiedBy: role.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (as Authstruct) CheckPermissionIdNotExist(roleperm *[]TblRolePermission, roleid int, permissionid []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Where("role_id=? and permission_id not in(?)", roleid, permissionid).Find(&roleperm).Error; err != nil {

		return err

	}
	return nil
}

/*bulk creation*/
func (as Authstruct) CreateRolePermission(roleper *[]TblRolePermission, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Create(&roleper).Error; err != nil {

		return err

	}

	return nil
}

func (as Authstruct) CheckPermissionIdExist(roleperm *[]TblRolePermission, roleid int, permissionid []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Where("role_id=? and permission_id in(?)", roleid, permissionid).Find(&roleperm).Error; err != nil {

		return err

	}
	return nil
}

func (as Authstruct) DeleteRolePermissionById(roleper *[]TblRolePermission, roleid int, DB *gorm.DB) error {

	if err := DB.Where("role_id=?", roleid).Delete(&roleper).Error; err != nil {

		return err

	}
	return nil
}

/*This is for assign permission*/
func (as Authstruct) GetAllModules(mod *[]TblModule, limit, offset, id int, filter Filter, DB *gorm.DB) (count int64) {

	query := DB.Table("tbl_modules").Where("parent_id!=0 or assign_permission=1").Preload("TblModulePermission", func(db *gorm.DB) *gorm.DB {
		return db.Where("assign_permission =0")
	}).Preload("TblModulePermission.TblRolePermission", func(db *gorm.DB) *gorm.DB {
		return db.Where("role_id = ?", id)
	})

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(module_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("id asc").Find(&mod)

	} else {

		query.Find(&mod).Count(&count)

		return count
	}

	return 0
}

/*Get role by id*/
func (as Authstruct) GetRoleById(role *TblRole, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Where("id=?", id).First(&role).Error; err != nil {

		return err

	}

	return nil
}

/*create json module permission*/
func (as Authstruct) CreateModulePermission(modpermission *TblModulePermission, DB *gorm.DB) (modper *TblModulePermission, err error) {

	if err := DB.Table("tbl_module_permissions").Create(&modpermission).Error; err != nil {

		return &TblModulePermission{}, err

	}

	return modpermission, nil
}

/*single creation*/
func (as Authstruct) CreateRolePermissionsingle(roleper *TblRolePermission, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Create(&roleper).Error; err != nil {

		return err

	}

	return nil
}

/*Get Id by RouteName*/
func (as Authstruct) GetIdByRouteName(tblmodper *TblModulePermission, id string, DB *gorm.DB) error {

	if err := DB.Table("tbl_module_permissions").Where("route_name=?", "/channel/entrylist/"+id).First(&tblmodper).Error; err != nil {

		return err
	}

	return nil
}

/*Delete Role Permission by id*/
func (as Authstruct) Deleterolepermission(TblRolePermission *TblRolePermission, id int, DB *gorm.DB) error {

	if err := DB.Where("permission_id=?", id).Delete(&TblRolePermission).Error; err != nil {

		return err
	}

	return nil
}

func (as Authstruct) DeleteModulePermissioninEntries(tblmodper *TblModulePermission, id string, DB *gorm.DB) error {

	if err := DB.Where("route_name=?", "/channel/entrylist/"+id).Delete(&tblmodper).Error; err != nil {

		return err
	}

	return nil
}

/*update channel entry permission*/
func (as Authstruct) UpdateChannelNameInEntries(modper *TblModulePermission, DB *gorm.DB) error {

	if err := DB.Table("tbl_module_permissions").Where("route_name=?", modper.RouteName).UpdateColumns(map[string]interface{}{
		"display_name": modper.DisplayName}).Error; err != nil {

		return err
	}

	return nil
}

/*Check role*/
func (as Authstruct) CheckRoleExists(role *TblRole, id int, name string, DB *gorm.DB) error {

	if id == 0 {
		if err := DB.Model(TblRole{}).Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and is_deleted = 0 ", name).First(&role).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Model(TblRole{}).Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and id not in(?) and is_deleted= 0 ", name, id).First(&role).Error; err != nil {

			return err
		}
	}
	return nil

}

/*update role isactive*/
func (as Authstruct) RoleIsActive(role *TblRole, id, val int, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_active": val, "modified_by": role.ModifiedBy, "modified_on": role.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Get PermissionId By RoleId*/
func (as Authstruct) GetPermissionId(perm *[]TblRolePermission, roleid int, DB *gorm.DB) error {

	if err := DB.Model(TblRolePermission{}).Where("role_id=?", roleid).Find(&perm).Error; err != nil {

		return err
	}

	return nil
}
