package memberaccess

import "time"

type TblAccessControl struct {
	Id                int `gorm:"primaryKey;auto_increment"`
	AccessControlName string
	AccessControlSlug string
	CreatedOn         time.Time
	CreatedBy         int
	ModifiedOn        time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy        int       `gorm:"DEFAULT:NULL"`
	IsDeleted         int
	DeletedOn         time.Time `gorm:"DEFAULT:NULL"`
}

type TblAccessControlPages struct {
	Id                       int `gorm:"primaryKey;auto_increment"`
	AccessControlUserGroupId int
	SpacesId                 int
	PageGroupId              int
	PageId                   int
	CreatedOn                time.Time
	CreatedBy                int
	ModifiedOn               time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy               int       `gorm:"DEFAULT:NULL"`
	IsDeleted                int
	DeletedOn                time.Time `gorm:"DEFAULT:NULL"`
}

type TblAccessControlUserGroup struct {
	Id              int `gorm:"primaryKey;auto_increment"`
	AccessControlId int
	MemberGroupId   int
	CreatedOn       time.Time
	CreatedBy       int
	ModifiedOn      time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	IsDeleted       int
	DeletedOn       time.Time `gorm:"DEFAULT:NULL"`
	SpacesId        []int     `gorm:"<-:false"`
	PagesId         []int     `gorm:"<-:false"`
	GroupsId         []int     `gorm:"<-:false"`
}

func GetSpaceByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int) error {

	if err := Access.Authority.DB.Table("tbl_access_control_user_group").Select("tbl_access_control_user_group.*tbl_access_control_pages.spaces_id as SpacesId").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Group("space_id").Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil

}

func GetPageByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int) error {

	if err := Access.Authority.DB.Table("tbl_access_control_user_group").Select("tbl_access_control_user_group.*tbl_access_control_pages.page_id as PagesId").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Group("member_group_id").Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}

func GetGroupByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int) error {

	if err := Access.Authority.DB.Table("tbl_access_control_user_group").Select("tbl_access_control_user_group.*tbl_access_control_pages.page_group_id as GroupsId").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Group("member_group_id").Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}