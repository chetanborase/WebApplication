package model

import (
	"WebApplication/internal/storage/model"
	"time"
)

const (
	TblUserInfoSizeOfIdXref     = 200
	TblUserInfoSizeOfFirstName  = 200
	TblUserInfoSizeOfMiddleName = 200
	TblUserInfoSizeOfLastName   = 200
	TblUserInfoSizeOfFullName   = 500
	TblUserInfoSizeOfUserName   = 200

	TblAddressInfoSizeOfAddress1 = 200
	TblAddressInfoSizeOfCity     = 200
	TblAddressInfoSizeOfCountry  = 200
	TblAddressInfoSizeOfState    = 200
	TblAddressInfoSizeOfPinCode  = 200

	TblContactInfoSizeOfPhoneNumber     = 200
	TblContactInfoSizeOfDialCode        = 200
	TblContactInfoSizeOfFullPhoneNumber = 200
	TblContactInfoSizeOfEmail           = 200
)

type (
	//todo store multiple address and contact per user
	User struct {
		Id         *int       `json:"-"`
		UserUUID   *string    `db:"UserUUID"`
		IdXref     *string    `db:"IdXref"`
		FirstName  *string    `db:"FirstName"`
		MiddleName *string    `db:"MiddleName"`
		LastName   *string    `db:"LastName"`
		FullName   *string    `db:"FullName"`
		Email      *string    `db:"Email"`
		UserName   *string    `db:"UserName"`
		Password   *string    `db:"Password"`
		CreatedAt  time.Time  `db:"CreatedAt"`
		UpdatedAt  time.Time  `db:"UpdatedAT"`
		DeletedAt  *time.Time `db:"DeletedAt"`
		Address
		Contact
	}
	Address struct {
		AddressUserUUID *string `db:"UserUUID"`
		Address1        *string `db:"Address1"`
		Address2        *string `db:"Address2"`
		Address3        *string `db:"Address3"`
		Area            *string `db:"Area"`
		City            *string `db:"City"`
		State           *string `db:"State"`
		Country         *string `db:"Country"`
		PinCode         *string `db:"PinCode"`
	}
	Contact struct {
		ContactUserUUID *string `db:"UserUUID"`
		SocialMediaID   *string `db:"SocialMediaID"`
		WebSite         *string `db:"WebSite"`
		DialCode        *string `db:"DialCode"`
		PhoneNumber     *string `db:"PhoneNumber"`
		FullPhoneNumber *string `db:"FullPhoneNumber"`
	}
	SearchResult struct {
		Users      []User `json:"users"`
		PageNo     int    `json:"page_no"`
		TotalPages int    `json:"total_pages"`
		IsLastPage bool   `json:"is_last_page"`
		LastID     int    `json:"max_id"`
		Time       string `json:"time"`
	}
)

func MapToDbModel(u model.User) User {
	dbUser := User{
		IdXref:     safePointer(u.IdXref),
		UserUUID:   safePointer(u.UserUUID),
		FirstName:  safePointer(u.FirstName),
		MiddleName: safePointer(u.MiddleName),
		LastName:   safePointer(u.LastName),
		FullName:   safePointer(u.FullName),
		Email:      safePointer(u.Email),
		UserName:   safePointer(u.UserName),
		Password:   safePointer(u.Password),
		UpdatedAt:  u.UpdatedAt,

		Address: Address{
			AddressUserUUID: safePointer(u.UserUUID),
			Address1:        safePointer(u.Addr.Address1),
			Address2:        safePointer(u.Addr.Address2),
			Address3:        safePointer(u.Addr.Address3),
			Area:            safePointer(u.Addr.Area),
			City:            safePointer(u.Addr.City),
			State:           safePointer(u.Addr.State),
			Country:         safePointer(u.Addr.Country),
			PinCode:         safePointer(u.Addr.PinCode),
		},
		Contact: Contact{
			ContactUserUUID: safePointer(u.UserUUID),
			SocialMediaID:   safePointer(u.Cont.SocialMediaID),
			WebSite:         safePointer(u.Cont.WebSite),
			DialCode:        safePointer(u.Cont.DialCode),
			PhoneNumber:     safePointer(u.Cont.PhoneNumber),
			FullPhoneNumber: safePointer(u.Cont.FullPhoneNumber),
		},
	}
	return dbUser
}
func safePointer(str string) *string{
	if str==""{
		return nil
	}
	return &str

}