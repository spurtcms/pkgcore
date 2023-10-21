package memberaccess

import (
	"time"

	"gorm.io/gorm"
)

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
	SpacesId        int       `gorm:"-:migration;<-:false"`
	PageId          int       `gorm:"-:migration;<-:false"`
	PageGroupId     int       `gorm:"-:migration;<-:false"`
}

func (at AccessType) GetSpaceByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.spaces_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Group("spaces_id").Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil

}

func (at AccessType) GetPageByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.page_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}

func (at AccessType) GetGroupByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.page_group_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=?", membergroupid).Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}
