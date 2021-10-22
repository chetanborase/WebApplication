package database

import (
	"WebApplication/internal/generator"
	"WebApplication/internal/storage/database/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

var (
	user model.User
)

func TestInsertUser(t *testing.T) {
	user = testGetUserModel()
	//expected no error as data is fresh and unique
	err := InsertUser(user)
	assert.NoError(t, err, "Insert User")

	//expected error because data already stored in above part and not unique anymore
	err = InsertUser(user)
	assert.Error(t, err, err)

	//user id should not be nil, in below case error expected - "non-empty user uuid required"
	u := testGetUserModel()
	u.UserUUID = nil
	assert.Error(t, errors.New("non-empty user uuid required"), InsertUser(u))

	// Change UserUUID,AddressUUID,ContactUUID
	// Cannot add or update a child row: a foreign key constraint fails (`accounts`.`Tbl_AddressInfo`, CONSTRAINT `Tbl_AddressInfo_ibfk_1` FOREIGN KEY (`UserUUID`) REFERENCES `Tbl_UserInfo` (`UserUUID`))"}
	u = testGetUserModel()
	u.UserUUID = generator.GetRandomString(15)
	u.AddressUserUUID = generator.GetRandomString(15)
	err = InsertUser(u)
	log.Println(err)
	if err == nil {
		t.Error(err)
		t.FailNow()
	} else {
		if !strings.Contains(err.Error(), "foreign key constraint fails") {
			t.FailNow()
		}
	}
}
func TestGetUser(t *testing.T) {
	_, err := GetUser("some random string to ensure user not exist with this string")
	assert.EqualError(t, err, "sql: no rows in result set")

	//below there is valid user uuid and active account
	//no error expected
	u, err := GetUser(*user.UserUUID)
	assert.NoError(t, err)
	if !testCheckIfUserStructEquals(u,user){
		t.Error("both are not equal")
		t.FailNow()
	}

	//below there is valid user uuid but user account has been disabled(soft deleted)
	//error expected - "User Not Found"
	_, err = GetUser("0002cb9d-b02d-459a-ab47-c3a3244b8dab")
	assert.EqualError(t, err, "sql: no rows in result set")
}

func TestUpdateUser(t *testing.T) {
	//changing some data
	user.FirstName = generator.GetStringPointer("Ramesh")
	user.MiddleName = generator.GetStringPointer("Suresh")

	//no error expected as we are sure that user is stored in database
	err := UpdateUser(user)
	assert.NoError(t, err, err)

	//check if record really updated or not
	u, _ := GetUser(*user.UserUUID)
	if !testCheckIfUserStructEquals(u,user){
		t.Error("record actually not updated in database")
		t.FailNow()
	}

	//let update user account that not even exist
	//error expected - "records not found"
	u = testGetUserModel()
	err = UpdateUser(u)
	assert.EqualError(t, err, "records not found")

}

func TestDeleteUser(t *testing.T) {
	//delete active user account
	//no error expected
	err := DeleteUser(*user.UserUUID)
	assert.NoError(t, err, err)

	//delete deleted User Account
	//error expected :no account found with given user id
	err = DeleteUser("0002cb9d-b02d-459a-ab47-c3a3244b8dab")
	assert.EqualError(t, err, "no account found with given user id")

	//delete an account that not exist
	//error expected :no account found with given user id
	err = DeleteUser("vnbifujsbvfjsvbfjdbnvjdvbf;k")
	assert.EqualError(t, err, "no account found with given user id")

}
func TestUpdateUser2(t *testing.T) {
	err := UpdateUser(user)
	assert.EqualError(t, err, "records not found")
}


func testGetUserModel() model.User {
	u := model.User{
		UserUUID:   generator.GetRandomString(15),
		IdXref:     generator.GetRandomString(15),
		FirstName:  generator.GetRandomString(15),
		MiddleName: generator.GetRandomString(15),
		LastName:   generator.GetRandomString(15),
		FullName:   generator.GetRandomString(15),
		Email:      generator.GetRandomString(15),
		UserName:   generator.GetRandomString(15),
		Password:   generator.GetRandomString(15),
		Address: model.Address{
			Address1: generator.GetRandomString(15),
			Address2: generator.GetRandomString(15),
			Address3: generator.GetRandomString(15),
			Area:     generator.GetRandomString(15),
			City:     generator.GetRandomString(15),
			State:    generator.GetRandomString(15),
			Country:  generator.GetRandomString(15),
			PinCode:  generator.GetRandomString(15),
		},
		Contact: model.Contact{
			SocialMediaID:   generator.GetRandomString(15),
			WebSite:         generator.GetRandomString(15),
			DialCode:        generator.GetRandomString(15),
			PhoneNumber:     generator.GetRandomString(15),
			FullPhoneNumber: generator.GetRandomString(15),
		},
	}
	u.AddressUserUUID = u.UserUUID
	u.ContactUserUUID = u.UserUUID
	return u
}
func testCheckIfUserStructEquals(actual model.User, expected model.User) bool {

	if *actual.UserUUID == *expected.UserUUID &&
		*actual.IdXref == *expected.IdXref &&
		*actual.FirstName == *expected.FirstName &&
		*actual.MiddleName == *expected.MiddleName &&
		*actual.LastName == *expected.LastName &&
		*actual.FullName == *expected.FullName &&
		*actual.Email == *expected.Email &&
		*actual.UserName == *expected.UserName &&
		*actual.Address1 == *expected.Address1 &&
		*actual.Address2 == *expected.Address2 &&
		*actual.Address3 == *expected.Address3 &&
		*actual.Area == *expected.Area &&
		*actual.City == *expected.City &&
		*actual.State == *expected.State &&
		*actual.Country == *expected.Country &&
		*actual.PinCode == *expected.PinCode &&
		*actual.SocialMediaID == *expected.SocialMediaID &&
		*actual.WebSite == *expected.WebSite &&
		*actual.DialCode == *expected.DialCode &&
		*actual.PhoneNumber == *expected.PhoneNumber &&
		*actual.FullPhoneNumber == *expected.FullPhoneNumber {
		return true
	} else {
		return false
	}
}
