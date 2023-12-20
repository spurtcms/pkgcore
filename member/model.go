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
	GroupName        string `gorm:"-:migration;<-:false"`
	Password         string
	DateString       string    `gorm:"-"`
	Username         string    `gorm:"DEFAULT:NULL"`
	Otp              int       `gorm:"DEFAULT:NULL"`
	OtpExpiry        time.Time `gorm:"DEFAULT:NULL"`
	ModifiedDate     string    `gorm:"-"`
	NameString       string    `gorm:"-"`
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
	Emailid  string
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

func (as Authstruct) MemberGroupList(membergroup []TblMemberGroup, limit int, offset int, filter Filter, getactive bool, DB *gorm.DB) (membergroupl []TblMemberGroup, TotalMemberGroup int64, err error) {

	query := DB.Model(TblMemberGroup{}).Where("is_deleted = 0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")

	}

	if getactive {

		query = query.Where("is_active=1")

	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&membergroup)

		return membergroup, 0, err

	}

	query.Find(&membergroup).Count(&TotalMemberGroup)

	return membergroup, TotalMemberGroup, err

}

// Member Group Insert
func (as Authstruct) MemberGroupCreate(membergroup *TblMemberGroup, DB *gorm.DB) error {

	if err := DB.Model(TblMemberGroup{}).Create(&membergroup).Error; err != nil {

		return err
	}

	return nil
}

// Member Group Update
func (as Authstruct) MemberGroupUpdate(membergroup *TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Model(TblMemberGroup{}).Where("id=?", id).Updates(TblMemberGroup{Name: membergroup.Name, Slug: membergroup.Slug, Description: membergroup.Description, Id: membergroup.Id, ModifiedOn: membergroup.ModifiedOn, ModifiedBy: membergroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

// Delete the member group data
func (as Authstruct) MemberGroupDelete(membergroup *TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Model(TblMemberGroup{}).Where("id=?", id).Updates(TblMemberGroup{IsDeleted: membergroup.IsDeleted}).Error; err != nil {

		return err

	}

	return nil
}

// Member list
func (as Authstruct) MembersList(member []TblMember, limit int, offset int, filter Filter, flag bool, DB *gorm.DB) (memberl []TblMember, Total_Member int64, err error) {

	query := DB.Model(TblMember{}).Select("tbl_members.id,tbl_members.uuid,tbl_members.member_group_id,tbl_members.first_name,tbl_members.last_name,tbl_members.email,tbl_members.mobile_no,tbl_members.profile_image,tbl_members.profile_image_path,tbl_members.created_on,tbl_members.created_by,tbl_members.modified_on,tbl_members.modified_by,tbl_members.is_active,tbl_members.is_deleted,tbl_members.deleted_on,tbl_members.deleted_by,tbl_member_groups.name as group_name").
		Joins("inner join tbl_member_groups on tbl_members.member_group_id = tbl_member_groups.id").Where("tbl_members.is_deleted=?", 0).Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("(LOWER(TRIM(tbl_members.first_name)) ILIKE LOWER(TRIM(?))"+" OR LOWER(TRIM(tbl_members.last_name)) ILIKE LOWER(TRIM(?))"+" OR LOWER(TRIM(tbl_member_groups.name)) ILIKE LOWER(TRIM(?)))"+" AND tbl_members.is_deleted=0"+" AND tbl_member_groups.is_deleted=0", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")

	}
	if flag {

		query.Find(&member)

		return member, 0, err

	}

	if limit != 0 && !flag {

		query.Offset(offset).Limit(limit).Order("id desc").Find(&member)

		return member, 0, err

	}
	query.Find(&member).Count(&Total_Member)

	return member, Total_Member, nil

}

func (as Authstruct) GetGroupData(membergroup []TblMemberGroup, DB *gorm.DB) (membergrouplists []TblMemberGroup, err error) {

	var membergrouplist []TblMemberGroup

	if err := DB.Model(TblMemberGroup{}).Where("is_deleted = 0 and is_active = 1").Find(&membergrouplist).Error; err != nil {

		return []TblMemberGroup{}, err

	}

	return membergrouplist, nil

}

// Member Insert
func (as Authstruct) MemberCreate(member *TblMember, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Create(&member).Error; err != nil {

		return err
	}

	return nil
}

// Update Member
func (as Authstruct) UpdateMember(member *TblMember, DB *gorm.DB) error {

	query := DB.Model(TblMember{}).Where("id=?", member.Id)

	if member.Password == "" && member.ProfileImage == "" && member.ProfileImagePath == "" {

		query.Omit("password , profile_image , profile_image_path").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy})

	} else {

		query.UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath, "password": member.Password})
	}
	return nil
}

// Get Member Details

func (as Authstruct) MemberDetails(member *TblMember, memberid int, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Where("id=?", memberid).First(&member).Error; err != nil {
		return err

	}

	return nil
}

// Delete Member
func (as Authstruct) DeleteMember(member *TblMember, id int, DB *gorm.DB) error {

	if err := DB.Model(&member).Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

// Check Email is already exists
func (AS Authstruct) CheckEmailInMember(member *TblMember, email string, userid int, DB *gorm.DB) error {

	if userid == 0 {
		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and is_deleted=0", email).First(&member).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and id not in (?) and is_deleted = 0 ", email, userid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

func (As Authstruct) CheckNumberInMember(member *TblMember, number string, userid int, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Model(TblMember{}).Where("mobile_no = ? and is_deleted = 0", number).First(&member).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblMember{}).Where("mobile_no = ? and id not in (?) and is_deleted=0", number, userid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

// upateotp
func (As Authstruct) UpdateOTP(tblmem *TblMember, otp int, memberid int, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Where("id=?", memberid).UpdateColumns(map[string]interface{}{"otp": tblmem.Otp, "otp_expiry": tblmem.OtpExpiry}).Error; err != nil {

		return err
	}

	return nil
}

// updateemail
func (As Authstruct) UpdateEmail(email string, memberid int, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Where("id=?", memberid).UpdateColumns(map[string]interface{}{"email": email}).Error; err != nil {

		return err
	}

	return nil
}

// updatePassword
func (As Authstruct) UpdatePassword(password string, memberid int, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Where("id=?", memberid).UpdateColumns(map[string]interface{}{"password": password}).Error; err != nil {

		return err
	}

	return nil
}

// Member la IsActive Function
func (As Authstruct) MemberIsActive(memberstatus TblMemberGroup, memberid int, status int, DB *gorm.DB) error {

	if err := DB.Model(TblMemberGroup{}).Where("id=?", memberid).UpdateColumns(map[string]interface{}{"is_active": status, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// Delete Popup
func (As Authstruct) MemberDeletePopup(id int, DB *gorm.DB) (member TblMember, err error) {

	if err := DB.Model(TblMember{}).Where("member_group_id=? and is_deleted = 0", id).Find(&member).Error; err != nil {

		return TblMember{}, err
	}

	return member, nil
}

// Get Member group data
func (As Authstruct) GetMemberById(membergroup TblMemberGroup, id int, DB *gorm.DB) (err error) {

	if err := DB.Model(TblMemberGroup{}).Where("id=?", id).First(&membergroup).Error; err != nil {

		return err
	}

	return nil
}

// Name already exists
func (As Authstruct) CheckNameInMember(member *TblMember, userid int, name string, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and is_deleted=0", name).First(&member).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0", name, userid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

func (AS Authstruct) MemberUpdate(member *TblMember, id int, DB *gorm.DB) error {

	if err := DB.Model(TblMember{}).Where("id=?", id).UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "mobile_no": member.MobileNo, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "email": member.Email, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath,
		"password": member.Password}).Error; err != nil {

		return err
	}

	return nil
}

// Group Name already exists
func (As Authstruct) CheckNameInMemberGroup(member *TblMemberGroup, userid int, name string, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and is_deleted=0", name).First(&member).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblMember{}).Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0", name, userid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

// Update Member Lms
func (as Authstruct) UpdateMemberLms(member *TblMember, DB *gorm.DB) error {

	query := DB.Model(TblMember{}).Where("id=?", member.Id)

	if member.ProfileImage == "" && member.ProfileImagePath == "" {

		query.Omit("profile_image , profile_image_path").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "mobile_no": member.MobileNo, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy})

	} else {

		query.UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "mobile_no": member.MobileNo, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath})
	}
	return nil
}
