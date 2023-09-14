package auth

import "time"

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
