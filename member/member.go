package member

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/spurtcms/spurtcms-core/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type Memberauth struct {
	Authority *auth.Authority
}

type Image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
}

func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(&TblMemberGroup{}, &TblMember{})

	db.Exec(`CREATE INDEX IF NOT EXISTS email_unique
    ON public.tbl_members USING btree
    (email COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE is_deleted = 0;`)

	db.Exec(`CREATE INDEX IF NOT EXISTS mobile_no_unique
    ON public.tbl_members USING btree
    (mobile_no COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE is_deleted = 0;`)

}

type Filters struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

/*List member group*/
func (a Memberauth) ListMemberGroup(offset, limit int, filter Filters) (membergroup []TblMemberGroup, MemberGroupCount int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMemberGroup{}, 0, checkerr
	}

	check, _ := a.Authority.IsGranted("Member Group", auth.Read)

	if check {

		var TotalMemberGroup int64

		var membergroup []TblMemberGroup

		query := a.Authority.DB.Table("tbl_member_group").Where("is_deleted = 0").Order("id desc")

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

	}
	return []TblMemberGroup{}, 0, errors.New("not authorized")

}

func (a Memberauth) GetGroupData() (membergrouplists []TblMemberGroup, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMemberGroup{}, checkerr
	}

	var membergrouplist []TblMemberGroup

	if err := a.Authority.DB.Table("tbl_member_group").Where("is_deleted = 0 and is_active = 1").Find(&membergrouplist).Error; err != nil {

		return []TblMemberGroup{}, err

	}

	return membergrouplist, nil

	// check, _ := a.Authority.IsGranted("Member Group", auth.Read)

	// if check {

	// var membergrouplist []TblMemberGroup

	// if err := a.Authority.DB.Table("tbl_member_group").Where("is_deleted = 0").Find(&membergrouplist).Error; err != nil {

	// 	return membergrouplist, err

	// }
	// }

}

/*Create Member Group*/
func (a Memberauth) CreateMemberGroup(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Create)

	if err != nil {

		return err
	}

	if check {

		if c.PostFormValue("membergroup_name") == "" || c.PostFormValue("membergroup_desc") == "" {

			return errors.New("given value is empty")
		}

		var membergroup TblMemberGroup

		membergroup.Name = c.PostFormValue("membergroup_name")

		membergroup.Slug = strings.ToLower(c.PostFormValue("membergroup_name"))

		membergroup.Description = c.PostFormValue("membergroup_desc")

		membergroup.CreatedBy = userid

		membergroup.IsActive = 1

		membergroup.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.Authority.DB.Table("tbl_member_group").Create(&membergroup).Error; err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

/*Update Member Group*/
func (a Memberauth) UpdateMemberGroup(c *http.Request, id int) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Update)

	if err != nil {

		return err
	}

	if check {

		if c.PostFormValue("membergroup_name") == "" || c.PostFormValue("membergroup_desc") == "" {

			return errors.New("given value is empty")
		}

		var membergroup TblMemberGroup

		membergroup.Id = id

		membergroup.Name = c.PostFormValue("membergroup_name")

		membergroup.Slug = strings.ToLower(c.PostFormValue("membergroup_name"))

		membergroup.Description = c.PostFormValue("membergroup_desc")

		membergroup.ModifiedBy = userid

		membergroup.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.Authority.DB.Table("tbl_member_group").Where("id=?", membergroup.Id).Updates(TblMemberGroup{Name: membergroup.Name, Slug: membergroup.Slug, Description: membergroup.Description, Id: membergroup.Id, ModifiedOn: membergroup.ModifiedOn, ModifiedBy: membergroup.ModifiedBy}).Error; err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

/*Delete Member Group*/
func (a Memberauth) DeleteMemberGroup(id int) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Update)

	if err != nil {

		return err
	}

	if check {

		var membergroup TblMemberGroup

		membergroup.IsDeleted = 1

		if id <= 0 {

			return errors.New("invalid id cannot delete")

		}

		if err := a.Authority.DB.Table("tbl_member_group").Where("id=?", id).Updates(TblMemberGroup{IsDeleted: membergroup.IsDeleted}).Error; err != nil {

			return err

		}

	} else {
		return errors.New("not authorized")
	}

	return nil
}

// list member
func (a Memberauth) ListMembers(offset, limit int, filter Filters, flag bool) (member []TblMember, totoalmember int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMember{}, 0, checkerr
	}
	check, _ := a.Authority.IsGranted("Member", auth.Read)

	if check {
		var member []TblMember

		var Total_Member int64

		query := a.Authority.DB.Table("tbl_members").Select("tbl_members.id,tbl_members.uuid,tbl_members.member_group_id,tbl_members.first_name,tbl_members.last_name,tbl_members.email,tbl_members.mobile_no,tbl_members.profile_image,tbl_members.profile_image_path,tbl_members.created_on,tbl_members.created_by,tbl_members.modified_on,tbl_members.modified_by,tbl_members.is_active,tbl_members.is_deleted,tbl_members.deleted_on,tbl_members.deleted_by").
			Joins("left join tbl_member_group on tbl_members.member_group_id = tbl_member_group.id").Where("tbl_members.is_deleted=?", 0)

		if filter.Keyword != "" {

			query = query.Where("(LOWER(TRIM(tbl_members.first_name)) ILIKE LOWER(TRIM(?))"+" OR LOWER(TRIM(tbl_members.last_name)) ILIKE LOWER(TRIM(?))"+" OR LOWER(TRIM(tbl_member_group.name)) ILIKE LOWER(TRIM(?)))"+" AND tbl_members.is_deleted=0"+" AND tbl_member_group.is_deleted=0", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")

		}
		if flag {

			query.Order("id desc").Find(&member)

			return member, 0, err

		}

		if limit != 0 && !flag{

			query.Offset(offset).Limit(limit).Order("id desc").Find(&member)

			return member, 0, err

		} else {
			query.Find(&member).Count(&Total_Member)

			return member, Total_Member, nil
		}

	}

	return []TblMember{}, 0, errors.New("not authorized")

}

// Create Member
func (a Memberauth) CreateMember(c *http.Request, imagename string, imagepath string) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Create)

	if err != nil {

		return err
	}

	if check {

		uvuid := (uuid.New()).String()

		var member TblMember

		member.Uuid = uvuid

		member.ProfileImage = imagename

		member.ProfileImagePath = imagepath

		member.MemberGroupId, _ = strconv.Atoi(c.PostFormValue("mem_group"))

		member.FirstName = c.PostFormValue("mem_name")

		member.LastName = c.PostFormValue("mem_lname")

		member.Email = c.PostFormValue("mem_email")

		member.MobileNo = c.PostFormValue("mem_mobile")

		member.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))

		member.Username = c.PostFormValue("mem_usrname")

		password := c.PostFormValue("mem_pass")

		if password != "" {

			hash_pass := hashingPassword(password)
			member.Password = hash_pass

		}

		member.CreatedBy = userid

		member.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.Authority.DB.Table("tbl_members").Create(&member).Error; err != nil {

			return err

		}

	} else {

		return errors.New("not authorized")
	}

	return nil

}

// Update Member
func (a Memberauth) UpdateMember(c *http.Request, id int, imageName string, imagePath string) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Update)

	if err != nil {

		return err
	}

	if check {

		uvuid := (uuid.New()).String()

		var member TblMember

		member.Uuid = uvuid

		member.Id = id

		member.MemberGroupId, _ = strconv.Atoi(c.PostFormValue("mem_group"))

		member.FirstName = c.PostFormValue("mem_name")

		member.LastName = c.PostFormValue("mem_lname")

		member.Email = c.PostFormValue("mem_email")

		member.MobileNo = c.PostFormValue("mem_mobile")

		member.ProfileImage = imageName

		member.ProfileImagePath = imagePath

		member.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))

		member.ModifiedBy = userid

		member.Username = c.PostFormValue("mem_usrname")

		password := c.PostFormValue("mem_pass")

		if password != "" {

			hash_pass := hashingPassword(password)
			member.Password = hash_pass

		}

		member.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		query := a.Authority.DB.Table("tbl_members").Where("id=?", member.Id)

		if member.Password == "" && member.ProfileImage == "" && member.ProfileImagePath == "" {

			query.Omit("password , profile_image , profile_image_path").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy})

			return err

		} else {
			query.UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath, "password": member.Password})
			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

// delete member
func (a Memberauth) DeleteMember(id int) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Delete)

	if err != nil {

		return err
	}

	if check {

		var member TblMember

		member.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		member.DeletedBy = userid

		if err := a.Authority.DB.Table("tbl_members").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

			return err

		}

	} else {

		return errors.New("not authorized")

	}

	return nil
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg"
}

// Process uploaded file into an image.
func Process(r http.Request, field string) (Image, error) {
	file, info, err := r.FormFile(field)

	if err != nil {
		return Image{}, err
	}

	contentType := info.Header.Get("Content-Type")

	if !okContentType(contentType) {
		return Image{}, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := io.ReadAll(file)

	if err != nil {
		return Image{}, err
	}

	_, _, err = image.Decode(bytes.NewReader(bs))

	if err != nil {
		return Image{}, err
	}

	i := &Image{
		Filename:    info.Filename,
		ContentType: contentType,
		Data:        bs,
		Size:        len(bs),
	}

	return *i, nil
}

// Check Email is already exits or not

func (a Memberauth) CheckEmailInMember(id int, email string) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Create)

	if err != nil {

		return err
	}

	if check {
		var member TblMember

		if id == 0 {

			if err := a.Authority.DB.Table("tbl_members").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and is_deleted=0", email).First(&member).Error; err != nil {

				return err
			}
		} else {

			if err := a.Authority.DB.Debug().Table("tbl_members").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and id not in (?) and is_deleted = 0 ", email, id).First(&member).Error; err != nil {

				return err
			}
		}

	}
	return nil
}

// Check Number is already exits or not

func (a Memberauth) CheckNumberInMember(id int, number string) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Create)

	if err != nil {

		return err
	}

	if check {

		var member TblMember

		if id == 0 {

			if err := a.Authority.DB.Table("tbl_members").Where("mobile_no = ? and is_deleted = 0", number).First(&member).Error; err != nil {

				return err
			}
		} else {

			if err := a.Authority.DB.Debug().Table("tbl_members").Where("mobile_no = ? and id not in (?) and is_deleted = 0", number, id).First(&member).Error; err != nil {

				return err
			}
		}

	}
	return nil
}

// member group delete popup
func (a Memberauth) MemberDeletePopup(id int) (member TblMember, err1 error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return TblMember{}, checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Read)

	if err != nil {

		return TblMember{}, err
	}

	if check {

		var member TblMember

		if err := a.Authority.DB.Table("tbl_members").Where("member_group_id=? and is_deleted = 0", id).Find(&member).Error; err != nil {

			return TblMember{}, err
		}
	}
	return member, nil

}

// member is_active

func (a Memberauth) MemberIsActive(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Read)

	if err != nil {

		return err
	}

	if check {

		id, _ := strconv.Atoi(c.PostFormValue("id"))

		val, _ := strconv.Atoi(c.PostFormValue("isactive"))

		var memberstatus TblMemberGroup

		memberstatus.ModifiedBy = userid

		memberstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.Authority.DB.Table("tbl_member_group").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_active": val, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {

		return errors.New("not authorized")
	}

	return nil

}

//

func (a Memberauth) GetMemberDetails(id int) (members TblMember, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return TblMember{}, checkerr
	}

	// check, err := a.Authority.IsGranted("Member", auth.Create)

	// if err != nil {

	// 	return TblMember{}, err
	// }

	// if check {

	var member TblMember

	fmt.Println("sfd", id)

	if err := a.Authority.DB.Table("tbl_members").Where("id=?", id).First(&member).Error; err != nil {

		return TblMember{}, err
	}

	// }
	return member, nil

}

func (a Memberauth) GetMemberById(id int) (membergroup TblMemberGroup, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return TblMemberGroup{}, checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Read)

	if err != nil {

		return TblMemberGroup{}, err
	}

	if check {

		var membergroup TblMemberGroup

		if err := a.Authority.DB.Table("tbl_member_group").Where("id=?", id).First(&membergroup).Error; err != nil {

			return TblMemberGroup{}, err
		}

	}
	return membergroup, nil

}

/*Create meber token*/
func CreateMemberToken(userid, roleid int, secretkey string) (string, error) {

	atClaims := jwt.MapClaims{}

	atClaims["member_id"] = userid

	atClaims["group_id"] = roleid

	atClaims["expiry_time"] = time.Now().Add(2 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	return token.SignedString([]byte(secretkey))
}

/*Member login*/
func CheckMemberLogin(memlogin MemberLogin, db *gorm.DB, secretkey string) (string, error) {

	username := memlogin.Username

	password := memlogin.Password

	var member TblMember

	if err := db.Table("tbl_members").Where("username = ?", username).First(&member).Error; err != nil {

		return "", err

	}

	passerr := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))

	if passerr != nil || passerr == bcrypt.ErrMismatchedHashAndPassword {

		return "", errors.New("invalid password")

	}

	token, err := CreateMemberToken(member.Id, member.MemberGroupId, secretkey)

	if err != nil {

		return "", err
	}

	return token, nil

}

// verify token
func VerifyToken(token string, secret string) (memberid, groupid int, err error) {
	Claims := jwt.MapClaims{}

	tkn, err := jwt.ParseWithClaims(token, Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(err)
			return 0, 0, errors.New("invalid token")
		}

		return 0, 0, errors.New(err.Error())
	}

	if !tkn.Valid {
		fmt.Println(tkn)
		return 0, 0, errors.New("invalid token")
	}

	usrid := Claims["member_id"]

	rolid := Claims["group_id"]

	return int(usrid.(float64)), int(rolid.(float64)), nil
}

// Check Name is already exits or not

func (a Memberauth) CheckNameInMember(id int, name string) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Create)

	if err != nil {

		return err
	}

	if check {

		var member TblMember

		if id == 0 {

			if err := a.Authority.DB.Table("tbl_members").Where("username = ? and is_deleted = 0", name).First(&member).Error; err != nil {

				return err
			}
		} else {

			if err := a.Authority.DB.Debug().Table("tbl_members").Where("username = ? and id not in (?) and is_deleted = 0", name, id).First(&member).Error; err != nil {

				return err
			}
		}

	}
	return nil
}

func hashingPassword(pass string) string {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		panic(err)

	}

	return string(passbyte)
}