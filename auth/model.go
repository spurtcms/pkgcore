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
	CreatedDate         string `gorm:"<-:false"`
	DefaultModule       int
	ParentId            int
	IconPath            string
	TblModulePermission []TblModulePermission `gorm:"<-:false; foreignKey:ModuleId"`
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
	ModuleName           string    `gorm:"<-:false"`
	FullAccessPermission int
	ParentId             int
	AssignPermission     int
	BreadcrumbName       string
	TblRolePermission    []TblRolePermission `gorm:"<-:false; foreignKey:PermissionId"`
}

type TblRolePermission struct {
	Id           int `gorm:"primaryKey;auto_increment"`
	RoleId       int
	PermissionId int
	CreatedBy    int
	CreatedOn    time.Time
	CreatedDate  string `gorm:"<-:false"`
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
	RouteName            string    `gorm:"<-:false"`
	DisplayName          string    `gorm:"<-:false"`
	Description          string    `gorm:"-"`
	ModuleId             int       `gorm:"<-:false"`
	PermissionId         int       `gorm:"-"`
	FullAccessPermission int       `gorm:"<-:false"`
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
	CreatedDate string    `gorm:"<-:false"`
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
	RouteName    string    `gorm:"<-:false"`
	DisplayName  string    `gorm:"<-"`
	Description  string    `gorm:"-"`
	ModuleId     int       `gorm:"<-"`
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
	Action     []string //create,edit,update,delete

}

type MultiPermissin struct {
	RoleId      int
	Permissions []Permission
}

/*get all roles*/
func GetAllRoles(role *[]TblRole, limit, offset int, filter Filter, DB *gorm.DB) (rolecount int64, err error) {

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

// Roels Insert
func RoleCreate(role *TblRole, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Create(role).Error; err != nil {

		return err
	}

	return nil
}

// Delete the role data
func RoleDelete(role *TblRole, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_roles").Where("id = ?", id).Update("is_deleted", 1).Error; err != nil {

		return err

	}

	return nil
}

/**/
func RoleUpdate(role *TblRole, DB *gorm.DB) error {

	if err := DB.Model(&role).Where("id=?", role.Id).Updates(TblRole{Name: role.Name, Description: role.Description, Slug: role.Slug, IsActive: role.IsActive, IsDeleted: role.IsDeleted, ModifiedOn: role.ModifiedOn, ModifiedBy: role.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

func CheckPermissionIdNotExist(roleperm *[]TblRolePermission, roleid int, permissionid []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Where("role_id=? and permission_id not in(?)", roleid, permissionid).Find(&roleperm).Error; err != nil {

		return err

	}
	return nil
}

/*bulk creation*/
func CreateRolePermission(roleper *[]TblRolePermission, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Create(&roleper).Error; err != nil {

		return err

	}

	return nil
}

func CheckPermissionIdExist(roleperm *[]TblRolePermission, roleid int, permissionid []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_role_permissions").Where("role_id=? and permission_id in(?)", roleid, permissionid).Find(&roleperm).Error; err != nil {

		return err

	}
	return nil
}

func DeleteRolePermissionById(roleper *[]TblRolePermission, roleid int, DB *gorm.DB) error {

	if err := DB.Where("role_id=?", roleid).Delete(&roleper).Error; err != nil {

		return err

	}
	return nil
}
