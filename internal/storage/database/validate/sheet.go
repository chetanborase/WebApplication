package validate

import (
	"WebApplication/internal/storage/database/model"
	"errors"
	"github.com/google/uuid"
	"strings"
)

func SheetData(row []string) (u model.User, err error) {
	errStringForRow := ""
	if l := len(row[0]); l < 1 || l > model.TblUserInfoSizeOfIdXref {
		errStringForRow += "SizeOf FirstName Should be between 1 to 200 character"
	}
	if l := len(row[1]); l < 1 || l > model.TblUserInfoSizeOfFirstName {
		errStringForRow += "SizeOf FirstName Should be between 1 to 200 character"
	}
	if l := len(row[2]); l < 1 || l > model.TblUserInfoSizeOfMiddleName {
		errStringForRow += "SizeOf MiddleName Should be between 1 to 200 character"
	}
	if l := len(row[3]); l < 1 || l > model.TblUserInfoSizeOfLastName {
		errStringForRow += "SizeOf LastName Should be between 1 to 200 character"
	}
	if l := len(row[4]); l < 1 || l > model.TblUserInfoSizeOfFullName {
		errStringForRow += "SizeOf FullName Should be between 1 to 200 character"
	}
	if l := len(row[5]); l < 1 || l > model.TblContactInfoSizeOfEmail {
		errStringForRow += "SizeOf FullName Should be between 1 to 200 character"
	}
	if l := len(row[6]); l < 1 || l > model.TblUserInfoSizeOfUserName {
		errStringForRow += "SizeOf FullName Should be between 1 to 200 character"
	}
	if l := len(row[8]); l < 1 || l > model.TblAddressInfoSizeOfAddress1 {
		errStringForRow += "SizeOf Address Should be between 1 to 200 character"
	}
	if l := len(row[12]); l < 1 || l > model.TblAddressInfoSizeOfCity {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if l := len(row[13]); l < 1 || l > model.TblAddressInfoSizeOfState {
		errStringForRow += "SizeOf DialCode Should be between 1 to 200 character"
	}
	if l := len(row[14]); l < 1 || l > model.TblAddressInfoSizeOfCountry {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if l := len(row[15]); l < 1 || l > model.TblAddressInfoSizeOfPinCode {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if l := len(row[18]); l < 1 || l > model.TblContactInfoSizeOfDialCode {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if l := len(row[19]); l < 1 || l > model.TblContactInfoSizeOfPhoneNumber {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if l := len(row[20]); l < 1 || l > model.TblContactInfoSizeOfFullPhoneNumber {
		errStringForRow += "SizeOf PhoneNumber Should be between 1 to 200 character"
	}
	if len(errStringForRow) > 0 {
		return model.User{}, errors.New(errStringForRow)
	}
	return MapXlsToStruct(row), nil
}
func MapXlsToStruct(row []string) (u model.User) {
	u = model.User{}
	u.UserUUID = new(string)
	*u.UserUUID = uuid.New().String()
	u.IdXref = &row[0]
	u.FirstName = &row[1]
	u.MiddleName = &row[2]
	u.LastName = &row[3]
	u.FullName = &row[4]
	u.Email = &row[5]
	u.UserName = &row[6]
	u.Password = &row[7]
	u.Address1 = &row[8]
	u.Address2 = &row[9]
	u.Address3 = &row[10]
	u.Area = &row[11]
	u.City = &row[12]
	u.State = &row[13]
	u.Country = &row[14]
	u.PinCode = &row[15]
	u.SocialMediaID = &row[16]
	u.WebSite = &row[17]
	u.DialCode = &row[18]
	u.PhoneNumber = &row[19]
	u.FullPhoneNumber = &row[20]
	u.AddressUserUUID = u.UserUUID
	u.ContactUserUUID = u.UserUUID
	return
}

func SheetHeader(row []string) (err error) {
	if strings.Compare(row[0], "IdXref") != 0 {
		goto ifErr
	}
	if row[1] != "FirstName" {
		goto ifErr
	}
	if row[2] != "MiddleName" {
		goto ifErr
	}
	if row[3] != "LastName" {
	}
	if row[4] != "FullName" {
		goto ifErr
	}
	if row[5] != "Email" {
		goto ifErr
	}
	if row[6] != "UserName" {
		goto ifErr
	}
	if row[7] != "Password" {
		goto ifErr
	}
	if row[8] != "Address1" {
		goto ifErr
	}
	if row[9] != "Address2" {
		goto ifErr
	}
	if row[10] != "Address3" {
		goto ifErr
	}
	if row[11] != "Area" {
		goto ifErr
	}
	if row[12] != "City" {
		goto ifErr
	}
	if row[13] != "State" {
		goto ifErr
	}
	if row[14] != "Country" {
		goto ifErr
	}
	if row[15] != "PinCode" {
		goto ifErr
	}
	if row[16] != "SocialMediaID" {
		goto ifErr
	}
	if row[17] != "WebSite" {
		goto ifErr
	}
	if row[18] != "DialCode" {
		goto ifErr
	}
	if row[19] != "PhoneNumber" {
		goto ifErr
	}
	if row[20] != "FullPhoneNumber" {
		goto ifErr
	}

	return nil

ifErr:
	err = errors.New("header Validation Failed")
	return

}
