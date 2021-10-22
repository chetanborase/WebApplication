package database

import (
	"WebApplication/internal/storage/model"
	"fmt"
	"strings"
)

var (
	insertUser = `insert  into Tbl_UserInfo (UserUUID, IdXref, FirstName, MiddleName, LastName, FullName, Email, UserName, Password, CreatedAt)
    VALUE  (?,?,?,?,?,?,?,?,?,?)`
	insertAddress = `insert into Tbl_AddressInfo (UserUUID, Address1, Address2, Address3, Area, City, State, Country, PinCode)
 VALUE (?,?,?,?,?,?,?,?,?)`
	insertContact = `insert into Tbl_ContactInfo (UserUUID, SocialMediaID, 
                         WebSite, DialCode, 
                         PhoneNumber, FullPhoneNumber)
                         value (?,?,?,?,?,?)`
	lockTable = `lock tables Tbl_UserInfo write,
            Tbl_AddressInfo write,
            Tbl_ContactInfo write`
	unlockTable = `unlock tables`
	updateUser  = `
update Tbl_UserInfo
SET IdXref     =?,
    FirstName  = ?,
    MiddleName = ?,
    LastName   = ?,
    FullName   = ?,
    Email      = ?,
    Password   = ?,
    UpdatedAt  =?
where UserUUID = ? AND DeletedAt is null `
	updateAddress = `update Tbl_AddressInfo
set   Address1       = ?,
      Address2       = ?,
      Address3       = ?,
      Area           = ?,
      City           = ?,
      State          = ?,
      Country        = ?,
      PinCode        = ?
    where UserUUID = ?;
`
	updateContact = `update Tbl_ContactInfo set
      SocialMediaID  = ?,
      WebSite        = ?,
      DialCode       = ?,
      PhoneNumber    = ?,
      FullPhoneNumber= ?
where UserUUID =?`
	deleteUser = `update Tbl_UserInfo set DeletedAt = ? where UserUUID = ? AND DeletedAt is null`
)

func getSearchQuery(values model.UrlRequest) (string, []interface{}) {
	var (
		queryExtender []string
		queryArgs     []interface{}
	)
	if values.Key != nil {
		queryExtender = append(queryExtender, SearchWithKeyCondition)
		key := values.Key[0]
		keys := []interface{}{key, key, key, key}
		queryArgs = append(queryArgs, keys...)
	}
	if values.City != nil {
		extender, args := extendQuery("City", values.City)
		queryExtender = append(queryExtender, " ("+strings.Join(extender, " OR ")+") ")
		queryArgs = append(queryArgs, args...)
	}
	if values.State != nil {
		extender, args := extendQuery("State", values.State)
		queryExtender = append(queryExtender, " ("+strings.Join(extender, " OR ")+") ")
		queryArgs = append(queryArgs, args...)
	}
	if values.Country != nil {
		extender, args := extendQuery("Country", values.Country)
		queryExtender = append(queryExtender, " ("+strings.Join(extender, " OR ")+") ")
		queryArgs = append(queryArgs, args...)
	}
	if values.LastId == 0 {
		values.LastId = values.Page * values.Limit
	}

	//extender, args := extendQuery("u.id >", []string{strconv.Itoa(values.LastId)})
	//queryExtender = append(queryExtender, " ("+strings.Join(extender, " OR ")+") ")
	queryArgs = append(queryArgs, values.LastId)
	conditionKeyword := ""
	if len(queryArgs) > 0 {
		conditionKeyword = " where "
	}
	//log.Fatal("value is  :  ", strings.Join(queryExtender, " AND "))
	query := fmt.Sprint(`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
		u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
		AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
		CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
		from Tbl_UserInfo as u
		inner join Tbl_AddressInfo AI
		on u.UserUUID = AI.UserUUID
		inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID`,
		conditionKeyword, strings.Join(queryExtender, " AND "), " AND u.id > ?", " limit ? ")
	limit := values.Limit
	queryArgs = append(queryArgs, limit)
	return query, queryArgs
}

func extendQuery(segmentName string, value []string) (query []string, args []interface{}) {
	for _, v := range value {
		query = append(query, segmentName+" = ? ")
		args = append(args, v)
	}
	return
}
