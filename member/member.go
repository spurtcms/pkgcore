package member

import (
	"errors"
	"fmt"
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
	Authority *auth.Authorization
}

type MemberAuth struct {
	Auth *auth.Authorization
}

type Image struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
}

type Authstruct struct{}

var AS Authstruct

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

/*List member group*/
func (a Memberauth) ListMemberGroup(offset, limit int, filter Filter) (membergroup []TblMemberGroup, MemberGroupCount int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMemberGroup{}, 0, checkerr
	}

	check, _ := a.Authority.IsGranted("Member Group", auth.Read)

	if check {

		if err != nil {

			return []TblMemberGroup{}, 0, err
		}

		var membergrplist []TblMemberGroup

		AS.MemberGroupList(membergrplist, limit, offset, filter, a.Authority.DB)

		_, membercounts, _ := AS.MemberGroupList(membergrplist, limit, offset, filter, a.Authority.DB)

		membergrouplist, _, _ := AS.MemberGroupList(membergrplist, limit, offset, filter, a.Authority.DB)

		var membergrouplists []TblMemberGroup

		for _, val := range membergrouplist {

			if !val.ModifiedOn.IsZero() {

				val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {
				val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}

			membergrouplists = append(membergrouplists, val)

		}

		return membergrouplists, membercounts, nil
	}
	return []TblMemberGroup{}, 0, errors.New("not authorized")

}

/*Create Member Group*/
func (a Memberauth) CreateMemberGroup(membergrpc MemberGroupCreation) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Create)

	if err != nil {

		return err
	}

	if check {

		if membergrpc.Name == "" || membergrpc.Description == "" {

			return errors.New("given value is empty")
		}

		var membergroup TblMemberGroup

		membergroup.Name = membergrpc.Name

		membergroup.Slug = strings.ToLower(membergrpc.Name)

		membergroup.Description = membergrpc.Description

		membergroup.CreatedBy = userid

		membergroup.IsActive = 1

		membergroup.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := AS.MemberGroupCreate(&membergroup, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil
}

/*Update Member Group*/
func (a Memberauth) UpdateMemberGroup(membergrpc MemberGroupCreation, id int) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Update)

	if err != nil {

		return err
	}
	if check {

		if membergrpc.Name == "" || membergrpc.Description == "" {

			return errors.New("given value is empty")
		}

		var membergroup TblMemberGroup

		membergroup.Name = membergrpc.Name

		membergroup.Slug = strings.ToLower(membergrpc.Name)

		membergroup.Description = membergrpc.Description

		membergroup.ModifiedBy = userid

		membergroup.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := AS.MemberGroupUpdate(&membergroup, id, a.Authority.DB)

		if err != nil {

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

		if id <= 0 {

			return errors.New("invalid id cannot delete")

		}
		membergroup.IsDeleted = 1

		err := AS.MemberGroupDelete(&membergroup, id, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {
		return errors.New("not authorized")
	}

	return nil
}

// list member
func (a Memberauth) ListMembers(offset int, limit int, filter Filter, flag bool) (member []TblMember, totoalmember int64, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMember{}, 0, checkerr
	}
	check, _ := a.Authority.IsGranted("Member", auth.Read)

	if check {

		var member []TblMember

		// var membergroup []TblMemberGroup

		memberlist, _, _ := AS.MembersList(member, offset, limit, filter, flag, a.Authority.DB)

		_, Total_users, _ := AS.MembersList(member, 0, 0, filter, flag, a.Authority.DB)

		// var members []TblMember

		// membergrouplist, _ := AS.GetGroupData(membergroup, a.Authority.DB)

		// for _, member_object := range memberlist {

		// 	member_object.CreatedDate = member_object.CreatedOn.Format("02 Jan 2006 03:04 PM")

		// 	for _, membergrp_object := range membergrouplist {

		// 		if member_object.MemberGroupId == membergrp_object.Id {

		// 			member_object.Group = append(member_object.Group, membergrp_object)

		// 			members = append(members, member_object)

		// 		}

		// 	}
		// }

		return memberlist, Total_users, nil

	}

	return []TblMember{}, 0, errors.New("not authorized")

}

func (a Memberauth) GetGroupData() (membergroup []TblMemberGroup, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMemberGroup{}, checkerr
	}

	check, err := a.Authority.IsGranted("Member Group", auth.Create)

	if err != nil {

		return []TblMemberGroup{}, err
	}

	if check {

		var membergroup []TblMemberGroup

		membergrouplist, _ := AS.GetGroupData(membergroup, a.Authority.DB)

		return membergrouplist, nil

	} else {

		return []TblMemberGroup{}, errors.New("not authorized")
	}

}

// Create Member
func (a Memberauth) CreateMember(Mc MemberCreation) error {

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

		member.ProfileImage = Mc.ProfileImage

		member.ProfileImagePath = Mc.ProfileImagePath

		member.MemberGroupId = Mc.GroupId

		member.FirstName = Mc.FirstName

		member.LastName = Mc.LastName

		member.Email = Mc.Email

		member.MobileNo = Mc.MobileNo

		member.IsActive = Mc.IsActive

		member.Username = Mc.Username

		if Mc.Password != "" {

			hash_pass := hashingPassword(Mc.Password)

			member.Password = hash_pass

		}

		member.CreatedBy = userid

		member.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := AS.MemberCreate(&member, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")
	}

	return nil

}

// Update Member
func (a Memberauth) UpdateMember(Mc MemberCreation, id int) error {

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

		member.MemberGroupId = Mc.GroupId

		member.FirstName = Mc.FirstName

		member.LastName = Mc.LastName

		member.Email = Mc.Email

		member.MobileNo = Mc.MobileNo

		member.ProfileImage = Mc.ProfileImage

		member.ProfileImagePath = Mc.ProfileImagePath

		member.IsActive = Mc.IsActive

		member.ModifiedBy = userid

		member.Username = Mc.Username

		password := Mc.Password

		if password != "" {

			hash_pass := hashingPassword(password)

			member.Password = hash_pass

		}

		member.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := AS.UpdateMember(&member, a.Authority.DB)

		if err != nil {

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

		err := AS.DeleteMember(&member, id, a.Authority.DB)

		if err != nil {

			return err
		}

	} else {

		return errors.New("not authorized")

	}

	return nil
}

// Check Email is already exits or not
func (a Memberauth) CheckEmailInMember(id int, email string) (bool, error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Create)

	if err != nil {

		return false, err
	}

	if check {

		var member TblMember

		err := AS.CheckEmailInMember(&member, email, id, a.Authority.DB)

		if err != nil {

			return false, err
		}

		return true, nil
	}
	return false, errors.New("not authorized")
}

// Check Email is already exits or not
func (a MemberAuth) CheckEmailInMember(id int, email string) (TblMember, bool, error) {

	var member TblMember

	err := AS.CheckEmailInMember(&member, email, id, a.Auth.DB)

	if err != nil {

		return TblMember{}, false, err
	}

	return member, true, nil
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

		err := AS.CheckNumberInMember(&member, number, id, a.Authority.DB)

		if err != nil {
			return err
		}

	}
	return nil
}

// Check Number is already exits or not
func (a MemberAuth) CheckNumberInMember(id int, number string) (bool, error) {

	// _, _, checkerr := auth.VerifyToken(a.Auth.Token, a.Auth.Secret)

	// if checkerr != nil {

	// 	return false, checkerr
	// }

	var member TblMember

	err := AS.CheckNumberInMember(&member, number, id, a.Auth.DB)

	if err != nil {

		return false, err
	}

	return true, nil

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

		fmt.Println("pkg namr", name)

		err := AS.CheckNameInMember(&member, id, name, a.Authority.DB)

		if err != nil {
			return err
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

		member, _ = AS.MemberDeletePopup(id, a.Authority.DB)

		if err != nil {
			return TblMember{}, err
		}
	}
	return member, nil

}

// member is_active
func (a Memberauth) MemberIsActive(memberid int, status int) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := a.Authority.IsGranted("Member", auth.Read)

	if err != nil {

		return err
	}

	if check {

		var memberstatus TblMemberGroup

		memberstatus.ModifiedBy = userid

		memberstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := AS.MemberIsActive(memberstatus, memberid, status, a.Authority.DB)

		if err != nil {
			return err
		}
	} else {

		return errors.New("not authorized")
	}

	return nil

}

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

	err = AS.MemberDetails(&member, id, a.Authority.DB)

	if err != nil {

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

		err1 := AS.GetMemberById(membergroup, id, a.Authority.DB)

		if err1 != nil {

			return TblMemberGroup{}, err1
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
func (M MemberAuth) CheckMemberLogin(memlogin MemberLogin, db *gorm.DB, secretkey string) (string, error) {

	username := memlogin.Emailid

	password := memlogin.Password

	var member TblMember

	if err := db.Table("tbl_members").Where("email = ? and is_deleted=0", username).First(&member).Error; err != nil {

		return "", errors.New("your email not registered")

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

func hashingPassword(pass string) string {

	passbyte, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	if err != nil {

		panic(err)

	}

	return string(passbyte)
}

// This OTP valid only 5 minutes
// updateOTP
func (M MemberAuth) UpdateOtp(otp int, memberid int) (bool, error) {

	// memberid, _, checkerr := VerifyToken(M.Auth.Token, M.Auth.Secret)

	// if checkerr != nil {

	// 	return false, checkerr
	// }

	var tblmember TblMember

	tblmember.Otp = otp

	tblmember.OtpExpiry, _ = time.Parse("2006-01-02 15:04:05", time.Now().Add(time.Duration(5)*time.Minute).In(IST).Format("2006-01-02 15:04:05"))

	err := AS.UpdateOTP(&tblmember, otp, memberid, M.Auth.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

// ChangeEmailid
func (M MemberAuth) ChangeEmailId(otp int, emailid string) (bool, error) {

	memberid, _, checkerr := VerifyToken(M.Auth.Token, M.Auth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var tblmember TblMember

	AS.MemberDetails(&tblmember, memberid, M.Auth.DB)

	if tblmember.Otp != otp {

		return false, errors.New("invalid otp")
	}

	if tblmember.OtpExpiry.Unix() < time.Now().Unix() {

		return false, errors.New("otp exipred")

	}

	err := AS.UpdateEmail(emailid, memberid, M.Auth.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

// ChangePassword
func (M MemberAuth) ChangePassword(otp int, memberid int, password string) (bool, error) {

	// memberid, _, checkerr := VerifyToken(M.Auth.Token, M.Auth.Secret)

	// if checkerr != nil {

	// 	return false, checkerr
	// }

	var tblmember TblMember

	AS.MemberDetails(&tblmember, memberid, M.Auth.DB)

	if tblmember.Otp != otp {

		return false, errors.New("invalid otp")
	}

	if tblmember.OtpExpiry.Unix() < time.Now().Unix() {

		return false, errors.New("otp exipred")

	}

	hashpass := hashingPassword(password)

	err := AS.UpdatePassword(hashpass, memberid, M.Auth.DB)

	if err != nil {

		return false, err
	}

	return true, nil
}

// get member details
func (M MemberAuth) GetMemberDetails() (members TblMember, err error) {

	memberid, _, checkerr := VerifyToken(M.Auth.Token, M.Auth.Secret)

	if checkerr != nil {

		return TblMember{}, checkerr
	}

	var member TblMember

	err1 := AS.MemberDetails(&member, memberid, M.Auth.DB)

	if err1 != nil {

		return TblMember{}, err
	}

	return member, nil

}

// register member
func (M MemberAuth) MemberRegister(MemC MemberCreation) (check bool, err error) {

	if MemC.FirstName != "" {

		return false, errors.New("firstname is empty can't register")

	} else if MemC.Email != "" {

		return false, errors.New("email is empty can't register")

	} else if MemC.MobileNo != "" {

		return false, errors.New("mobile number is empty can't register")

	} else if MemC.Password != "" {

		return false, errors.New("password is empty can't register")
	}

	Pass := hashingPassword(MemC.Password)

	var member TblMember

	member.FirstName = MemC.FirstName

	member.LastName = MemC.LastName

	member.Email = MemC.Email

	member.MobileNo = MemC.MobileNo

	member.IsActive = 1

	member.MemberGroupId = 1

	member.Username = MemC.Username

	member.Password = Pass

	member.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	member.CreatedBy = 1

	err1 := AS.MemberCreate(&member, M.Auth.DB)

	if err1 != nil {

		return false, err
	}

	return true, nil

}

// Update member
func (M MemberAuth) MemberUpdate(MemC MemberCreation) (check bool, err error) {

	memberid, _, checkerr := VerifyToken(M.Auth.Token, M.Auth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var member TblMember

	member.FirstName = MemC.FirstName

	member.LastName = MemC.LastName

	member.MobileNo = MemC.MobileNo

	member.ProfileImage = MemC.ProfileImage

	member.ProfileImagePath = MemC.ProfileImagePath

	member.Password = hashingPassword(MemC.Password)

	member.Email = MemC.Email

	// member.IsActive = MemC.IsActive

	// member.Username = MemC.Username

	// member.MemberGroupId = MemC.GroupId

	member.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	member.ModifiedBy = 1

	err1 := AS.MemberUpdate(&member, memberid, M.Auth.DB)

	if err1 != nil {

		return false, err
	}

	return true, nil

}
