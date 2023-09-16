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

	"github.com/google/uuid"
	"github.com/spurtcms/spurtcms-core/auth"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type Memberauth struct {
	Authority auth.Authority
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

type MemberRequired struct {
	Request *http.Request
	Param   string
	Limit   int
	Offset  int
	filter  Filter
}

type Filter struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

func MemberInstance(a *auth.Option) auth.Authority {

	auth := auth.Authority{
		DB:     a.DB,
		Token:  a.Token,
		Secret: a.Secret,
	}

	MigrateTables(auth.DB)

	return auth
}

/*List member group*/
func (a Memberauth) ListMemberGroup(mem MemberRequired) (membergrouplis []TblMemberGroup, err error) {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return []TblMemberGroup{}, checkerr
	}

	var membergrouplist []TblMemberGroup

	check, _ := auth.Authority.IsGranted(a.Authority, "Member Group", auth.Create)

	if check {

		query := a.Authority.DB.Table("tbl_member_group").Where("is_deleted = 0").Order("id desc")

		if mem.filter.Keyword != "" {

			query = query.Where("LOWER(TRIM(name)) ILIKE LOWER(TRIM(?))", "%"+mem.filter.Keyword+"%")

		}

		if mem.Limit != 0 {

			query.Limit(mem.Limit).Offset(mem.Offset).Find(&membergrouplist)

		}

		if err := query.Error; err != nil {

			return []TblMemberGroup{}, err
		}

	} else {

		return []TblMemberGroup{}, errors.New("not authorized")
	}

	return membergrouplist, nil
}

/*Create Member Group*/
func (a Memberauth) CreateMemberGroup(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := auth.Authority.IsGranted(a.Authority, "Member Group", auth.Create)

	if err != nil {

		return err
	}

	if check {

		if c.PostFormValue("membergroup_name") == "" || c.PostFormValue(
			"membergroup_desc") == "" {

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
func (a Memberauth) UpdateMemberGroup(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := auth.Authority.IsGranted(a.Authority, "Member Group", auth.Update)

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

		membergroup.ModifiedBy = userid

		membergroup.Id, _ = strconv.Atoi(c.PostFormValue("membergroup_id"))

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
func (a Memberauth) DeleteMemberGroup(c *http.Request) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := auth.Authority.IsGranted(a.Authority, "Member Group", auth.Update)

	if err != nil {

		return err
	}

	if check {

		var membergroup TblMemberGroup

		membergroup.IsDeleted = 1

		MemberGroupId, _ := strconv.Atoi(c.URL.Query().Get("id"))

		if MemberGroupId == 0 {

			return errors.New("internal error server")

		}

		if err := a.Authority.DB.Table("tbl_member_group").Where("id=?", MemberGroupId).Updates(TblMemberGroup{IsDeleted: membergroup.IsDeleted}).Error; err != nil {

			return err

		}

	} else {
		return errors.New("not authorized")
	}

	return nil
}

// list member
func (a Memberauth) ListMembers(mem MemberRequired) error {

	_, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	return nil
}

// Create Member
func (a Memberauth) CreateMember(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := auth.Authority.IsGranted(a.Authority, "Member", auth.Create)

	if err != nil {

		return err
	}

	if check {

		uvuid := (uuid.New()).String()

		var member TblMember

		member.Uuid = uvuid

		member.MemberGroupId, _ = strconv.Atoi(c.PostFormValue("mem_group"))

		member.FirstName = c.PostFormValue("mem_name")

		member.LastName = c.PostFormValue("mem_lname")

		member.Email = c.PostFormValue("mem_email")

		member.MobileNo = c.PostFormValue("mem_mobile")

		member.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))

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
func (a Memberauth) UpdateMember(c *http.Request) error {

	userid, _, checkerr := auth.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := auth.Authority.IsGranted(a.Authority, "Member", auth.Update)

	if err != nil {

		return err
	}

	if check {

		uvuid := (uuid.New()).String()

		var member TblMember

		member.Uuid = uvuid

		member.MemberGroupId, _ = strconv.Atoi(c.PostFormValue("mem_group"))

		member.FirstName = c.PostFormValue("mem_name")

		member.LastName = c.PostFormValue("mem_lname")

		member.Email = c.PostFormValue("mem_email")

		member.MobileNo = c.PostFormValue("mem_mobile")

		member.IsActive, _ = strconv.Atoi(c.PostFormValue("mem_activestat"))

		member.ModifiedBy = userid

		member.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		if err := a.Authority.DB.Model(&member).Updates(TblMember{FirstName: member.FirstName, IsActive: member.IsActive,
			LastName: member.LastName, ModifiedOn: member.ModifiedOn, ModifiedBy: member.ModifiedBy, MobileNo: member.MobileNo, MemberGroupId: member.MemberGroupId}).Error; err != nil {

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

	check, err := auth.Authority.IsGranted(a.Authority, "Member", auth.Delete)

	if err != nil {

		return err
	}

	if check {

		var member TblMember

		member.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		member.DeletedBy = userid

		if err := a.Authority.DB.Model(&member).Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

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
func Process(r *http.Request, field string) (*Image, error) {
	file, info, err := r.FormFile(field)

	if err != nil {
		return nil, err
	}

	contentType := info.Header.Get("Content-Type")

	if !okContentType(contentType) {
		return nil, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	_, _, err = image.Decode(bytes.NewReader(bs))

	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    info.Filename,
		ContentType: contentType,
		Data:        bs,
		Size:        len(bs),
	}

	return i, nil
}
