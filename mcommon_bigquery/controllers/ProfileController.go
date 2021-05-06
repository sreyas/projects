package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func GetProfileDetails() {
	tableID := "Profiles"
	ctx := context.Background()

	bigqueryclient, err := bigquery.NewClient(ctx, ProjectID)

	defer bigqueryclient.Close()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	tableFound := false

	ts := bigqueryclient.Dataset(DatasetID).Tables(ctx)
	for {
		t, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

		}
		if t.TableID == tableID {
			tableFound = true
		}

	}
	if !tableFound {
		ProfileSchema := bigquery.Schema{
			{Name: "ProfileId", Type: bigquery.StringFieldType, Required: false},
			{Name: "FirstName", Type: bigquery.StringFieldType, Required: false},
			{Name: "LastName", Type: bigquery.StringFieldType, Required: false},
			{Name: "PhoneNumber", Type: bigquery.StringFieldType, Required: false},
			{Name: "Email", Type: bigquery.StringFieldType, Required: false},
			{Name: "Status", Type: bigquery.StringFieldType, Required: false},
			{Name: "CreatedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "UpdatedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "OptedOutAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "OptedOutSource", Type: bigquery.StringFieldType, Required: false},
			{Name: "Street1", Type: bigquery.StringFieldType, Required: false},
			{Name: "Street2", Type: bigquery.StringFieldType, Required: false},
			{Name: "City", Type: bigquery.StringFieldType, Required: false},
			{Name: "State", Type: bigquery.StringFieldType, Required: false},
			{Name: "PostalCode", Type: bigquery.StringFieldType, Required: false},
			{Name: "Country", Type: bigquery.StringFieldType, Required: false},
			{Name: "Latitude", Type: bigquery.StringFieldType, Required: false},
			{Name: "Longitude", Type: bigquery.StringFieldType, Required: false},
			{Name: "Precision", Type: bigquery.StringFieldType, Required: false},
			{Name: "LastSavedCity", Type: bigquery.StringFieldType, Required: false},
			{Name: "LastSavedState", Type: bigquery.StringFieldType, Required: false},
			{Name: "LastSavedPostalCode", Type: bigquery.StringFieldType, Required: false},
			{Name: "LastSavedCountry", Type: bigquery.StringFieldType, Required: false},
			{Name: "CongressionalDistrict", Type: bigquery.StringFieldType, Required: false},
			{Name: "StateUpperDistrict", Type: bigquery.StringFieldType, Required: false},
			{Name: "StateLowerDistrict", Type: bigquery.StringFieldType, Required: false},
			{Name: "SplitDistrict", Type: bigquery.StringFieldType, Required: false},
			{Name: "ConstituentID", Type: bigquery.StringFieldType, Required: false},
			{Name: "IntegrationType", Type: bigquery.StringFieldType, Required: false},
			{Name: "SynchronizedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "CustomColumns", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: ProfileSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(tableID)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}

	startDate, toDate := GetTableQueryParams(tableID)

	num := 1000
	page := 1
	for num == 1000 {
		url := "https://secure.mcommons.com/api/profiles?limit=1000&page=" + strconv.Itoa(page) + "&from=" + startDate + "&to=" + toDate
		page += 1
		log.Println(url)
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.SetBasicAuth(Username, Password)
		start := time.Now()
		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		res.Body.Close()
		profileResponse := ProfileTotalResponse{}
		_ = xml.Unmarshal([]byte(body), &profileResponse)
		elapsed := time.Since(start).Seconds()
		log.Println(url, "elapsed time", elapsed)
		num, _ = strconv.Atoi(profileResponse.Profiles.Num)
		go ProcessProfile(profileResponse, bigqueryclient, ctx, tableID)
	}

}
func ProcessProfile(profileResponse ProfileTotalResponse, bigqueryclient *bigquery.Client, ctx context.Context, tableID string) {
	TotalProfiles := []Profile{}
	for _, v := range profileResponse.Profiles.Profile {
		var singleProfile Profile
		singleProfile.ProfileId = v.ID
		singleProfile.FirstName = v.FirstName
		singleProfile.LastName = v.LastName
		singleProfile.PhoneNumber = v.PhoneNumber
		singleProfile.Email = v.Email
		singleProfile.Status = v.Status
		singleProfile.CreatedAt = v.CreatedAt
		singleProfile.UpdatedAt = v.UpdatedAt
		singleProfile.OptedOutAt = v.OptedOutAt
		singleProfile.OptedOutSource = v.OptedOutSource
		singleProfile.Street1 = v.Address.Street1
		singleProfile.Street2 = v.Address.Street2
		singleProfile.City = v.Address.City
		singleProfile.State = v.Address.State
		singleProfile.PostalCode = v.Address.PostalCode
		singleProfile.Country = v.Address.Country
		singleProfile.Latitude = v.LastSavedLocation.Latitude
		singleProfile.Longitude = v.LastSavedLocation.Longitude
		singleProfile.Precision = v.LastSavedLocation.Precision
		singleProfile.LastSavedCity = v.LastSavedLocation.City
		singleProfile.LastSavedState = v.LastSavedLocation.State
		singleProfile.LastSavedPostalCode = v.LastSavedLocation.PostalCode
		singleProfile.LastSavedCountry = v.LastSavedLocation.Country
		singleProfile.CongressionalDistrict = v.LastSavedDistricts.CongressionalDistrict
		singleProfile.StateUpperDistrict = v.LastSavedDistricts.StateUpperDistrict
		singleProfile.StateLowerDistrict = v.LastSavedDistricts.StateLowerDistrict
		singleProfile.SplitDistrict = v.LastSavedDistricts.SplitDistrict
		singleProfile.ConstituentID = v.Integrations.Integration.ConstituentID
		singleProfile.IntegrationType = v.Integrations.Integration.Type
		singleProfile.SynchronizedAt = v.Integrations.Integration.SynchronizedAt
		for ci, c := range v.CustomColumns.CustomColumn {
			if ci != 0 {
				singleProfile.CustomColumns += ","
			}
			singleProfile.CustomColumns += c.Name
		}
		existingProfile := getExistingProfiles(ctx, bigqueryclient, singleProfile.ProfileId)
		if _, ok := existingProfile[singleProfile.ProfileId]; !ok {
			TotalProfiles = append(TotalProfiles, singleProfile)
		} else {
			log.Println("founddddddddddddd profileeeeeeee")
			TotalProfiles = append(TotalProfiles, singleProfile)
			updateqry := GenrerateProfileUpdateQuery(singleProfile)
			res, err := bigqueryclient.Query(updateqry).Read(ctx)
			log.Println(err)
			log.Println(res)
		}
		inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
		if err := inserter.Put(ctx, TotalProfiles); err != nil {
			log.Println("insertion error")
			log.Println(err)
		}
	}
	inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
	if err := inserter.Put(ctx, TotalProfiles); err != nil {
		log.Println("insertion error")
		log.Println(err)
	}
}

type Profile struct {
	ProfileId             string
	FirstName             string
	LastName              string
	PhoneNumber           string
	Email                 string
	Status                string
	CreatedAt             string
	UpdatedAt             string
	OptedOutAt            string
	OptedOutSource        string
	Street1               string
	Street2               string
	City                  string
	State                 string
	PostalCode            string
	Country               string
	Latitude              string
	Longitude             string
	Precision             string
	LastSavedCity         string
	LastSavedState        string
	LastSavedPostalCode   string
	LastSavedCountry      string
	CongressionalDistrict string
	StateUpperDistrict    string
	StateLowerDistrict    string
	SplitDistrict         string
	ConstituentID         string
	IntegrationType       string
	SynchronizedAt        string
	CustomColumns         string
}
type ProfileResponse struct {
	Text           string `xml:",chardata"`
	ID             string `xml:"id,attr"`
	FirstName      string `xml:"first_name"`
	LastName       string `xml:"last_name"`
	PhoneNumber    string `xml:"phone_number"`
	Email          string `xml:"email"`
	Status         string `xml:"status"`
	CreatedAt      string `xml:"created_at"`
	UpdatedAt      string `xml:"updated_at"`
	OptedOutAt     string `xml:"opted_out_at"`
	OptedOutSource string `xml:"opted_out_source"`
	Source         struct {
		Text        string `xml:",chardata"`
		Type        string `xml:"type,attr"`
		Name        string `xml:"name,attr"`
		ID          string `xml:"id,attr"`
		OptInPathID string `xml:"opt_in_path_id,attr"`
		MessageID   string `xml:"message_id,attr"`
		Email       string `xml:"email,attr"`
	} `xml:"source"`
	Address struct {
		Text       string `xml:",chardata"`
		Street1    string `xml:"street1"`
		Street2    string `xml:"street2"`
		City       string `xml:"city"`
		State      string `xml:"state"`
		PostalCode string `xml:"postal_code"`
		Country    string `xml:"country"`
	} `xml:"address"`
	LastSavedLocation struct {
		Text       string `xml:",chardata"`
		Latitude   string `xml:"latitude"`
		Longitude  string `xml:"longitude"`
		Precision  string `xml:"precision"`
		City       string `xml:"city"`
		State      string `xml:"state"`
		PostalCode string `xml:"postal_code"`
		Country    string `xml:"country"`
	} `xml:"last_saved_location"`
	LastSavedDistricts struct {
		Text                  string `xml:",chardata"`
		CongressionalDistrict string `xml:"congressional_district"`
		StateUpperDistrict    string `xml:"state_upper_district"`
		StateLowerDistrict    string `xml:"state_lower_district"`
		SplitDistrict         string `xml:"split_district"`
	} `xml:"last_saved_districts"`
	CustomColumns struct {
		Text         string `xml:",chardata"`
		CustomColumn []struct {
			Text      string `xml:",chardata"`
			Name      string `xml:"name,attr"`
			CreatedAt string `xml:"created_at,attr"`
			UpdatedAt string `xml:"updated_at,attr"`
		} `xml:"custom_column"`
	} `xml:"custom_columns"`
	Subscriptions struct {
		Text         string `xml:",chardata"`
		Subscription []struct {
			Text                string `xml:",chardata"`
			CampaignID          string `xml:"campaign_id,attr"`
			CampaignName        string `xml:"campaign_name,attr"`
			CampaignDescription string `xml:"campaign_description,attr"`
			OptInPathID         string `xml:"opt_in_path_id,attr"`
			Status              string `xml:"status,attr"`
			OptInSource         string `xml:"opt_in_source,attr"`
			CreatedAt           string `xml:"created_at,attr"`
			ActivatedAt         string `xml:"activated_at,attr"`
			OptedOutAt          string `xml:"opted_out_at,attr"`
			OptOutSource        string `xml:"opt_out_source,attr"`
		} `xml:"subscription"`
	} `xml:"subscriptions"`
	Integrations struct {
		Text        string `xml:",chardata"`
		Integration struct {
			Text           string `xml:",chardata"`
			ConstituentID  string `xml:"constituent_id,attr"`
			Type           string `xml:"type,attr"`
			SynchronizedAt string `xml:"synchronized_at,attr"`
		} `xml:"integration"`
	} `xml:"integrations"`
	Clicks struct {
		Text  string `xml:",chardata"`
		Click []struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id,attr"`
			CreatedAt   string `xml:"created_at"`
			URL         string `xml:"url"`
			RemoteAddr  string `xml:"remote_addr"`
			HTTPReferer string `xml:"http_referer"`
			UserAgent   string `xml:"user_agent"`
		} `xml:"click"`
	} `xml:"clicks"`
}
type ProfileTotalResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Success  string   `xml:"success,attr"`
	Profiles struct {
		Text    string            `xml:",chardata"`
		Num     string            `xml:"num,attr"`
		Page    string            `xml:"page,attr"`
		Profile []ProfileResponse `xml:"profile"`
	} `xml:"profiles"`
}

func getExistingProfiles(ctx context.Context, bigqueryclient *bigquery.Client, profileId string) map[string]Profile {
	tableId := "Profiles"
	query := bigqueryclient.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableId + ` where ProfileId="` + profileId + `"	`)
	bigqryRes, err := query.Read(ctx)
	log.Println(err)
	Profiles := make(map[string]Profile)
	for {
		var row Profile
		err := bigqryRes.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		Profiles[row.ProfileId] = row

	}
	return Profiles

}
func GenrerateProfileUpdateQuery(profile Profile) string {
	updateqry := ""
	updateqry += fmt.Sprintf("update .%s.%s set ", DatasetID, "Profiles")
	v := reflect.ValueOf(profile)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		limiter := ","
		if i == v.NumField()-1 {
			limiter = ""
		}
		updateqry += fmt.Sprintf(` %s="%s" %s`, typeOfS.Field(i).Name, v.Field(i).Interface(), limiter)
	}
	updateqry += fmt.Sprintf(` where ProfileId="%s" `, profile.ProfileId)
	return updateqry
}
