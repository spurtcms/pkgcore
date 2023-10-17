package member

import (
	"time"

	"gorm.io/gorm"
)

type TblMember struct {
	Id               int `gorm:"primaryKey;auto_increment;"`
	Uuid             string
	FirstName        string
	LastName         string
	Email            string
	MobileNo         string
	IsActive         int
	ProfileImage     string
	ProfileImagePath string
	LastLogin        int
	IsDeleted        int
	DeletedOn        time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	CreatedOn        time.Time `gorm:"DEFAULT:NULL"`
	CreatedDate      string    `gorm:"-"`
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	MemberGroupId    int
	Group            []TblMemberGroup `gorm:"-"`
	Password         string
	Username         string
}

type TblMemberGroup struct {
	Id          int `gorm:"primaryKey;auto_increment;"`
	Name        string
	Slug        string
	Description string
	IsActive    int
	IsDeleted   int
	CreatedOn   time.Time `gorm:"DEFAULT:NULL"`
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	DateString  string    `gorm:"-"`
}

type MemberLogin struct {
	Username string
	Password string
}

type MemberCreation struct {
	FirstName        string
	LastName         string
	Email            string
	MobileNo         string
	IsActive         int
	ProfileImage     string
	ProfileImagePath string
	Username         string
	Password         string
	GroupId          int
}

type MemberGroupCreation struct {
	Name        string
	Description string
}

type Filter struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}


// Member Group List

func (as Authstruct) MemberGroupList(membergroup []TblMemberGroup, limit int, offset int, filter Filter, DB *gorm.DB) (membergroupl []TblMemberGroup, TotalMemberGroup int64, err error) {

	query := DB.Table("tbl_member_group").Where("is_deleted = 0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")

	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&membergroup)

		return membergroup, 0, err

	} else {

		query.Find(&membergroup).Count(&TotalMemberGroup)

		return membergroup, TotalMemberGroup, err
	}

	return []TblMemberGroup{}, 0, nil
}

// Member Group Insert
func (as Authstruct) MemberGroupCreate(membergroup *TblMemberGroup, DB *gorm.DB) error {

	if err := DB.Table("tbl_members").Create(&membergroup).Error; err != nil {

		return err
	}

	return nil
}

// Member Group Update
func (as Authstruct) MemberGroupUpdate(membergroup *TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_group").Where("id=?", id).Updates(TblMemberGroup{Name: membergroup.Name, Slug: membergroup.Slug, Description: membergroup.Description, Id: membergroup.Id, ModifiedOn: membergroup.ModifiedOn, ModifiedBy: membergroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

// Delete the member group data
func (as Authstruct) MemberGroupDelete(membergroup *TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_group").Where("id=?", id).Updates(TblMemberGroup{IsDeleted: membergroup.IsDeleted}).Error; err != nil {

		return err

	}

	return nil
}
