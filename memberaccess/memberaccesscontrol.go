package memberaccess

import (
	"errors"
	"log"

	"github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/member"
	"gorm.io/gorm"
)

type AccessAuth struct {
	Authority auth.Authorization
}

type AccessType struct{}

var AT AccessType

func MigrateTables(db *gorm.DB) {
	db.AutoMigrate(
		&TblAccessControl{},
		&TblAccessControlPages{},
		&TblAccessControlUserGroup{},
	)
}

/*Get  Space Id */
func (a AccessAuth) GetSpace() (spaceid []int, err error) {

	_, Groupid, err := member.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if err != nil {

		return []int{}, err
	}

	var TblaccessControlusergroup []TblAccessControlUserGroup

	err = AT.GetSpaceByMemberId(&TblaccessControlusergroup, Groupid, a.Authority.DB)

	if err != nil {

		return []int{}, err
	}

	var spid []int

	for _, val := range TblaccessControlusergroup {

		spid = append(spid, val.SpacesId)
	}

	return spid, nil
}

/*Get Page Id*/
func (a AccessAuth) GetPage() (pageid []int, err error) {

	_, groupid, err := member.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if err != nil {

		return []int{}, err
	}

	var TblaccessControlusergroup []TblAccessControlUserGroup

	err = AT.GetPageByMemberId(&TblaccessControlusergroup, groupid, a.Authority.DB)

	if err != nil {

		return []int{}, err
	}

	var pageids []int

	for _, val := range TblaccessControlusergroup {

		pageids = append(pageids, val.PageId)
	}

	return pageids, nil

}

/*Get Page Id*/
func (a AccessAuth) GetGroup() (pagegroupid []int, err error) {

	_, groupid, err := member.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if err != nil {

		return []int{}, err
	}

	var TblaccessControlusergroup []TblAccessControlUserGroup

	err = AT.GetGroupByMemberId(&TblaccessControlusergroup, groupid, a.Authority.DB)

	if err != nil {

		return []int{}, err
	}

	var pgroupid []int

	for _, val := range TblaccessControlusergroup {

		pgroupid = append(pgroupid, val.PageGroupId)
	}

	return pgroupid, nil

}

/*Check LoginPage*/
func (a AccessAuth) CheckPageLogin(pageid int) (bool, error) {

	var page []TblAccessControlUserGroup

	AT.CheckPageRestrict(&page, pageid, a.Authority.DB)

	_, groupid, err := member.VerifyToken(a.Authority.Token, a.Authority.Secret)

	if err != nil {

		log.Println(err)

	}

	var loginflg = false

	var MemberNot bool

	for _, val := range page {

		if groupid == 0 && val.PageId == pageid {

			loginflg = true

			break

		}

	}

	for _, val := range page {

		if groupid == val.MemberGroupId && val.PageId == pageid {

			MemberNot = true

			break

		} else {

			MemberNot = false

		}

	}

	if len(page) > 0 {

		if loginflg {

			return false, errors.New("login required")
		}

		if !MemberNot {

			return false, errors.New("not permitted")
		}

	}
	return true, nil
}
