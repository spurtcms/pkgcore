package memberaccess

import (
	"time"

	"github.com/spurtcms/spurtcms-core/member"
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
	DeletedOn         time.Time               `gorm:"DEFAULT:NULL"`
	DeletedBy         int                     `gorm:"DEFAULT:NULL"`
	Username          string                  `gorm:"column:username;<-:false"`
	Rolename          string                  `gorm:"column:name;<-:false"`
	MemberGroups      []member.TblMemberGroup `gorm:"-"`
	DateString        string                  `gorm:"-"`
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
	DeletedBy                int       `gorm:"DEFAULT:NULL"`
	ParentPageId             int       `gorm:"column:parent_id;<-:false"`
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
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
}

type Filter struct {
	Keyword string
}

type SubPage struct {
	Id       string `json:"id"`
	GroupId  string `json:"groupId"`
	ParentId string `json:"parentId"`
	SpaceId  string `json:"spaceId"`
}

type Page struct {
	Id      string `json:"id"`
	GroupId string `json:"groupId"`
	SpaceId string `json:"spaceId"`
}

type PageGroup struct {
	Id      string `json:"id"`
	SpaceId string `json:"spaceId"`
}

type MemberAccessControlRequired struct {
	Title          string
	Pages          []Page
	SubPage        []SubPage
	Group          []PageGroup
	SpacesIds      []int
	MemberGroupIds []int
}

type TblPage struct {
	Id          int `gorm:"primaryKey;auto_increment"`
	SpacesId    int
	PageGroupId int
	ParentId    int
	CreatedOn   time.Time
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	DeletedOn   time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy   int       `gorm:"DEFAULT:NULL"`
	IsDeleted   int       `gorm:"DEFAULT:0"`
}

func (at AccessType) GetSpaceByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.spaces_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=? and tbl_access_control_user_group.is_deleted=0", membergroupid).Group("spaces_id").Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil

}

func (at AccessType) GetPageByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.page_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=? and tbl_access_control_user_group.is_deleted=0", membergroupid).Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}

func (at AccessType) GetGroupByMemberId(tblaccess *[]TblAccessControlUserGroup, membergroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_pages.page_group_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id =tbl_access_control_user_group.id").Where("member_group_id=? and tbl_access_control_user_group.is_deleted=0", membergroupid).Find(&tblaccess).Error; err != nil {

		return err
	}
	return nil
}

/**/
func (at AccessType) CheckPageRestrict(page *[]TblAccessControlUserGroup, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Select("tbl_access_control_user_group.*,tbl_access_control_pages.page_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id = tbl_access_control_user_group.id").Where("page_id=? and tbl_access_control_pages.is_deleted=0", pageid).Find(&page).Error; err != nil {

		return err
	}

	return nil
}

// Get all content access list
func (at AccessType) GetContentAccessList(contentAccessList *[]TblAccessControl, limit, offset int, filter Filter, DB *gorm.DB) (list []TblAccessControl, count int64) {

	query := DB.Table("tbl_access_control").Select("tbl_access_control.*,tbl_users.username,tbl_roles.name").Joins("left join tbl_users on tbl_users.id = tbl_access_control.created_by").
		Joins("left join tbl_roles on tbl_roles.id = tbl_users.role_id").
		Where("tbl_access_control.is_deleted = 0 AND tbl_users.is_deleted = 0 AND tbl_roles.is_deleted = 0").Order("tbl_access_control.id DESC")

	if filter.Keyword != "" {

		query.Where("(LOWER(TRIM(tbl_access_control.access_control_name)) ILIKE LOWER(TRIM(?)))", "%"+filter.Keyword+"%")
	}

	// if q.DataAccess == 1 {

	// 	query = query.Where("tbl_access_control.created_by = ?", q.UserId)
	// }

	if limit != 0 {

		query.Offset(offset).Limit(limit).Find(&contentAccessList)

		return *contentAccessList, 0

	} else {

		var count int64

		var emptyslice []TblAccessControl

		query.Find(&contentAccessList).Count(&count)

		return emptyslice, count
	}
}

func (at AccessType) GetAccessGrantedMemberGroups(memberGroups *[]TblAccessControlUserGroup, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Where("is_deleted = 0 and access_control_id = ?", accessId).Find(&memberGroups).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetMemberGroupsByContentAccessMemId(memgrp *member.TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_groups").Where("is_deleted = 0 and id = ?", id).First(&memgrp).Error; err != nil {

		return err

	}

	return nil

}

/*Create Access*/
func (at AccessType) NewContentAccessEntry(contentAccess *TblAccessControl, DB *gorm.DB) (*TblAccessControl, error) {

	if err := DB.Table("tbl_access_control").Create(&contentAccess).Error; err != nil {

		return &TblAccessControl{}, err
	}

	return contentAccess, nil
}

func (at AccessType) GrantAccessToMemberGroups(memberGrpAccess *TblAccessControlUserGroup, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Create(&memberGrpAccess).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetMemberGrpByAccessControlId(memberGrpAccess *[]TblAccessControlUserGroup, content_access_id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Where("is_deleted = 0 and access_control_id = ?", content_access_id).Find(&memberGrpAccess).Error; err != nil {

		return err
	}

	return nil

}

func (at AccessType) InsertPageEntries(spg_access *TblAccessControlPages, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Create(&spg_access).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetPagesUnderSpaces(tblPagesData *[]TblPage, space_id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Where("is_deleted = 0 and spaces_id = ?", space_id).Find(&tblPagesData).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) UpdateContentAccessId(contentAccess *TblAccessControl, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control").Where("is_deleted = 0 and id = ?", contentAccess.Id).UpdateColumns(map[string]interface{}{"access_control_name": contentAccess.AccessControlName, "access_control_slug": contentAccess.AccessControlSlug, "modified_on": contentAccess.ModifiedOn, "modified_by": contentAccess.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil

}

func (at AccessType) CheckPresenceOfAccessGrantedMemberGroups(count *int64, mem_id, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Where("is_deleted = 0 and member_group_id = ? and access_control_id = ?", mem_id, accessId).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) UpdateContentAccessMemberGroup(accessmemgrp *TblAccessControlUserGroup, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Where("is_deleted = 0 and access_control_id = ? and member_group_id = ?", accessmemgrp.AccessControlId, accessmemgrp.MemberGroupId).UpdateColumns(map[string]interface{}{"modified_on": accessmemgrp.ModifiedOn, "modified_by": accessmemgrp.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) CheckPresenceOfPageInContentAccess(memid, pgid, groupid, spid int, count *int64, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and spaces_id = ? and page_group_id = ? and page_id = ?", memid, spid, groupid, pgid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) UpdatePagesInContentAccess(page_access *TblAccessControlPages, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and spaces_id = ? and page_group_id = ? and page_id = ?", page_access.AccessControlUserGroupId, page_access.SpacesId, page_access.PageGroupId, page_access.PageId).UpdateColumns(map[string]interface{}{"modified_on": page_access.ModifiedOn, "modified_by": page_access.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) RemoveMemberGroupsNotUnderContentAccessRights(memgrp_access *TblAccessControlUserGroup, memgrp_array []int, access_id int, DB *gorm.DB) error {

	if err := DB.Exec(`
		WITH updated_user_groups AS (
			UPDATE tbl_access_control_user_group
			SET is_deleted = (?),
			deleted_by = (?),
			deleted_on = (?)
			WHERE tbl_access_control_user_group.IS_DELETED =0 and tbl_access_control_user_group.access_control_id=? and tbl_access_control_user_group.member_group_id not in(?)
			RETURNING id
		)
		UPDATE tbl_access_control_pages
		SET is_deleted = (?),
		deleted_by = (?),
		deleted_on = (?)
		FROM updated_user_groups
		WHERE tbl_access_control_pages.access_control_user_group_id = (
			SELECT id
			FROM tbl_access_control_user_group
			WHERE tbl_access_control_user_group.id = updated_user_groups.id
		)`, memgrp_access.IsDeleted, memgrp_access.DeletedBy, memgrp_access.DeletedOn, access_id, memgrp_array, memgrp_access.IsDeleted, memgrp_access.DeletedBy, memgrp_access.DeletedOn).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) RemovePagesNotUnderContentAccess(pg_access *TblAccessControlPages, pageIds []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and page_id NOT IN ?", pg_access.AccessControlUserGroupId, pageIds).UpdateColumns(map[string]interface{}{"is_deleted": pg_access.IsDeleted, "deleted_on": pg_access.DeletedOn, "deleted_by": pg_access.DeletedBy}).Error; err != nil {

		return err
	}

	return nil
}

// Delete Access Control tbl

func (at AccessType) DeleteControl(accesscontrol *TblAccessControl, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control").Where("id = ?", id).UpdateColumns(map[string]interface{}{"deleted_by": accesscontrol.DeletedBy, "deleted_on": accesscontrol.DeletedOn, "is_deleted": accesscontrol.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

// Delete Access Control User Group tbl

func (at AccessType) DeleteInAccessUserGroup(accessusergrp *TblAccessControlUserGroup, Id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_group").Where("access_control_id = ?", Id).UpdateColumns(map[string]interface{}{"deleted_by": accessusergrp.DeletedBy, "deleted_on": accessusergrp.DeletedOn, "is_deleted": accessusergrp.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

// To Get Deleted id in access control user group tbl

func (at AccessType) GetDeleteIdInAccessUserGroup(controlaccessgrp *[]TblAccessControlUserGroup, Id int, DB *gorm.DB) (*[]TblAccessControlUserGroup, error) {

	if err := DB.Table("tbl_access_control_user_group").Where("access_control_id = ?", Id).Find(&controlaccessgrp).Error; err != nil {

		return &[]TblAccessControlUserGroup{}, err
	}

	return controlaccessgrp, nil
}

// Delete Access Control Pages tbl

func (at AccessType) DeleteAccessControlPages(pg_access *TblAccessControlPages, Id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Where("access_control_user_group_id IN ?", Id).UpdateColumns(map[string]interface{}{"deleted_by": pg_access.DeletedBy, "deleted_on": pg_access.DeletedOn, "is_deleted": pg_access.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetAccessGrantedMemberGroupsList(memgrps *[]int, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_groups").Select("tbl_member_groups.id").
		Joins("left join tbl_access_control_user_group on tbl_access_control_user_group.member_group_id =  tbl_member_groups.id and tbl_access_control_user_group.is_deleted = 0 ").
		Where("tbl_member_groups.is_deleted = 0 and tbl_access_control_user_group.access_control_id = ?", accessId).Find(&memgrps).Error; err != nil {

		return err

	}

	return nil

}

func (at AccessType) GetContentAccessByAccessId(accesscontrol *TblAccessControl, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control").Where("is_deleted = 0 and id = ?", id).First(&accesscontrol).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetPagesAndPageGroupsInContentAccess(contentAccessPages *[]TblAccessControlPages, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Select("tbl_access_control_pages.*,tbl_page.parent_id").Joins("left join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_group.is_deleted = 0").
		Joins("inner join tbl_access_control on tbl_access_control.id = tbl_access_control_user_group.access_control_id and tbl_access_control.is_deleted = 0").
		Joins("inner join tbl_page on tbl_page.id = tbl_access_control_pages.page_id and tbl_access_control_pages.is_deleted = 0").
		Where("tbl_access_control.id = ?", accessId).Find(&contentAccessPages).Error; err != nil {

		return err
	}

	return nil
}

func (at AccessType) GetPagesUnderPageGroup(pagesinPgg *[]TblPage, PageGroupId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Joins("inner join tbl_page_aliases on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page.is_deleted = 0 and tbl_page.page_group_id = ? and tbl_page_aliases.status = 'publish' and tbl_page_aliases.is_deleted = 0", PageGroupId).Find(&pagesinPgg).Error; err != nil {

		return err
	}

	return nil

}

func (at AccessType) GetContentAccessSpaces(spaceIds *[]int, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Select("distinct(tbl_access_control_pages.spaces_id)").
		Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_access_control on tbl_access_control.id = tbl_access_control_user_group.access_control_id").
		Where("tbl_access_control.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_control_user_group.is_deleted = 0 and tbl_access_control.id = ?", accessId).Find(&spaceIds).Error; err != nil {

		return err
	}
	return nil
}

func (at AccessType) GetcontentAccessPagesBySpaceId(ContentAccessPages *[]int, spid, accessid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_pages").Select("distinct(tbl_access_control_pages.page_id)").
		Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_access_control on tbl_access_control.id = tbl_access_control_user_group.access_control_id").
		Where("tbl_access_control_pages.is_deleted = 0 and tbl_access_control_pages.spaces_id = ? and tbl_access_control.id = ?", spid, accessid).Find(&ContentAccessPages).Error; err != nil {

		return err
	}

	return nil
}
