package teams

import (
	"time"

	"gorm.io/gorm"
)

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
	LastLogin            time.Time `gorm:"DEFAULT:NULL"`
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
	RoleName             string    `gorm:"-:migration;<-:false"`
	DefaultLanguageId    int
	NameString           string `gorm:"-"`
}

type Filters struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

type TeamCreate struct {
	FirstName        string
	LastName         string
	RoleId           int
	Email            string
	Username         string
	Password         string
	IsActive         int
	DataAccess       int
	MobileNo         string
	ProfileImage     string
	ProfileImagePath string
}

// This func will help to create a user in your database
func (t Team) CreateUser(user *TblUser, DB *gorm.DB) error {

	if err := DB.Create(&user).Error; err != nil {

		return err

	}

	return nil
}

func (t Team) GetUsersList(users *[]TblUser, offset, limit int, filter Filters, flag bool, DB *gorm.DB) ([]TblUser, int64) {

	var Total_users int64

	query := DB.Table("tbl_users").Select("tbl_users.id,tbl_users.uuid,tbl_users.role_id,tbl_users.first_name,tbl_users.last_name,tbl_users.email,tbl_users.password,tbl_users.username,tbl_users.mobile_no,tbl_users.profile_image,tbl_users.profile_image_path,tbl_users.created_on,tbl_users.created_by,tbl_users.modified_on,tbl_users.modified_by,tbl_users.is_active,tbl_users.is_deleted,tbl_users.deleted_on,tbl_users.deleted_by,tbl_users.data_access,tbl_roles.name as role_name").
		Joins("inner join tbl_roles on tbl_users.role_id = tbl_roles.id").Where("tbl_users.is_deleted=?", 0)

	if filter.Keyword != "" {

		query = query.Where("(LOWER(TRIM(tbl_users.first_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
			Or("LOWER(TRIM(tbl_users.last_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
			Or("LOWER(TRIM(tbl_roles.name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%").
			Or("LOWER(TRIM(tbl_users.username)) ILIKE LOWER(TRIM(?)))", "%"+filter.Keyword+"%")

	}

	if flag {

		query.Order("id desc").Find(&users)

		return *users, 0

	}

	if limit != 0 && !flag {

		query.Offset(offset).Limit(limit).Order("id desc").Find(&users)

		return *users, 0

	} else {

		query.Find(&users).Count(&Total_users)

		return []TblUser{}, Total_users
	}

}

func (t Team) GetUserDetailsTeam(user *TblUser, id int, DB *gorm.DB) error {

	if err := DB.Where("id=?", id).First(&user).Error; err != nil {

		return err
	}
	return nil
}

func (t Team) UpdateUser(user *TblUser, imgdata string, DB *gorm.DB) error {

	query := DB.Table("tbl_users").Where("id=?", user.Id)

	if imgdata == "" || user.Password == "" {

		if user.Password == "" && imgdata == "" {

			query = query.Omit("password", "profile_image", "profile_image_path").UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "data_access": user.DataAccess})

		} else if imgdata == "" {

			query = query.Omit("profile_image", "profile_image_path").UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "data_access": user.DataAccess, "password": user.Password})

		} else if user.Password == "" {

			query = query.Omit("password").UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "profile_image": user.ProfileImage, "profile_image_path": user.ProfileImagePath, "data_access": user.DataAccess})
		}

		if err := query.Error; err != nil {

			return err
		}

	} else {

		if err := query.UpdateColumns(map[string]interface{}{"first_name": user.FirstName, "last_name": user.LastName, "role_id": user.RoleId, "email": user.Email, "username": user.Username, "mobile_no": user.MobileNo, "is_active": user.IsActive, "modified_on": user.ModifiedOn, "modified_by": user.ModifiedBy, "profile_image": user.ProfileImage, "profile_image_path": user.ProfileImagePath, "data_access": user.DataAccess, "password": user.Password}).Error; err != nil {

			return err
		}

	}

	return nil
}

func (t Team) DeleteUser(user *TblUser, DB *gorm.DB) error {

	if err := DB.Model(&user).Where("id=?", user.Id).Updates(TblUser{IsDeleted: user.IsDeleted, DeletedOn: user.DeletedOn, DeletedBy: user.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

func (t Team) CheckEmail(user *TblUser, email string, userid int, DB *gorm.DB) error {

	if userid == 0 {
		if err := DB.Table("tbl_users").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and is_deleted = 0 ", email).First(&user).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_users").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and id not in(?) and is_deleted= 0 ", email, userid).First(&user).Error; err != nil {

			return err
		}
	}
	return nil
}

func (t Team) CheckNumber(user *TblUser, mobile string, userid int, DB *gorm.DB) error {
	if userid == 0 {
		if err := DB.Table("tbl_users").Where("mobile_no = ? and is_deleted=0", mobile).First(&user).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_users").Where("mobile_no = ? and id not in (?) and is_deleted=0", mobile, userid).First(&user).Error; err != nil {

			return err
		}

	}

	return nil
}

func (t Team) CheckValidation(user *TblUser, email, username, mobile string, userid int, DB *gorm.DB) error {
	if userid == 0 {
		if err := DB.Table("tbl_users").Where("mobile_no = ? or LOWER(TRIM(email))=LOWER(TRIM(?)) or username = ?   and is_deleted=0", mobile, email, username).First(&user).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_users").Where("mobile_no = ? or LOWER(TRIM(email))=LOWER(TRIM(?)) or username = ? and id not in (?) and is_deleted=0", mobile, email, username, userid).First(&user).Error; err != nil {

			return err
		}

	}

	return nil
}

func (t Team) CheckUsername(user *TblUser, username string, userid int, DB *gorm.DB) error {
	if userid == 0 {
		if err := DB.Table("tbl_users").Where("username = ? and is_deleted=0", username).First(&user).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_users").Where("username = ? and id not in (?) and is_deleted=0", username, userid).First(&user).Error; err != nil {

			return err
		}

	}

	return nil
}

func (t Team) ChangePasswordById(user *TblUser, DB *gorm.DB) error {

	if err := DB.Model(&user).Where("id=?", user.Id).Updates(TblUser{Password: user.Password, ModifiedOn: user.ModifiedOn, ModifiedBy: user.ModifiedBy}).Error; err != nil {

		return err
	}
	return nil
}

// Delete the role data
func (t Team) UserUsedRoleCheck(user *TblUser, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_users").Where("role_id=? and is_deleted =0", id).Find(user).Error; err != nil {

		return err
	}

	return nil
}

// Rolechekc
func (t Team) CheckRoleUsed(user *TblUser, roleid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_users").Where("role_id=? and is_deleted =0", roleid).Find(user).Error; err != nil {
		return err
	}
	return nil

}

func (t Team) UserCount(DB *gorm.DB) (count int64, err error) {

	if err := DB.Table("tbl_users").Where("is_deleted = 0 ").Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}

func (t Team) NewuserCount(DB *gorm.DB) (count int64, err error) {

	if err := DB.Table("tbl_users").Where("is_deleted = 0 AND created_on >=?", time.Now().AddDate(0, 0, -10)).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}

func (t Team) Lastlogin(id int, log_time time.Time, DB *gorm.DB) error {

	if err := DB.Table("tbl_users").Where("id=? ", id).Update("last_login", log_time).Error; err != nil {

		return err
	}
	return nil

}
