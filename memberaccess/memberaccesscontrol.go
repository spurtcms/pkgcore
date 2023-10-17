package memberaccess

import (
	"github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/member"
	"gorm.io/gorm"
)

type AccessAuth struct {
	Authority auth.Authority
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

var Access AccessAuth

func AccessInstance(a *auth.Option) auth.Authority {

	auth := auth.Authority{
		DB:     a.DB,
		Token:  a.Token,
		Secret: a.Secret,
	}

	MigrateTables(auth.DB)

	return auth
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
