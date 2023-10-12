package memberaccess

import (
	"github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/member"
	"gorm.io/gorm"
)

type AccessAuth struct {
	Authority auth.Authority
}

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

	err = GetSpaceByMemberId(&TblaccessControlusergroup, Groupid)

	if err != nil {

		return []int{}, err
	}

	var spid []int

	for _, val := range TblaccessControlusergroup {

		spid = append(spid, val.SpacesId...)
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

	err = GetPageByMemberId(&TblaccessControlusergroup, groupid)

	if err != nil {

		return []int{}, err
	}

	var pageids []int

	for _, val := range TblaccessControlusergroup {

		pageids = append(pageids, val.PagesId...)
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

	err = GetGroupByMemberId(&TblaccessControlusergroup, groupid)

	if err != nil {

		return []int{}, err
	}

	var pgroupid []int

	for _, val := range TblaccessControlusergroup {

		pgroupid = append(pgroupid, val.GroupsId...)
	}

	return pgroupid, nil

}
