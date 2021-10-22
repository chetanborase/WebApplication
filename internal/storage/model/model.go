package model

import (
	"net/url"
	"strconv"
	"time"
)

type (
	User struct {
		UserUUID   string    `json:"user_uuid"`
		IdXref     string    `json:"id_xref"`
		FirstName  string    `json:"first_name"`
		MiddleName string    `json:"middle_name"`
		LastName   string    `json:"last_name"`
		FullName   string    `json:"full_name"`
		Email      string    `json:"email"`
		UserName   string    `json:"user_name"`
		Password   string    `json:"password"`
		UpdatedAt  time.Time `json:"updated_at"`
		Addr       Address
		Cont       Contact
	}
	Address struct {
		AddressUserUUID string `json:"address_user_uuid"`
		Address1        string `json:"address_1"`
		Address2        string `json:"address_2"`
		Address3        string `json:"address_3"`
		Area            string `json:"area"`
		City            string `json:"city"`
		State           string `json:"state"`
		Country         string `json:"country"`
		PinCode         string `json:"pin_code"`
	}
	Contact struct {
		ContactUserUUID string `json:"contact_user_uuid"`
		SocialMediaID   string `json:"social_media_id"`
		WebSite         string `json:"web_site"`
		DialCode        string `json:"dial_code"`
		PhoneNumber     string `json:"phone_number"`
		FullPhoneNumber string `json:"full_phone_number"`
	}
)

type UrlRequest struct {
	Key     []string
	City    []string
	State   []string
	Country []string
	Limit   int
	Page    int
	LastId  int
}

func MapUrlToStruct(values url.Values) UrlRequest {
	var limit, page, lastId int
	pageLimit, found := values["limit"]
	if !found {
		limit = 10
	} else {
		limit, _ = strconv.Atoi(pageLimit[0])
	}
	pageNumber, found := values["page"]
	if !found {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageNumber[0])
	}
	id, found := values["lastid"]
	if !found {
		lastId = (page - 1) * limit
	} else {
		lastId, _ = strconv.Atoi(id[0])
	}
	return UrlRequest{
		Key:     values["key"],
		City:    values["city"],
		State:   values["state"],
		Country: values["country"],
		Limit:   limit,
		Page:    page,
		LastId:  lastId,
	}
}
