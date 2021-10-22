package database

import (
	"WebApplication/internal/storage/model"
	"testing"
)

func TestGetSearchQuery(t *testing.T) {
	inputs := []model.UrlRequest{
		{
			Key:     []string{"ab"},
			City:    []string{"vapi"},
			State:   nil,
			Country: nil,
			Limit:   10,
			Page:    0,
			LastId:  0,
		},
		{
			Key:     nil,
			City:    []string{"vapi,valsad"},
			State:   []string{"Gujarat"},
			Country: nil,
			Limit:   10,
			Page:    1,
			LastId:  100,
		},
		{
			Key:     nil,
			City:    []string{"vapi,valsad"},
			State:   []string{"Gujarat"},
			Country: nil,
			Limit:   10,
			Page:    1,
			LastId:  100,
		},
		{
			Key:     nil,
			City:    []string{"vapi,valsad"},
			State:   []string{"Gujarat"},
			Country: nil,
			Limit:   10,
			Page:    1,
			LastId:  100,
		},
	}
	expectedOutput := []string{
		`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
		u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
		AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
		CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
		from Tbl_UserInfo as u
		inner join Tbl_AddressInfo AI
		on u.UserUUID = AI.UserUUID
		inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID
		where (City= ? OR City= ?) AND State= ? limit 10'`,

		`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
                u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
                AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
                CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
                from Tbl_UserInfo as u
                inner join Tbl_AddressInfo AI
                on u.UserUUID = AI.UserUUID
                inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID where (FirstName like concat('%%',?,'%%') OR
                        MiddleName like concat('%%',?,'%%') OR
                        LastName like concat('%%',?,'%%') OR
                        FullName like concat('%%',?,'%%')) AND  (City = ? )  limit ?  [ab ab ab ab vapi 10]`,
		`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
                u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
                AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
                CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
                from Tbl_UserInfo as u
                inner join Tbl_AddressInfo AI
                on u.UserUUID = AI.UserUUID
                inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID where (FirstName like concat('%%',?,'%%') OR
                        MiddleName like concat('%%',?,'%%') OR
                        LastName like concat('%%',?,'%%') OR
                        FullName like concat('%%',?,'%%')) limit ?  [ab ab ab ab 10]`,
		`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
                u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
                AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
                CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
                from Tbl_UserInfo as u
                inner join Tbl_AddressInfo AI
                on u.UserUUID = AI.UserUUID
                inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID where (FirstName like concat('%%',?,'%%') OR
                        MiddleName like concat('%%',?,'%%') OR
                        LastName like concat('%%',?,'%%') OR
                        FullName like concat('%%',?,'%%')) AND  (City = ? )  limit ?  [ab ab ab ab vapi 10]`,

		`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
                u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
                AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
                CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
                from Tbl_UserInfo as u
                inner join Tbl_AddressInfo AI
                on u.UserUUID = AI.UserUUID
                inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID where (FirstName like concat('%%',?,'%%') OR
                        MiddleName like concat('%%',?,'%%') OR
                        LastName like concat('%%',?,'%%') OR
                        FullName like concat('%%',?,'%%')) AND  (City = ?  OR City = ? )  limit ?  [ab ab ab ab vapi navsari 10]`,
	}
	_, _ = inputs, expectedOutput
}

func getStringPointer(str string) *string {
	return &str
}
