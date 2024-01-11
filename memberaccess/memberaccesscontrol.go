package memberaccess

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/pkgcontent/channels"
	"github.com/spurtcms/pkgcore/auth"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

type AccessAuth struct {
	Authority auth.Authorization
}

type AccessAdminAuth struct {
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

/*List */
func (access AccessAdminAuth) ContentAccessList(limit int, offset int, filter Filter) (accesslist []TblAccessControl, totalCount int64, err error) {

	_, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return []TblAccessControl{}, 0, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return []TblAccessControl{}, 0, err
	}

	if check {

		var contentAccessList []TblAccessControl

		AT.GetContentAccessList(&contentAccessList, limit, offset, filter, access.Authority.DB)

		var final_content_accesslist []TblAccessControl

		for _, contentAccess := range contentAccessList {

			var access_grant_memgrps []TblAccessControlUserGroup

			AT.GetAccessGrantedMemberGroups(&access_grant_memgrps, contentAccess.Id, access.Authority.DB)

			for _, memgrp := range access_grant_memgrps {

				if memgrp.AccessControlId == contentAccess.Id {

					var memberGroup member.TblMemberGroup

					AT.GetMemberGroupsByContentAccessMemId(&memberGroup, memgrp.MemberGroupId, access.Authority.DB)

					contentAccess.MemberGroups = append(contentAccess.MemberGroups, memberGroup)
				}
			}

			var entriesCount, pageCount int64

			AT.GetaccessGrantedEntriesCount(&entriesCount, contentAccess.Id, access.Authority.DB)

			AT.GetaccessGrantedPageCount(&pageCount, contentAccess.Id, access.Authority.DB)

			if entriesCount > 0 {

				contentAccess.AccessGrantedModules = append(contentAccess.AccessGrantedModules, "Channel")
			}

			if pageCount > 0 {

				contentAccess.AccessGrantedModules = append(contentAccess.AccessGrantedModules, "Space")
			}

			if !contentAccess.ModifiedOn.IsZero() {

				contentAccess.DateString = contentAccess.ModifiedOn.UTC().Format("02 Jan 2006 03:04 PM")

			} else {

				contentAccess.DateString = contentAccess.CreatedOn.UTC().Format("02 Jan 2006 03:04 PM")

			}

			final_content_accesslist = append(final_content_accesslist, contentAccess)

		}

		var contentAccessList1 []TblAccessControl

		_, content_access_count := AT.GetContentAccessList(&contentAccessList1, 0, 0, filter, access.Authority.DB)

		return final_content_accesslist, content_access_count, nil

	}

	return []TblAccessControl{}, 0, errors.New("not authorized")
}

/*Get Access by id*/
func (access AccessAdminAuth) GetControlAccessById(accessid int) (accesslist TblAccessControl, pg []Page, spage []SubPage, pgroup []PageGroup, selectedspacesid []int, MembergroupIds []int, channelid []int, channelEntries []Entry, err error) {

	_, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return TblAccessControl{}, []Page{}, []SubPage{}, []PageGroup{}, []int{}, []int{}, []int{}, []Entry{}, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return TblAccessControl{}, []Page{}, []SubPage{}, []PageGroup{}, []int{}, []int{}, []int{}, []Entry{}, err
	}

	if check {

		var AccessControl TblAccessControl

		AT.GetContentAccessByAccessId(&AccessControl, accessid, access.Authority.DB)

		var pages []Page

		var subpages []SubPage

		var pagegroups []PageGroup

		var contentAccessPages []TblAccessControlPages

		AT.GetPagesAndPageGroupsInContentAccess(&contentAccessPages, accessid, access.Authority.DB)

		var pageArrContainer [][]TblPage

		var pgArrContainer []TblAccessControlPages

		seen1 := make(map[int]bool)

		seen2 := make(map[int]bool)

		seen3 := make(map[int]bool)

		seen4 := make(map[int]bool)

		for _, pagz := range contentAccessPages {

			if pagz.ParentPageId == 0 {

				if !seen1[pagz.PageId] {

					var pg Page

					pg.Id = strconv.Itoa(pagz.PageId)

					pg.GroupId = strconv.Itoa(pagz.PageGroupId)

					pg.SpaceId = strconv.Itoa(pagz.SpacesId)

					pages = append(pages, pg)

					seen1[pagz.PageId] = true

				}
			}

			if pagz.ParentPageId != 0 {

				if !seen2[pagz.PageId] {

					var spg SubPage

					spg.Id = strconv.Itoa(pagz.PageId)

					spg.GroupId = strconv.Itoa(pagz.PageGroupId)

					spg.ParentId = strconv.Itoa(pagz.ParentPageId)

					spg.SpaceId = strconv.Itoa(pagz.SpacesId)

					subpages = append(subpages, spg)

					seen2[pagz.PageId] = true
				}

			}

			if pagz.PageGroupId != 0 {

				if !seen3[pagz.PageGroupId] {

					var pagesinPgg []TblPage

					AT.GetPagesUnderPageGroup(&pagesinPgg, pagz.PageGroupId, access.Authority.DB)

					pageArrContainer = append(pageArrContainer, pagesinPgg)

					seen3[pagz.PageGroupId] = true

				}

				if !seen4[pagz.PageId] {

					pgArrContainer = append(pgArrContainer, pagz)

					seen4[pagz.PageId] = true

				}

			}

		}

		// log.Println("orgpgg", pageArrContainer)

		// log.Println("pggchk", pgArrContainer)

		groupedObjects := make(map[int][]TblAccessControlPages)

		OriginalPgg := make(map[int][]TblPage)

		for _, pgz := range pgArrContainer {

			if _, exists := groupedObjects[pgz.PageGroupId]; !exists {

				groupedObjects[pgz.PageGroupId] = []TblAccessControlPages{}

			}

			groupedObjects[pgz.PageGroupId] = append(groupedObjects[pgz.PageGroupId], pgz)
		}

		for _, pggArr := range pageArrContainer {

			for _, pages := range pggArr {

				OriginalPgg[pages.PageGroupId] = pggArr

			}

		}

		for pggId, array := range OriginalPgg {

			for pggid := range groupedObjects {

				if len(OriginalPgg[pggId]) == len(groupedObjects[pggid]) && pggId == pggid {

					for index, result := range array {

						if index == 0 {

							var pgg PageGroup

							pgg.Id = strconv.Itoa(result.PageGroupId)

							pgg.SpaceId = strconv.Itoa(result.SpacesId)

							pagegroups = append(pagegroups, pgg)

							break
						}

					}

				}
			}
		}

		// var access_grant_memgrps_list []int

		var accessGrantedMemgrps []int

		AT.GetAccessGrantedMemberGroupsList(&accessGrantedMemgrps, accessid, access.Authority.DB)

		var spaceIds []int

		AT.GetContentAccessSpaces(&spaceIds, accessid, access.Authority.DB)

		var accessSpaceIds []int

		for _, spaceId := range spaceIds {

			var contentAccessPages []int

			AT.GetcontentAccessPagesBySpaceId(&contentAccessPages, spaceId, accessid, access.Authority.DB)

			var tblPageData []TblPage

			AT.GetPagesUnderSpaces(&tblPageData, spaceId, access.Authority.DB)

			if len(tblPageData) == len(contentAccessPages) {

				accessSpaceIds = append(accessSpaceIds, spaceId)
			}

		}

		var channelEntries []Entry

		var contentAccessEntries []TblAccessControlPages

		AT.GetAccessGrantedEntries(&contentAccessEntries, accessid, access.Authority.DB)

		channelMap := make(map[int][]TblAccessControlPages)

		for _, accessEntry := range contentAccessEntries {

			chanEntry := Entry{Id: strconv.Itoa(accessEntry.EntryId), ChannelId: strconv.Itoa(accessEntry.ChannelId)}

			channelEntries = append(channelEntries, chanEntry)

			if _, exists := channelMap[accessEntry.ChannelId]; !exists {

				channelMap[accessEntry.ChannelId] = []TblAccessControlPages{}

			}

			channelMap[accessEntry.ChannelId] = append(channelMap[accessEntry.ChannelId], accessEntry)

		}

		var channelIds []int

		for channelId, entriesArr := range channelMap {

			var entriesCountInChannel int64

			AT.GetEntriesCountUnderChannel(&entriesCountInChannel, channelId, access.Authority.DB)

			if int(entriesCountInChannel) == len(entriesArr) {

				channelIds = append(channelIds, channelId)
			}
		}

		return AccessControl, pages, subpages, pagegroups, accessGrantedMemgrps, accessSpaceIds, channelIds, channelEntries, nil

	}

	return TblAccessControl{}, []Page{}, []SubPage{}, []PageGroup{}, []int{}, []int{}, []int{}, []Entry{}, errors.New("not authorized")
}

/*Create Access control*/
func (access AccessAdminAuth) CreateMemberAccessControl(control MemberAccessControlRequired) (bool, error) {

	userid, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var contentAccess TblAccessControl

		contentAccess.AccessControlName = control.Title

		contentAccess.AccessControlSlug = strings.ToLower(control.Title)

		contentAccess.CreatedBy = userid

		contentAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		contentAccess.IsDeleted = 0

		contentAccessId, err := AT.NewContentAccessEntry(&contentAccess, access.Authority.DB)

		if err != nil {

			log.Println(err)

		}

		for _, memgrp_id := range control.MemberGroupIds {

			var memberGrpAccess TblAccessControlUserGroup

			memberGrpAccess.AccessControlId = contentAccessId.Id

			memberGrpAccess.MemberGroupId = memgrp_id

			memberGrpAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			memberGrpAccess.CreatedBy = userid

			memberGrpAccess.IsDeleted = 0

			err = AT.GrantAccessToMemberGroups(&memberGrpAccess, access.Authority.DB)

			if err != nil {

				log.Println(err)
			}

		}

		var memberGrpAccess []TblAccessControlUserGroup

		AT.GetMemberGrpByAccessControlId(&memberGrpAccess, contentAccessId.Id, access.Authority.DB)

		var pagesNotInSpaces []TblAccessControlPages

		seen_page := make(map[int]bool)

		for _, memgrp := range memberGrpAccess {

			for _, pg := range control.Pages {

				pageId, _ := strconv.Atoi(pg.Id)

				var page_access TblAccessControlPages

				page_access.AccessControlUserGroupId = memgrp.Id

				page_access.SpacesId, _ = strconv.Atoi(pg.SpaceId)

				page_access.PageGroupId, _ = strconv.Atoi(pg.GroupId)

				page_access.PageId = pageId

				page_access.CreatedBy = userid

				page_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				page_access.IsDeleted = 0

				err = AT.InsertPageEntries(&page_access, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}

				if !seen_page[pageId] {

					pagesNotInSpaces = append(pagesNotInSpaces, page_access)

					seen_page[pageId] = true
				}

			}

			for _, subpage := range control.SubPage {

				spgId, _ := strconv.Atoi(subpage.Id)

				var spg_access TblAccessControlPages

				spg_access.AccessControlUserGroupId = memgrp.Id

				spg_access.SpacesId, _ = strconv.Atoi(subpage.SpaceId)

				spg_access.PageGroupId, _ = strconv.Atoi(subpage.GroupId)

				spg_access.PageId = spgId

				spg_access.CreatedBy = userid

				spg_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				spg_access.IsDeleted = 0

				err = AT.InsertPageEntries(&spg_access, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}

				if !seen_page[spgId] {

					pagesNotInSpaces = append(pagesNotInSpaces, spg_access)

					seen_page[spgId] = true
				}

			}

			for _, entry := range control.ChannelEntries {

				chanId, _ := strconv.Atoi(entry.ChannelId)

				entryId, _ := strconv.Atoi(entry.Id)

				var channelAccess TblAccessControlPages

				channelAccess.AccessControlUserGroupId = memgrp.Id

				channelAccess.ChannelId = chanId

				channelAccess.EntryId = entryId

				channelAccess.CreatedBy = userid

				channelAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				channelAccess.IsDeleted = 0

				err = AT.InsertPageEntries(&channelAccess, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}
			}

		}

		var tblPagesData [][]TblPage

		if len(pagesNotInSpaces) != 0 {

			for _, page := range pagesNotInSpaces {

				for _, space_id := range control.SpacesIds {

					if page.SpacesId != space_id {

						var tblpagedata []TblPage

						err = AT.GetPagesUnderSpaces(&tblpagedata, space_id, access.Authority.DB)

						if err != nil {

							log.Println(err)
						}

						tblPagesData = append(tblPagesData, tblpagedata)

					}

				}
			}

		} else {

			for _, space_id := range control.SpacesIds {

				var tblpagedata []TblPage

				err = AT.GetPagesUnderSpaces(&tblpagedata, space_id, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}

				tblPagesData = append(tblPagesData, tblpagedata)

			}

		}

		for _, memgrp := range memberGrpAccess {

			for _, pagedatas := range tblPagesData {

				for _, pagedata := range pagedatas {

					var page_access TblAccessControlPages

					page_access.AccessControlUserGroupId = memgrp.Id

					page_access.SpacesId = pagedata.SpacesId

					page_access.PageGroupId = pagedata.PageGroupId

					page_access.PageId = pagedata.Id

					page_access.CreatedBy = userid

					page_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					page_access.IsDeleted = 0

					err = AT.InsertPageEntries(&page_access, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}
				}

			}
		}

		return true, nil

	}

	return false, errors.New("not authorized")
}

/* update */
func (access AccessAdminAuth) UpdateMemberAccessControl(control MemberAccessControlRequired, accesscontrolid int) (bool, error) {

	userid, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var contentAccess TblAccessControl

		contentAccess.Id = accesscontrolid

		contentAccess.AccessControlName = control.Title

		contentAccess.AccessControlSlug = strings.ToLower(control.Title)

		contentAccess.ModifiedBy = userid

		contentAccess.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err = AT.UpdateContentAccessId(&contentAccess, access.Authority.DB)

		if err != nil {

			log.Println(err)

		}

		for _, memgrp_id := range control.MemberGroupIds {

			var access_count int64

			err = AT.CheckPresenceOfAccessGrantedMemberGroups(&access_count, memgrp_id, accesscontrolid, access.Authority.DB)

			if err != nil {

				log.Println(err)

			}

			var memberGrpAccess TblAccessControlUserGroup

			memberGrpAccess.AccessControlId = accesscontrolid

			memberGrpAccess.MemberGroupId = memgrp_id

			if access_count == 0 {

				memberGrpAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				memberGrpAccess.CreatedBy = userid

				memberGrpAccess.IsDeleted = 0

				err = AT.GrantAccessToMemberGroups(&memberGrpAccess, access.Authority.DB)

				if err != nil {

					log.Println(err)

				}

			} else if access_count == 1 {

				memberGrpAccess.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				memberGrpAccess.ModifiedBy = userid

				err = AT.UpdateContentAccessMemberGroup(&memberGrpAccess, access.Authority.DB)

				if err != nil {

					log.Println(err)

				}

			}

		}

		// var access_grant_memgrps_list []int

		// for _, memgrpid := range control.MemberGroupIds {

		// 	memId, _ := strconv.Atoi(memgrpid)

		// 	access_grant_memgrps_list = append(access_grant_memgrps_list, memId)

		// }

		var MemGrpAccess []TblAccessControlUserGroup

		AT.GetMemberGrpByAccessControlId(&MemGrpAccess, accesscontrolid, access.Authority.DB)

		var pagesNotInSpaces []TblAccessControlPages

		var pageIds []int

		var entryIds []int

		seen_page := make(map[int]bool)

		seen_entry := make(map[int]bool)

		for _, memgrp := range MemGrpAccess {

			for _, pg := range control.Pages {

				var page_count int64

				page_id, _ := strconv.Atoi(pg.Id)

				group_id, _ := strconv.Atoi(pg.GroupId)

				space_id, _ := strconv.Atoi(pg.SpaceId)

				err = AT.CheckPresenceOfPageInContentAccess(memgrp.Id, page_id, group_id, space_id, &page_count, access.Authority.DB)

				// log.Println("pg_count", page_count)

				if err != nil {

					log.Println(err)

				}

				var page_access TblAccessControlPages

				page_access.AccessControlUserGroupId = memgrp.Id

				page_access.SpacesId, _ = strconv.Atoi(pg.SpaceId)

				page_access.PageGroupId, _ = strconv.Atoi(pg.GroupId)

				page_access.PageId, _ = strconv.Atoi(pg.Id)

				if page_count == 0 {

					page_access.CreatedBy = userid

					page_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					page_access.IsDeleted = 0

					err = AT.InsertPageEntries(&page_access, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				} else if page_count == 1 {

					page_access.ModifiedBy = userid

					page_access.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					err = AT.UpdatePagesInContentAccess(&page_access, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				}

				if !seen_page[page_id] {

					pageIds = append(pageIds, page_id)

					pagesNotInSpaces = append(pagesNotInSpaces, page_access)

					seen_page[page_id] = true
				}

			}

			for _, subpage := range control.SubPage {

				var spg_count int64

				page_id, _ := strconv.Atoi(subpage.Id)

				group_id, _ := strconv.Atoi(subpage.GroupId)

				space_id, _ := strconv.Atoi(subpage.SpaceId)

				err = AT.CheckPresenceOfPageInContentAccess(memgrp.Id, page_id, group_id, space_id, &spg_count, access.Authority.DB)

				// log.Println("spg_count", spg_count)

				if err != nil {

					log.Println(err)

				}

				var spg_access TblAccessControlPages

				spg_access.AccessControlUserGroupId = memgrp.Id

				spg_access.SpacesId, _ = strconv.Atoi(subpage.SpaceId)

				spg_access.PageGroupId, _ = strconv.Atoi(subpage.GroupId)

				spg_access.PageId, _ = strconv.Atoi(subpage.Id)

				if spg_count == 0 {

					spg_access.CreatedBy = userid

					spg_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					spg_access.IsDeleted = 0

					err = AT.InsertPageEntries(&spg_access, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				} else if spg_count == 1 {

					spg_access.ModifiedBy = userid

					spg_access.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					err = AT.UpdatePagesInContentAccess(&spg_access, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				}

				if !seen_page[page_id] {

					pageIds = append(pageIds, page_id)

					pagesNotInSpaces = append(pagesNotInSpaces, spg_access)

					seen_page[page_id] = true
				}

			}

			for _, entry := range control.ChannelEntries {

				chanId, _ := strconv.Atoi(entry.ChannelId)

				entryId, _ := strconv.Atoi(entry.Id)

				var entryCount int64

				err = AT.CheckPresenceOfChannelEntriesInContentAccess(&entryCount, memgrp.Id, chanId, entryId, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}

				var channelAccess TblAccessControlPages

				channelAccess.AccessControlUserGroupId = memgrp.Id

				channelAccess.ChannelId = chanId

				channelAccess.EntryId = entryId

				if entryCount == 0 {

					channelAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					channelAccess.CreatedBy = userid

					channelAccess.IsDeleted = 0

					err = AT.InsertPageEntries(&channelAccess, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				} else if entryCount == 1 {

					channelAccess.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					channelAccess.ModifiedBy = userid

					err = AT.UpdateChannelEntriesInContentAccess(&channelAccess, access.Authority.DB)

					if err != nil {

						log.Println(err)
					}

				}

				if !seen_entry[entryId] {

					entryIds = append(entryIds, entryId)

					seen_entry[entryId] = true

				}
			}

		}

		var tblPagesData [][]TblPage

		if len(pagesNotInSpaces) != 0 {

			for _, page := range pagesNotInSpaces {

				for _, space_id := range control.SpacesIds {

					if page.SpacesId != space_id {

						var tblPagedata []TblPage

						err = AT.GetPagesUnderSpaces(&tblPagedata, space_id, access.Authority.DB)

						if err != nil {

							log.Println(err)
						}

						tblPagesData = append(tblPagesData, tblPagedata)

					}

				}
			}

		} else {

			for _, space_id := range control.SpacesIds {

				var tblPagedata []TblPage

				err = AT.GetPagesUnderSpaces(&tblPagedata, space_id, access.Authority.DB)

				if err != nil {

					log.Println(err)
				}

				tblPagesData = append(tblPagesData, tblPagedata)

			}

		}

		for _, memgrp := range MemGrpAccess {

			for _, pagedatas := range tblPagesData {

				for _, pagedata := range pagedatas {

					var page_count int64

					err = AT.CheckPresenceOfPageInContentAccess(memgrp.Id, pagedata.Id, pagedata.PageGroupId, pagedata.SpacesId, &page_count, access.Authority.DB)

					if err != nil {

						log.Println(err)

					}

					var page_access TblAccessControlPages

					page_access.AccessControlUserGroupId = memgrp.Id

					page_access.SpacesId = pagedata.SpacesId

					page_access.PageGroupId = pagedata.PageGroupId

					page_access.PageId = pagedata.Id

					if page_count == 0 {

						page_access.CreatedBy = userid

						page_access.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

						page_access.IsDeleted = 0

						err = AT.InsertPageEntries(&page_access, access.Authority.DB)

						if err != nil {

							log.Println(err)
						}

					} else if page_count == 1 {

						page_access.ModifiedBy = userid

						page_access.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

						err = AT.UpdatePagesInContentAccess(&page_access, access.Authority.DB)

						if err != nil {

							log.Println(err)
						}
					}

					if !seen_page[pagedata.Id] {

						pageIds = append(pageIds, pagedata.Id)

						seen_page[pagedata.Id] = true
					}

				}

			}
		}

		var memgrp_access1 TblAccessControlUserGroup

		memgrp_access1.IsDeleted = 1

		memgrp_access1.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		memgrp_access1.DeletedBy = userid

		AT.RemoveMemberGroupsNotUnderContentAccessRights(&memgrp_access1, control.MemberGroupIds, accesscontrolid, access.Authority.DB)

		for _, memgrp := range MemGrpAccess {

			var pg_access1 TblAccessControlPages

			pg_access1.AccessControlUserGroupId = memgrp.Id

			pg_access1.IsDeleted = 1

			pg_access1.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			pg_access1.DeletedBy = userid

			err = AT.RemovePagesNotUnderContentAccess(&pg_access1, pageIds, access.Authority.DB)

			if err != nil {

				log.Println(err)
			}

			err = AT.RemoveChannelEntriesNotUnderContentAccess(&pg_access1, entryIds, access.Authority.DB)

			if err != nil {

				log.Println(err)
			}

		}

		return true, nil

	}

	return false, errors.New("not authorized")

}

/*Delete member access control*/
func (access AccessAdminAuth) DeleteMemberAccessControl(accesscontrolid int) (bool, error) {

	userid, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var accesscontrol TblAccessControl

		accesscontrol.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		accesscontrol.DeletedBy = userid

		accesscontrol.IsDeleted = 1

		err := AT.DeleteControl(&accesscontrol, accesscontrolid, access.Authority.DB)

		var acusergrp TblAccessControlUserGroup

		acusergrp.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		acusergrp.DeletedBy = userid

		acusergrp.IsDeleted = 1

		AT.DeleteInAccessUserGroup(&acusergrp, accesscontrolid, access.Authority.DB)

		var accessgrp []TblAccessControlUserGroup

		AT.GetDeleteIdInAccessUserGroup(&accessgrp, accesscontrolid, access.Authority.DB)

		var pgid []int

		for _, v := range accessgrp {

			pgid = append(pgid, v.Id)

		}

		var accesscontrolpg TblAccessControlPages

		accesscontrolpg.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		accesscontrolpg.DeletedBy = userid

		accesscontrolpg.IsDeleted = 1

		AT.DeleteAccessControlPages(&accesscontrolpg, pgid, access.Authority.DB)

		if err != nil {

			return false, err

		}

		return true, nil

	}
	return false, errors.New("not authorized")
}

func (access AccessAdminAuth) GetChannelsWithEntries() ([]channels.TblChannel, error) {

	_, _, checkerr := auth.VerifyToken(access.Authority.Token, access.Authority.Secret)

	if checkerr != nil {

		return []channels.TblChannel{}, checkerr
	}

	check, err := access.Authority.IsGranted("Member-Restrict", auth.CRUD)

	if err != nil {

		return []channels.TblChannel{}, err
	}

	if check {

		var channel_contents []channels.TblChannel

		err := AT.GetChannels(&channel_contents, access.Authority.DB)

		if err != nil {

			log.Println(err)

			return []channels.TblChannel{}, err
		}

		var FinalChannellist []channels.TblChannel

		for _, channel := range channel_contents {

			var channel_entries []channels.TblChannelEntries

			AT.GetChannelEntriesByChannelId(&channel_entries, channel.Id, access.Authority.DB)

			if len(channel_entries) > 0 {

				channel.ChannelEntries = channel_entries

				FinalChannellist = append(FinalChannellist, channel)

			}
		}

		return FinalChannellist, nil

	}

	return []channels.TblChannel{}, errors.New("not authorized")

}
