/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type serverConfig struct {
	CCID    string
	Address string
}
type SmartContract struct {
	contractapi.Contract
}

type commonData struct {
	UserType     string `json:"userType"`
	UserID       string `json:"userID"`
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
}

type AdminData struct {
	TotalNumberOfPatient                   int `json:"totalNumberOfPatient"`
	TotalNumberOfDoctor                    int `json:"totalNumberOfDoctor"`
	TotalNumberOfPathologist               int `json:"totalNumberOfPathologist"`
	TotalNumberOfPharmacyCompany           int `json:"totalNumberOfPharmacyCompany"`
	TotalNumberOfMedicalResearchLab        int `json:"totalNumberOfMedicalResearchLab"`
	TotalNumberOfPremiumPatient            int `json:"totalNumberOfPremiumPatient"`
	TotalNumberOfPremiumDoctor             int `json:"totalNumberOfPremiumDoctor"`
	TotalNumberOfPremiumPathologist        int `json:"totalNumberOfPremiumPathologist"`
	TotalNumberOfPremiumPharmacyCompany    int `json:"totalNumberOfPremiumPharmacyCompany"`
	TotalNumberOfPremiumMedicalResearchLab int `json:"totalNumberOfPremiumMedicalResearchLab"`
}

func (s *SmartContract) GetAllAdmindata(ctx contractapi.TransactionContextInterface) (AdminData, error) {

	adminDataJSON, err := ctx.GetStub().GetState("user-analytics")
	if err != nil {
		return AdminData{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if adminDataJSON == nil {
		return AdminData{}, nil // Return an empty AdminData struct
	}

	// Unmarshal the JSON data into an Admin object
	var adminData AdminData
	err = json.Unmarshal(adminDataJSON, &adminData)
	if err != nil {
		return AdminData{}, fmt.Errorf("failed to unmarshal admin data: %v", err)
	}
	return adminData, nil
}

type Admin struct {
	AdminID         string              `json:"adminID"`
	Name            string              `json:"name"`
	Gender          string              `json:"gender"`
	ProfilePic      string              `json:"profilePic"`
	AllPrescription map[string][]string `json:"allPrescription"`

	IsAdded               bool     `json:"isAdded"`
	UserType              string   `json:"userType"`
	Birthday              string   `json:"birthday"`
	EmailAddress          string   `json:"emailAddress"`
	Age                   int      `json:"age"`
	Location              string   `json:"location"`
	PatientToAdmin        []string `json:"patientToAdmin"`
	SharedAllUsersAddress []string `json:"sharedAllUsersAddress"`
}

func NewAdmin(adminID string, name string, gender string, birthday string, emailAddress string, location string) Admin {
	// Return the Admin object with default values for unspecified fields
	return Admin{
		AdminID:               adminID,
		Name:                  name,
		Gender:                gender,
		Birthday:              birthday,
		EmailAddress:          emailAddress,              // Default profile picture
		AllPrescription:       make(map[string][]string), // Default empty list for prescriptions
		IsAdded:               true,
		UserType:              "admin",    // Default user type
		Age:                   30,         // Default age
		Location:              location,   // Default location
		PatientToAdmin:        []string{}, // Default empty list for patient associations
		SharedAllUsersAddress: []string{}, // Default empty list for shared users
	}
}
func (s *SmartContract) SetAdmin(ctx contractapi.TransactionContextInterface, adminID string, name string, gender string, birthday string, emailAddress string, location string) error {
	// Create a new Admin object
	admin := NewAdmin(adminID, name, gender, birthday, emailAddress, location)

	// Marshal the admin object to JSON
	adminJSON, err := json.Marshal(admin)
	if err != nil {
		return fmt.Errorf("failed to marshal admin: %v", err)
	}

	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[adminID] = string(TAdmin)
	allAdmins, err := s.GetAllAdmins(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve all admins: %v", err)
	}
	allAdmins = append(allAdmins, adminID)
	if err := s.putState(ctx, AllAdminsKey, allAdmins); err != nil {
		return fmt.Errorf("failed to store  allAdmins update: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}

	return ctx.GetStub().PutState(adminID, adminJSON)

}

// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
//     patients := []Patient{
//         {
//             PatientID:           "patient1",
//             Name:                "John Doe",
//             Gender:              "Male",
//             Age:                 30,
//             Location:            "New York",
//             IsAdded:             true,
//             UserType:            "Regular",
//             ImgUrl:              []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
//             PatientPersonalData: PatientPersonalData{ /* Add appropriate fields */ },
//             ProfilePic:          "https://example.com/profile1.jpg",
//             Birthday:            "1993-01-15",
//             EmailAddress:        "john.doe@example.com",
//             SharedAllDoctorAddress: []string{"doctor1@example.com", "doctor2@example.com"},
//             PersonalDoctor:      []string{"doctor3@example.com"},
//         },
//         {
//             PatientID:           "patient2",
//             Name:                "Jane Smith",
//             Gender:              "Female",
//             Age:                 28,
//             Location:            "Los Angeles",
//             IsAdded:             true,
//             UserType:            "VIP",
//             ImgUrl:              []string{"https://example.com/image3.jpg"},
//             PatientPersonalData: PatientPersonalData{ /* Add appropriate fields */ },
//             ProfilePic:          "https://example.com/profile2.jpg",
//             Birthday:            "1995-03-20",
//             EmailAddress:        "jane.smith@example.com",
//             SharedAllDoctorAddress: []string{"doctor4@example.com"},
//             PersonalDoctor:      []string{"doctor5@example.com", "doctor6@example.com"},
//         },
//         // Add more patient records as needed
//     }

//     for _, patient := range patients {
//         patientJSON, err := json.Marshal(patient)
//         if err != nil {
//             return fmt.Errorf("failed to marshal patient: %v", err)
//         }

//         err = ctx.GetStub().PutState(patient.PatientID, patientJSON)
//         if err != nil {
//             return fmt.Errorf("failed to put patient into world state: %v", err)
//         }
//     }

//     return nil
// }

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	accounts := make(map[string]string)

	accounts["user1"] = string(TAdmin)
	accounts["user2"] = string(TPatient)

	accountsJSON, err := json.Marshal(accounts)
	if err != nil {
		return fmt.Errorf("failed to serialize accounts: %v", err)
	}

	err = ctx.GetStub().PutState("accounts", accountsJSON)
	if err != nil {
		return fmt.Errorf("failed to store accounts in the ledger: %v", err)
	}

	return nil
}
func (s *SmartContract) GetAccounts(ctx contractapi.TransactionContextInterface) (map[string]string, error) {
	accountsJSON, err := ctx.GetStub().GetState("accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve accounts from the ledger: %v", err)
	}
	if accountsJSON == nil {
		return nil, fmt.Errorf("accounts not found in the ledger")
	}

	var accounts map[string]string
	err = json.Unmarshal(accountsJSON, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize accounts: %v", err)
	}

	return accounts, nil
}

func (s *SmartContract) GetAdmin(ctx contractapi.TransactionContextInterface, id string) (Admin, error) {
	// Retrieve the admin data from the world state using the provided admin ID
	adminJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Admin{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if adminJSON == nil {
		return Admin{}, fmt.Errorf("admin with ID %s does not exist", id)
	}

	// Unmarshal the JSON data into an Admin object
	var admin Admin
	err = json.Unmarshal(adminJSON, &admin)
	if err != nil {
		return Admin{}, fmt.Errorf("failed to unmarshal admin data: %v", err)
	}
	return admin, nil
}

type PatientPersonalData struct {
	Height           float32 `json:"height"`
	Blood            string  `json:"blood"`
	PreviousDiseases string  `json:"previousDiseases"`
	Medicinedrugs    string  `json:"medicinedrugs"`
	BadHabits        string  `json:"badHabits"`
	ChronicDiseases  string  `json:"chronicDiseases"`
	HealthAllergies  string  `json:"healthAllergies"`
	BirthDefects     string  `json:"birthDefects"`
}

type Patient struct {
	PatientID              string              `json:"patientID"`
	Name                   string              `json:"name"`
	Gender                 string              `json:"gender"`
	Age                    int                 `json:"age"`
	Location               string              `json:"location"`
	IsAdded                bool                `json:"isAdded"`
	UserType               string              `json:"userType"`
	Prescriptions          map[string][]string `json:"prescriptions"`
	PatientPersonalData    PatientPersonalData `json:"patientPersonalData"`
	ProfilePic             string              `json:"profilePic"`
	Birthday               string              `json:"birthday"`
	EmailAddress           string              `json:"emailAddress"`
	SharedAllDoctorAddress []string            `json:"sharedAllDoctorAddress"`
	PersonalDoctor         []string            `json:"personalDoctor"`
}

func NewPatient(patientID string, name string, gender string, age int, location string, birthday string, emailAddress string) Patient {
	return Patient{
		PatientID:              patientID,
		Name:                   name,
		Gender:                 gender,
		Age:                    age,
		IsAdded:                true,
		UserType:               "patient",
		Location:               location,
		Birthday:               birthday,
		EmailAddress:           emailAddress,
		Prescriptions:          make(map[string][]string),
		SharedAllDoctorAddress: []string{},
		PersonalDoctor:         []string{},
	}
}

type Doctor struct {
	DoctorID            string              `json:"doctorID"`
	Name                string              `json:"name"`
	Specialty           string              `json:"specialty"`
	ConsultationFee     float32             `json:"consultationFee"`
	BMDCNumber          int                 `json:"bmdcNumber"`
	YearOfExperience    int                 `json:"yearOfExperience"`
	PatientToDoctor     []string            `json:"patientToDoctor"`
	PatientTest         []string            `json:"patientTest"`
	IsAdded             bool                `json:"isAdded"`
	TreatedPatient      []string            `json:"treatedPatient"`
	UserType            string              `json:"userType"`
	ProfilePic          string              `json:"profilePic"`
	Birthday            string              `json:"birthday"`
	UserData            map[string]UserData `json:"userData"`
	SenderPathologist   []string            `json:"senderPathologist"`
	ReceiverPathologist []string            `json:"receiverPathologist"`
	EmailAddress        string              `json:"emailAddress"`
}

// Constructor function with default values
func NewDoctor(doctorID string, name string, specialty string, consultationFee float32, BMDCNumber int, yearOfExperience int, birthday string, emailAddress string) Doctor {
	return Doctor{
		DoctorID:            doctorID,
		Name:                name,
		Specialty:           specialty,
		ConsultationFee:     consultationFee,
		BMDCNumber:          BMDCNumber,
		YearOfExperience:    yearOfExperience,
		PatientToDoctor:     []string{},            // Default to empty array
		PatientTest:         []string{},            // Default to empty array
		TreatedPatient:      []string{},            // Default to empty array
		IsAdded:             true,                  // Default to true
		UserType:            "doctor",              // Default userType
		SenderPathologist:   []string{},            // Default to empty array
		ReceiverPathologist: []string{},            // Default to empty array
		UserData:            map[string]UserData{}, // Default to empty map
		Birthday:            birthday,
		EmailAddress:        emailAddress,
	}
}

func (s *SmartContract) GetAllUserTypeData(ctx contractapi.TransactionContextInterface, userId string) (commonData, error) {
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return commonData{}, err
	}

	userType, exists := accounts[userId]
	if !exists {
		return commonData{}, fmt.Errorf("user ID %s does not exist in accounts", userId)
	}

	userData := commonData{}
	switch userType {
	case string(TDoctor):
		doctor, err := s.GetDoctor(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("doctor with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     doctor.UserType,
			UserID:       doctor.DoctorID,
			Username:     doctor.Name,
			EmailAddress: doctor.EmailAddress,
		}

		return userData, nil

	case string(TMedicalResearchLab):
		medLab, err := s.GetMedicalResearchLab(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("medical research lab with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     medLab.UserType,
			UserID:       medLab.LabID,
			Username:     medLab.Name,
			EmailAddress: medLab.EmailAddress,
		}

		return userData, nil

	case string(TPharmacyCompany):
		pharmacy, err := s.GetPharmacyCompany(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("pharmacy company with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     pharmacy.UserType,
			UserID:       pharmacy.CompanyID,
			Username:     pharmacy.Name,
			EmailAddress: pharmacy.EmailAddress,
		}

		return userData, nil

	case string(TPathologist):
		pathologist, err := s.GetPathologist(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("pathologist with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     pathologist.UserType,
			UserID:       pathologist.PathologistID,
			Username:     pathologist.Name,
			EmailAddress: pathologist.EmailAddress,
		}

		return userData, nil

	case string(TPatient):
		patient, err := s.GetPatient(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("patient with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     patient.UserType,
			UserID:       patient.PatientID,
			Username:     patient.Name,
			EmailAddress: patient.EmailAddress,
		}

		return userData, nil

	case string(TAdmin):
		admin, err := s.GetAdmin(ctx, userId)
		if err != nil {
			return commonData{}, fmt.Errorf("admin with ID %s not found", userId)
		}
		userData = commonData{
			UserType:     admin.UserType,
			UserID:       admin.AdminID,
			Username:     admin.Name,
			EmailAddress: admin.EmailAddress,
		}

		return userData, nil

	default:
		return commonData{}, fmt.Errorf("unknown user type for address %d", userId)
	}
}

func (s *SmartContract) AddDisease(ctx contractapi.TransactionContextInterface, disease string) error {
	diseaseDataJSON, err := ctx.GetStub().GetState("diseaseData")
	if err != nil {
		return fmt.Errorf("failed to retrieve disease data: %v", err)
	}

	var diseaseData map[string]bool
	if diseaseDataJSON != nil {
		err = json.Unmarshal(diseaseDataJSON, &diseaseData)
		if err != nil {
			return fmt.Errorf("failed to parse disease data: %v", err)
		}
	} else {
		diseaseData = make(map[string]bool)
	}

	// Check if the disease already exists
	if _, exists := diseaseData[disease]; exists {
		return fmt.Errorf("disease already exists: %s", disease)
	}

	// Add new disease
	diseaseData[disease] = true

	// Store updated disease data
	diseaseDataJSON, err = json.Marshal(diseaseData)
	if err != nil {
		return fmt.Errorf("failed to serialize disease data: %v", err)
	}
	return ctx.GetStub().PutState("diseaseData", diseaseDataJSON)
}

func (s *SmartContract) GetDiseaseNames(ctx contractapi.TransactionContextInterface) ([]string, error) {
	diseaseDataJSON, err := ctx.GetStub().GetState("diseaseData")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve disease data: %v", err)
	}

	if diseaseDataJSON == nil {
		return []string{}, nil
	}

	var diseaseData map[string]bool
	err = json.Unmarshal(diseaseDataJSON, &diseaseData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disease data: %v", err)
	}
	diseaseList := make([]string, 0, len(diseaseData))
	for disease := range diseaseData {
		diseaseList = append(diseaseList, disease)
	}

	return diseaseList, nil
}

func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, imgurl string, disease string, user1Id string, user2Id string) error {
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	user1Type, exists := accounts[user1Id]
	if !exists {
		return fmt.Errorf("user ID %s does not exist in accounts", user1Id)
	}

	switch user1Type {
	case string(TPatient):
		if user1Type == string(TPatient) && user2Id == "-1" {
			patient, err := s.GetPatient(ctx, user1Id)
			if err != nil {
				return fmt.Errorf("patient with ID %s not found", user1Id)
			}

			if imgurl == "" || patient.Prescriptions == nil || patient.Prescriptions[disease] == nil || !contains(patient.Prescriptions[disease], imgurl) {
				return fmt.Errorf("prescription not found for patient with ID %s and disease %s", patient.PatientID, disease)
			}

			patient.Prescriptions[disease] = removeElement(patient.Prescriptions[disease], imgurl)

			// Persist updated state to the ledger
			patientJSON, err := json.Marshal(patient)
			if err != nil {
				return fmt.Errorf("failed to marshal patient data: %v", err)
			}
			err = ctx.GetStub().PutState(user1Id, patientJSON)
			if err != nil {
				return fmt.Errorf("failed to update patient state: %v", err)
			}
		} else {
			return fmt.Errorf("invalid user type for user ")
		}

	case string(TDoctor):
		doctor, err := s.GetDoctor(ctx, user1Id)
		if err != nil {
			return fmt.Errorf("doctor with ID %s not found", user1Id)
		}
		user2Type, exists := accounts[user2Id]
		if !exists {
			return fmt.Errorf("user ID %s does not exist in accounts", user2Id)
		}
		if user2Type == string(TPatient) {
			patient, err := s.GetPatient(ctx, user2Id)
			if err != nil {
				return fmt.Errorf("patient with ID %s not found", user2Id)
			}

			if userData, exists := doctor.UserData[patient.PatientID]; exists {
				if imgurl == "" || !contains(userData.ImagesUrl, imgurl) {
					return fmt.Errorf("prescription not found for patient with ID %s", user2Id)
				}

				userData.ImagesUrl = removeElement(userData.ImagesUrl, imgurl)
			} else {
				return fmt.Errorf("user data for patient with ID %s not found in doctor records", user2Id)
			}

		} else if user2Type == string(TPathologist) {
			pathologist, err := s.GetPathologist(ctx, user2Id)
			if err != nil {
				return fmt.Errorf("pathologist with ID %s not found", user2Id)
			}

			if userData, exists := doctor.UserData[pathologist.PathologistID]; exists {
				if imgurl == "" || !contains(userData.ImagesUrl, imgurl) {
					return fmt.Errorf("prescription not found for pathologist with ID %s", user2Id)
				}

				userData.ImagesUrl = removeElement(userData.ImagesUrl, imgurl)
			} else {
				return fmt.Errorf("user data for pathologist with ID %s not found in doctor records", user2Id)
			}
		} else {
			return fmt.Errorf("invalid user type for user ID %s", user2Id)
		}

		// Persist updated state to the ledger
		doctorJSON, err := json.Marshal(doctor)
		if err != nil {
			return fmt.Errorf("failed to marshal doctor data: %v", err)
		}
		err = ctx.GetStub().PutState(user1Id, doctorJSON)
		if err != nil {
			return fmt.Errorf("failed to update doctor state: %v", err)
		}

	case string(TPathologist):
		pathologist, err := s.GetPathologist(ctx, user1Id)
		if err != nil {
			return fmt.Errorf("pathologist with ID %s not found", user1Id)
		}
		user2Type, exists := accounts[user2Id]
		if !exists {
			return fmt.Errorf("user ID %s does not exist in accounts", user2Id)
		}
		if user2Type == string(TDoctor) {
			doctor, err := s.GetDoctor(ctx, user2Id)
			if err != nil {
				return fmt.Errorf("doctor with ID %s not found", user2Id)
			}

			if userData, exists := pathologist.UserData[doctor.DoctorID]; exists {
				if imgurl == "" || !contains(userData.ImagesUrl, imgurl) {
					return fmt.Errorf("prescription not found for doctor with ID %s", user2Id)
				}

				userData.ImagesUrl = removeElement(userData.ImagesUrl, imgurl)
			} else {
				return fmt.Errorf("user data for doctor with ID %s not found in pathologist records", user2Id)
			}
		} else {
			return fmt.Errorf("invalid user type for user ID %s", user2Id)
		}

		// Persist updated state to the ledger
		pathologistJSON, err := json.Marshal(pathologist)
		if err != nil {
			return fmt.Errorf("failed to marshal pathologist data: %v", err)
		}
		err = ctx.GetStub().PutState(user1Id, pathologistJSON)
		if err != nil {
			return fmt.Errorf("failed to update pathologist state: %v", err)
		}

	default:
		return errors.New("unauthorized access")
	}

	return nil
}

func (s *SmartContract) RevokeAccessData(ctx contractapi.TransactionContextInterface, senderId string, ruserId string) error {
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	sUserType, senderExists := accounts[senderId]
	if !senderExists {
		return fmt.Errorf("sender address %s does not exist in accounts", sUserType)
	}

	// Check if the user exists in accounts
	rUserType, userExists := accounts[ruserId]
	if !userExists {
		return fmt.Errorf("ruser address %s does not exist in accounts", rUserType)
	}

	switch sUserType {
	case string(TPatient):
		patient, err := s.GetPatient(ctx, senderId)
		if err != nil {
			return fmt.Errorf("patient with address %d not found", senderId)
		}

		if rUserType == string(TDoctor) {
			doctor, err := s.GetDoctor(ctx, ruserId)
			if err != nil {
				return fmt.Errorf("doctor with address %d not found", ruserId)
			}

			// Remove senderAddress from PatientToDoctor
			doctor.PatientToDoctor = removeElement(doctor.PatientToDoctor, senderId)
			// Remove userAddress from sharedAllDoctorAddress
			patient.SharedAllDoctorAddress = removeElement(patient.SharedAllDoctorAddress, ruserId)
			// Delete doctor's state from ledger
			doctorJSON, err := json.Marshal(doctor)
			if err != nil {
				return fmt.Errorf("failed to marshal doctor data: %v", err)
			}
			err = ctx.GetStub().PutState(ruserId, doctorJSON)
			if err != nil {
				return fmt.Errorf("failed to update doctor state: %v", err)
			}
			patientsJSON, err := json.Marshal(patient)
			if err != nil {
				return fmt.Errorf("failed to marshal patient data: %v", err)
			}
			err = ctx.GetStub().PutState(senderId, patientsJSON)
			if err != nil {
				return fmt.Errorf("failed to update patient state: %v", err)
			}

		}

	case string(TAdmin):
		admin, err := s.GetAdmin(ctx, senderId)
		if err != nil {
			return fmt.Errorf("admin with address %d not found", senderId)
		}

		switch rUserType {
		case string(TMedicalResearchLab):
			medLab, err := s.GetMedicalResearchLab(ctx, ruserId)
			if err != nil {
				return fmt.Errorf("medical research lab with address %d not found", ruserId)
			}
			admin.SharedAllUsersAddress = removeElement(admin.SharedAllUsersAddress, senderId)
			// Remove senderAddress from adminToMedRcLab
			medLab.AdminToMedRcLab = removeElement(medLab.AdminToMedRcLab, senderId)

			// Delete medical research lab state from ledger
			medLabJSON, err := json.Marshal(medLab)
			if err != nil {
				return fmt.Errorf("failed to marshal medical research lab data: %v", err)
			}
			err = ctx.GetStub().PutState(ruserId, medLabJSON)
			if err != nil {
				return fmt.Errorf("failed to update medical research lab state: %v", err)
			}
			adminJSON, err := json.Marshal(admin)
			if err != nil {
				return fmt.Errorf("failed to marshal admin data: %v", err)
			}
			err = ctx.GetStub().PutState(senderId, adminJSON)
			if err != nil {
				return fmt.Errorf("failed to update admin state: %v", err)
			}
		case string(TPharmacyCompany):
			pharma, err := s.GetPharmacyCompany(ctx, ruserId)
			if err != nil {
				return fmt.Errorf("pharmacy company with address %s not found", ruserId)
			}

			// Remove senderAddress from adminToPharmacy
			pharma.AdminToPharmacy = removeElement(pharma.AdminToPharmacy, senderId)
			admin.SharedAllUsersAddress = removeElement(admin.SharedAllUsersAddress, senderId)
			// Delete pharmacy company state from ledger
			pharmaJSON, err := json.Marshal(pharma)
			if err != nil {
				return fmt.Errorf("failed to marshal pharmacy company data: %v", err)
			}
			err = ctx.GetStub().PutState(ruserId, pharmaJSON)
			if err != nil {
				return fmt.Errorf("failed to update pharmacy company state: %v", err)
			}
		}

	default:
		return fmt.Errorf("sender address %s does not have permission to revoke access", senderId)
	}

	return nil
}

//	func isAddressInAllAdmins(AllAdmins []int, userID int) bool {
//		userExists := contains(AllAdmins, userID)
//		if !userExists {
//			fmt.Println("User not found in all admins")
//		}
//		return userExists
//	}
func removeElement(slice []string, item string) []string {
	for i, elem := range slice {
		if elem == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// // problem with this function

func (s *SmartContract) AddProfilePic(ctx contractapi.TransactionContextInterface, userId string, url string) error {
	// Determine the entity type and update the profile picture
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	switch accounts[userId] {
	case string(TDoctor):
		doctor, err := s.GetDoctor(ctx, userId)
		if err != nil || !doctor.IsAdded {
			return fmt.Errorf("Doctor with ID=%s not found or not added: %v", userId, err)
		}
		doctor.ProfilePic = url

		// Serialize the updated doctor object
		userBytes, err := json.Marshal(doctor)
		if err != nil {
			return fmt.Errorf("Failed to serialize doctor data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated doctor data for ID=%s: %v", userId, err)
		}

		return nil

	case string(TPatient):
		patient, err := s.GetPatient(ctx, userId)
		if err != nil || !patient.IsAdded {
			return fmt.Errorf("Patient with ID=%s not found or not added: %v", userId, err)
		}
		patient.ProfilePic = url

		// Serialize the updated patient object
		userBytes, err := json.Marshal(patient)
		if err != nil {
			return fmt.Errorf("Failed to serialize patient data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated patient data for ID=%s: %v", userId, err)
		}

		return nil

	case string(TPathologist):
		pathologist, err := s.GetPathologist(ctx, userId)
		if err != nil || !pathologist.IsAdded {
			return fmt.Errorf("Pathologist with ID=%s not found or not added.", userId)
		}
		pathologist.ProfilePic = url

		// Serialize the updated pathologist object
		userBytes, err := json.Marshal(pathologist)
		if err != nil {
			return fmt.Errorf("Failed to serialize pathologist data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated pathologist data for ID=%s: %v", userId, err)
		}

		return nil

	case string(TMedicalResearchLab):
		medicalResearchLab, err := s.GetMedicalResearchLab(ctx, userId)
		if err != nil || !medicalResearchLab.IsAdded {
			return fmt.Errorf("MedicalResearchLab with ID=%s not found or not added.", userId)
		}
		medicalResearchLab.ProfilePic = url

		// Serialize the updated medical research lab object
		userBytes, err := json.Marshal(medicalResearchLab)
		if err != nil {
			return fmt.Errorf("Failed to serialize medical research lab data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated medical research lab data for ID=%s: %v", userId, err)
		}

		return nil

	case string(TPharmacyCompany):
		pharmacyCompany, err := s.GetPharmacyCompany(ctx, userId)
		if err != nil || !pharmacyCompany.IsAdded {
			return fmt.Errorf("PharmacyCompany with ID=%s not found or not added.", userId)
		}
		pharmacyCompany.ProfilePic = url

		// Serialize the updated pharmacy company object
		userBytes, err := json.Marshal(pharmacyCompany)
		if err != nil {
			return fmt.Errorf("Failed to serialize pharmacy company data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated pharmacy company data for ID=%s: %v", userId, err)
		}

		return nil

	case string(TAdmin):
		admin, err := s.GetAdmin(ctx, userId)
		if err != nil || !admin.IsAdded {
			return fmt.Errorf("Admin with ID=%s not found or not added.", userId)
		}
		admin.ProfilePic = url

		// Serialize the updated admin object
		userBytes, err := json.Marshal(admin)
		if err != nil {
			return fmt.Errorf("Failed to serialize admin data: %v", err)
		}

		err = ctx.GetStub().PutState(userId, userBytes)
		if err != nil {
			return fmt.Errorf("Failed to save updated admin data for ID=%s: %v", userId, err)
		}

		return nil

	default:
		return fmt.Errorf("Unsupported entity type for adding a profile picture")
	}
}

func (s *SmartContract) AddTopMedicine(ctx contractapi.TransactionContextInterface, userId string, medicine string) error {

	pharmacyCompany, err := s.GetPharmacyCompany(ctx, userId)
	if err != nil || !pharmacyCompany.IsAdded {
		return fmt.Errorf("PharmacyCompany with ID=%s not found or not added.", userId)
	}
	pharmacyCompany.TopMedichine = append(pharmacyCompany.TopMedichine, medicine)
	userBytes, err := json.Marshal(pharmacyCompany)
	if err != nil {
		return fmt.Errorf("Failed to serialize pharmacy company data: %v", err)
	}

	err = ctx.GetStub().PutState(userId, userBytes)
	if err != nil {
		return fmt.Errorf("Failed to save updated pharmacy company data for ID=%s: %v", userId, err)
	}

	return nil

}

func (s *SmartContract) SetPatient(ctx contractapi.TransactionContextInterface, patientID string, name string, gender string, age int, location string, birthday string, emailAddress string) error {
	// Create a new Patient object
	patient := NewPatient(patientID, name, gender, age, location, birthday, emailAddress)

	// Store the patient in the world state using putState
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal patient: %v", err)
	}
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[patientID] = string(TPatient)
	allPatients, err := s.GetAllPatients(ctx)
if err != nil {
	return fmt.Errorf("failed to retrieve all patients: %v", err)
}
allPatients = append(allPatients, patientID)
if err := s.putState(ctx, AllPatientsKey, allPatients); err != nil {
	return fmt.Errorf("failed to store allPatients update: %v", err)
}

	adminData, err := s.GetAllAdmindata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin data: %v", err)
	}
	adminData.TotalNumberOfPatient += 1
	adminDataJSON, err := json.Marshal(adminData)
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %v", err)
	}
	err = ctx.GetStub().PutState("user-analytics", adminDataJSON)
	if err != nil {
		return fmt.Errorf("failed to save admin data: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}
	return ctx.GetStub().PutState(patientID, patientJSON)
}

func (s *SmartContract) SetPatientPersonalData(ctx contractapi.TransactionContextInterface, patientID string, height float32, blood string, previousDiseases string, medicinedrugs string, badHabits string, chronicDiseases string, healthAllergies string, birthDefects string) error {
	// Retrieve the existing patient from the world state
	patientJSON, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %v", err)
	}
	if patientJSON == nil {
		return fmt.Errorf("patient not found: %v", patientID)
	}

	// Unmarshal the patient data
	var patient Patient
	err = json.Unmarshal(patientJSON, &patient)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient data: %v", err)
	}

	// Update the patient's personal data
	patient.PatientPersonalData.Height = height
	patient.PatientPersonalData.Blood = blood
	patient.PatientPersonalData.PreviousDiseases = previousDiseases
	patient.PatientPersonalData.Medicinedrugs = medicinedrugs
	patient.PatientPersonalData.BadHabits = badHabits
	patient.PatientPersonalData.ChronicDiseases = chronicDiseases
	patient.PatientPersonalData.HealthAllergies = healthAllergies
	patient.PatientPersonalData.BirthDefects = birthDefects

	// Marshal the updated patient object
	updatedPatientJSON, err := json.Marshal(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient: %v", err)
	}

	// Store the updated patient in the world state using PutState
	err = ctx.GetStub().PutState(patientID, updatedPatientJSON)
	if err != nil {
		return fmt.Errorf("failed to put updated patient in world state: %v", err)
	}

	return nil
}
func (sc *SmartContract) AddPrescription(ctx contractapi.TransactionContextInterface, disease string, sUserId string, rUserId string, url []string) error {
	accounts, err := sc.GetAccounts(ctx)
	if err != nil {
		return err
	}

	// Fetch sender's user type and check if the sender is a doctor
	senderType := accounts[sUserId]
	receiverType := accounts[rUserId]

	// Check if the sender is a doctor
	if senderType == string(TDoctor) {
		// Fetch doctor data
		doctor, err := sc.GetDoctor(ctx, sUserId)
		if err != nil || !doctor.IsAdded {
			return fmt.Errorf("doctor with ID=%s not found or not added: %v", sUserId, err)
		}

		// Check receiver type and process accordingly
		switch receiverType {
		case string(TPathologist):
			// Receiver is a Pathologist
			pathologist, err := sc.GetPathologist(ctx, rUserId)
			if err != nil || !pathologist.IsAdded {
				return fmt.Errorf("pathologist with ID=%s not found or not added: %v", rUserId, err)
			}

			// Initialize UserData if not already done
			if doctor.UserData == nil {
				doctor.UserData = make(map[string]UserData)
			}

			// Add URLs to the UserData
			for _, imageUrl := range url {
				data := doctor.UserData[rUserId]
				data.ImagesUrl = append(data.ImagesUrl, imageUrl)
				doctor.UserData[rUserId] = data
			}

			// Check and update relationships
			if !contains(doctor.ReceiverPathologist, rUserId) {
				doctor.ReceiverPathologist = append(doctor.ReceiverPathologist, rUserId)
				pathologist.SenderDoctor = append(pathologist.SenderDoctor, sUserId)
			}
			doctorData, err := json.Marshal(doctor)
			if err != nil {
				return fmt.Errorf("error serializing doctor data: %v", err)
			}
			err = ctx.GetStub().PutState(sUserId, doctorData)
			if err != nil {
				return fmt.Errorf("error saving doctor data to ledger: %v", err)
			}

			// Update pathologist data
			pathologistData, err := json.Marshal(pathologist)
			if err != nil {
				return fmt.Errorf("error serializing pathologist data: %v", err)
			}
			err = ctx.GetStub().PutState(rUserId, pathologistData)
			if err != nil {
				return fmt.Errorf("error saving pathologist data to ledger: %v", err)
			}

		case string(TPatient):
			// Receiver is a Patient
			patient, err := sc.GetPatient(ctx, rUserId)
			if err != nil || !patient.IsAdded {
				return fmt.Errorf("patient with ID=%s not found or not added: %v", rUserId, err)
			}

			// Initialize UserData if not already done
			if doctor.UserData == nil {
				doctor.UserData = make(map[string]UserData)
			}

			// Add URLs to the UserData
			for _, imageUrl := range url {
				data := doctor.UserData[rUserId]
				data.ImagesUrl = append(data.ImagesUrl, imageUrl)
				doctor.UserData[rUserId] = data
			}

			// Check and update relationships
			if !contains(doctor.TreatedPatient, rUserId) {
				doctor.TreatedPatient = append(doctor.TreatedPatient, rUserId)
				patient.PersonalDoctor = append(patient.PersonalDoctor, sUserId)
			}
			doctorData, err := json.Marshal(doctor)
			if err != nil {
				return fmt.Errorf("error serializing doctor data: %v", err)
			}
			err = ctx.GetStub().PutState(sUserId, doctorData)
			if err != nil {
				return fmt.Errorf("error saving doctor data to ledger: %v", err)
			}

			// Update patient data
			patientData, err := json.Marshal(patient)
			if err != nil {
				return fmt.Errorf("error serializing patient data: %v", err)
			}
			err = ctx.GetStub().PutState(rUserId, patientData)
			if err != nil {
				return fmt.Errorf("error saving patient data to ledger: %v", err)
			}

		default:
			return fmt.Errorf("unsupported receiver type")
		}

	} else if senderType == string(TPathologist) {
		// Sender is a Pathologist
		pathologist, err := sc.GetPathologist(ctx, sUserId)
		if err != nil || !pathologist.IsAdded {
			return fmt.Errorf("pathologist with ID=%s not found or not added: %v", sUserId, err)
		}

		// Receiver is a Doctor
		if receiverType == string(TDoctor) {
			doctor, err := sc.GetDoctor(ctx, rUserId)
			if err != nil || !doctor.IsAdded {
				return fmt.Errorf("doctor with ID=%s not found or not added: %v", rUserId, err)
			}

			// Initialize UserData if not already done
			if pathologist.UserData == nil {
				pathologist.UserData = make(map[string]UserData)
			}

			// Add URLs to the UserData
			for _, imageUrl := range url {
				data := pathologist.UserData[rUserId]
				data.ImagesUrl = append(data.ImagesUrl, imageUrl)
				pathologist.UserData[rUserId] = data
			}

			// Check and update relationships
			if !contains(pathologist.ReceiverDoctor, rUserId) {
				pathologist.ReceiverDoctor = append(pathologist.ReceiverDoctor, rUserId)
				doctor.SenderPathologist = append(doctor.SenderPathologist, sUserId)
			}

			pathologistData, err := json.Marshal(pathologist)
			if err != nil {
				return fmt.Errorf("error serializing pathologist data: %v", err)
			}
			err = ctx.GetStub().PutState(sUserId, pathologistData)
			if err != nil {
				return fmt.Errorf("error saving pathologist data to ledger: %v", err)
			}

			// Update doctor data
			doctorData, err := json.Marshal(doctor)
			if err != nil {
				return fmt.Errorf("error serializing doctor data: %v", err)
			}
			err = ctx.GetStub().PutState(rUserId, doctorData)
			if err != nil {
				return fmt.Errorf("error saving doctor data to ledger: %v", err)
			}

		} else {
			return fmt.Errorf("unsupported receiver type for pathologist")
		}

	} else if rUserId == "-1" {
		if senderType == string(TPatient) {
			// Sender is a Patient
			patient, err := sc.GetPatient(ctx, sUserId)
			if err != nil || !patient.IsAdded {
				return fmt.Errorf("patient with ID=%s not found or not added: %v", sUserId, err)
			}
			if patient.Prescriptions == nil {
				patient.Prescriptions = make(map[string][]string)
			}
			patient.Prescriptions[disease] = append(patient.Prescriptions[disease], url...)

			// Serialize and put state for patient
			patientData, err := json.Marshal(patient)
			if err != nil {
				return fmt.Errorf("error serializing patient data: %v", err)
			}
			err = ctx.GetStub().PutState(sUserId, patientData)
			if err != nil {
				return fmt.Errorf("error saving patient data to ledger: %v", err)
			}

		} else if senderType == string(TMedicalResearchLab) {
			// Sender is a Medical Research Lab
			medicalResearchLab, err := sc.GetMedicalResearchLab(ctx, sUserId)
			if err != nil || !medicalResearchLab.IsAdded {
				return fmt.Errorf("medical research lab with ID=%s not found or not added: %v", sUserId, err)
			}
			medicalResearchLab.ImgUrl = append(medicalResearchLab.ImgUrl, url...)

			medicalResearchLabData, err := json.Marshal(medicalResearchLab)
			if err != nil {
				return fmt.Errorf("error serializing medical research lab data: %v", err)
			}
			err = ctx.GetStub().PutState(sUserId, medicalResearchLabData)
			if err != nil {
				return fmt.Errorf("error saving medical research lab data to ledger: %v", err)
			}

		}

	}
	return nil

}

// Helper function to check if a value exists in a slice
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func (s *SmartContract) ShareData(
	ctx contractapi.TransactionContextInterface,
	sUserId string, rUserId string,
) error {
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	// Retrieve sender and receiver types
	senderType, senderExists := accounts[sUserId]
	if !senderExists {
		return fmt.Errorf("sender with ID=%s not found", sUserId)
	}
	receiverType, receiverExists := accounts[rUserId]
	if !receiverExists {
		return fmt.Errorf("receiver with ID=%s not found", rUserId)
	}

	// Process sharing logic based on sender type
	switch senderType {
	case string(TPatient):
		// Fetch patient data
		patient, err := s.GetPatient(ctx, sUserId)
		if err != nil || !patient.IsAdded {
			return fmt.Errorf("patient with ID=%s not found or not added: %v", sUserId, err)
		}

		switch receiverType {
		case string(TDoctor):
			// Fetch doctor data
			doctor, err := s.GetDoctor(ctx, rUserId)
			if err != nil || !doctor.IsAdded {
				return fmt.Errorf("doctor with ID=%s not found or not added: %v", rUserId, err)
			}

			// Check if data is already shared
			if contains(doctor.PatientToDoctor, sUserId) {
				return fmt.Errorf("data already shared with doctor ID=%s", rUserId)
			} else {
				doctor.PatientToDoctor = append(doctor.PatientToDoctor, sUserId)
				patient.SharedAllDoctorAddress = append(patient.SharedAllDoctorAddress, rUserId)
			}

			// Save updated entities back to the ledger
			if err := s.putState(ctx, sUserId, patient); err != nil {
				return fmt.Errorf("failed to update patient data: %v", err)
			}
			if err := s.putState(ctx, rUserId, doctor); err != nil {
				return fmt.Errorf("failed to update doctor data: %v", err)
			}

		case string(TAdmin):
			// Fetch admin data
			admin, err := s.GetAdmin(ctx, rUserId)
			if err != nil || !admin.IsAdded {
				return fmt.Errorf("admin with ID=%s not found or not added: %v", rUserId, err)
			}

			// Check if data is already shared
			if contains(admin.PatientToAdmin, sUserId) {
				return fmt.Errorf("data already shared with admin ID=%s", rUserId)
			} else {
				admin.PatientToAdmin = append(admin.PatientToAdmin, sUserId)
				if admin.AllPrescription == nil {
					admin.AllPrescription = make(map[string][]string)
				}

				for disease, prescriptions := range patient.Prescriptions {
					if _, exists := admin.AllPrescription[disease]; !exists {
						admin.AllPrescription[disease] = []string{}
					}
					admin.AllPrescription[disease] = append(admin.AllPrescription[disease], prescriptions...)
				}

			}

			// Save updated admin back to the ledger
			if err := s.putState(ctx, rUserId, admin); err != nil {
				return fmt.Errorf("failed to update admin data: %v", err)
			}

		default:
			return fmt.Errorf("invalid receiver type for patient ID=%s", sUserId)
		}

	case string(TAdmin):
		if sUserId == ownerId && receiverType == string(TMedicalResearchLab) || sUserId == ownerId && receiverType == string(TPharmacyCompany) {

			pendingTx, err := s.GetPendingTx(ctx)
			if err != nil {
				return fmt.Errorf("failed to retrieve pending transactions: %v", err)
			}
			if contains(pendingTx[ownerId], rUserId) {
				return fmt.Errorf("transaction already in pending for receiver ID=%s", rUserId)
			}
			if _, exists := pendingTx[ownerId]; !exists {
				pendingTx[ownerId] = make([]string, 0) // Initialize empty slice
			}
			pendingTx[ownerId] = append(pendingTx[ownerId], rUserId)
			medicalResearchLab, err := s.GetMedicalResearchLab(ctx, rUserId)
			if err != nil || !medicalResearchLab.IsAdded {
				return fmt.Errorf("medical research lab with ID=%s not found or not added: %v", rUserId, err)
			}
			pharmacyCompany, err := s.GetPharmacyCompany(ctx, rUserId)
			if err != nil || !pharmacyCompany.IsAdded {
				return fmt.Errorf("pharmacy company with ID=%s not found or not added: %v", rUserId, err)
			}
			if contains(medicalResearchLab.AdminToMedRcLab, ownerId) || contains(pharmacyCompany.AdminToPharmacy, ownerId) {
				return fmt.Errorf("admin ID=%s already has access to receiver ID=%s", sUserId, rUserId)

			}
			transaction := Transaction{
				From:          ownerId,
				To:            rUserId,
				Executed:      false,
				Confirmations: 0,
			}
			transactions, err := s.GetTransactions(ctx)
			if err != nil {
				return fmt.Errorf("failed to retrieve transactions: %v", err)
			}

			// Update the transactions map with the new transaction
			transactions[rUserId] = transaction

			if err := s.putState(ctx, TransactionsKey, transactions); err != nil {
				return fmt.Errorf("failed to store transaction data: %v", err)
			}
			// pendingTxKey := "PendingTx"
			if err := s.putState(ctx, PendingTxKey, pendingTx); err != nil {
				return fmt.Errorf("failed to store pending transactions: %v", err)
			}
		} else {
			return fmt.Errorf("invalid receiver type for  owner ID=%s", sUserId)
		}

	default:
		return fmt.Errorf("invalid sender type ID=%s", sUserId)
	}

	return nil
}

func (s *SmartContract) putState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize data for key=%s: %v", key, err)
	}
	if err := ctx.GetStub().PutState(key, valueBytes); err != nil {
		return fmt.Errorf("failed to store data for key=%s: %v", key, err)
	}
	return nil
}

func (s *SmartContract) GetPathologistDataFromDoctor(ctx contractapi.TransactionContextInterface, doctorId string, pathologistId string) ([]string, error) {
	doctor, err := s.GetDoctor(ctx, doctorId)
	if err != nil || !doctor.IsAdded {
		return nil, fmt.Errorf("Doctor with ID=%d not found or not added.", doctorId)
	}
	pathologist, err := s.GetPathologist(ctx, pathologistId)
	if err != nil || !pathologist.IsAdded {
		return nil, fmt.Errorf("Pathologist with ID=%d not found or not added.", pathologistId)
	}
	userData, exists := doctor.UserData[pathologist.PathologistID]
	if !exists {
		return nil, fmt.Errorf("No data found for Pathologist ID=%d in Doctor ID=%d's records.", pathologistId, doctorId)
	}
	return userData.ImagesUrl, nil
}

func (s *SmartContract) GetPatientDataFromDoctor(ctx contractapi.TransactionContextInterface, doctorId string, patientID string) ([]string, error) {
	doctor, err := s.GetDoctor(ctx, doctorId)
	if err != nil || !doctor.IsAdded {
		return nil, fmt.Errorf("Doctor with ID=%d not found or not added.", doctorId)
	}
	patient, err := s.GetPatient(ctx, patientID)
	if err != nil || !patient.IsAdded {
		return nil, fmt.Errorf("Patient with ID=%d not found or not added.", patientID)
	}

	userData, exists := doctor.UserData[patient.PatientID]
	if !exists {
		return nil, fmt.Errorf("No data found for Pathologist ID=%d in Doctor ID=%d's records.", patientID, doctorId)
	}
	return userData.ImagesUrl, nil
}

func (s *SmartContract) GetDoctorDataFromPathologist(ctx contractapi.TransactionContextInterface, pathologistId string, doctorId string) ([]string, error) {

	pathologist, err := s.GetPathologist(ctx, pathologistId)
	if err != nil || !pathologist.IsAdded {
		return nil, fmt.Errorf("Pathologist with ID=%d not found or not added.", pathologistId)
	}
	doctor, err := s.GetDoctor(ctx, doctorId)
	if err != nil || !doctor.IsAdded {
		return nil, fmt.Errorf("Doctor with ID=%d not found or not added.", doctorId)
	}
	userData, exists := pathologist.UserData[doctor.DoctorID]
	if !exists {
		return nil, fmt.Errorf("No data found for Doctor ID=%d in Pathologist ID=%d's records.", doctorId, pathologistId)
	}
	return userData.ImagesUrl, nil
}

func (s *SmartContract) GetPatient(ctx contractapi.TransactionContextInterface, id string) (Patient, error) {
	// Retrieve the patient data from the world state using the provided patient ID
	patientJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Patient{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if patientJSON == nil {
		return Patient{}, fmt.Errorf("patient with ID %s does not exist", id)
	}

	// Unmarshal the JSON data into a Patient object
	var patient Patient
	err = json.Unmarshal(patientJSON, &patient)
	if err != nil {
		return Patient{}, fmt.Errorf("failed to unmarshal patient data: %v", err)
	}

	// Return the patient data
	return patient, nil
}

func (s *SmartContract) SetDoctor(ctx contractapi.TransactionContextInterface, doctorID string, name string, specialty string, consultationFee float32, BMDCNumber int, yearOfExperience int, birthday string, emailAddress string) error {
	// Create a new Doctor object
	doctor := NewDoctor(doctorID, name, specialty, consultationFee, BMDCNumber, yearOfExperience, birthday, emailAddress)

	// Marshal the doctor object to JSON
	doctorJSON, err := json.Marshal(doctor)
	if err != nil {
		return fmt.Errorf("failed to marshal doctor: %v", err)
	}
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[doctorID] = string(TDoctor)
	adminData, err := s.GetAllAdmindata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin data: %v", err)
	}
	adminData.TotalNumberOfDoctor += 1
	adminDataJSON, err := json.Marshal(adminData)
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %v", err)
	}
	err = ctx.GetStub().PutState("user-analytics", adminDataJSON)
	if err != nil {
		return fmt.Errorf("failed to save admin data: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}
	return ctx.GetStub().PutState(doctorID, doctorJSON)
}

func (s *SmartContract) GetDoctor(ctx contractapi.TransactionContextInterface, id string) (Doctor, error) {
	// Retrieve the doctor data from the world state using the provided doctor ID
	doctorJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Doctor{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if doctorJSON == nil {
		return Doctor{}, fmt.Errorf("doctor with ID %s does not exist", id)
	}

	// Unmarshal the JSON data into a Doctor object
	var doctor Doctor
	err = json.Unmarshal(doctorJSON, &doctor)
	if err != nil {
		return Doctor{}, fmt.Errorf("failed to unmarshal doctor data: %v", err)
	}

	return doctor, nil
}

type Pathologist struct {
	PathologistID        string              `json:"pathologistID"`
	Name                 string              `json:"name"`
	LicenseNumber        int                 `json:"licenseNumber"`
	SpecializationArea   string              `json:"specializationArea"`
	TotalExperience      int                 `json:"totalExperience"`
	IsAdded              bool                `json:"isAdded"`
	PatientToPathologist []string            `json:"patientToPathologist"` // allPatientsAddressSharedToPathologist
	UserType             string              `json:"userType"`
	ProfilePic           string              `json:"profilePic"`
	Birthday             string              `json:"birthday"`
	UserData             map[string]UserData `json:"userData"`
	SenderDoctor         []string            `json:"senderDoctor"` // allDoctorsAddressSharedToPathologist
	ReceiverDoctor       []string            `json:"receiverDoctor"`
	EmailAddress         string              `json:"emailAddress"`
}

func NewPathologist(pathologistID string, name string, licenseNumber int, specializationArea string, totalExperience int, birthday string, emailAddress string) Pathologist {
	return Pathologist{
		PathologistID:        pathologistID,
		Name:                 name,
		LicenseNumber:        licenseNumber,
		SpecializationArea:   specializationArea,
		TotalExperience:      totalExperience,
		IsAdded:              true,          // Default to true
		PatientToPathologist: []string{},    // Initialize as an empty slice
		UserType:             "pathologist", // Default user           // Default profile picture
		Birthday:             birthday,
		UserData:             map[string]UserData{},
		SenderDoctor:         []string{}, // Initialize as an empty slice
		ReceiverDoctor:       []string{}, // Initialize as an empty slice
		EmailAddress:         emailAddress,
	}
}

func (s *SmartContract) SetPathologist(ctx contractapi.TransactionContextInterface, pathologistID string, name string, licenseNumber int, specializationArea string, totalExperience int, birthday string, emailAddress string) error {
	// Create a new Pathologist object
	pathologist := NewPathologist(pathologistID, name, licenseNumber, specializationArea, totalExperience, birthday, emailAddress)

	// Marshal the pathologist object to JSON
	pathologistJSON, err := json.Marshal(pathologist)
	if err != nil {
		return fmt.Errorf("failed to marshal pathologist: %v", err)
	}
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[pathologistID] = string(TPathologist)
	adminData, err := s.GetAllAdmindata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin data: %v", err)
	}
	adminData.TotalNumberOfPathologist += 1
	adminDataJSON, err := json.Marshal(adminData)
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %v", err)
	}
	err = ctx.GetStub().PutState("user-analytics", adminDataJSON)
	if err != nil {
		return fmt.Errorf("failed to save admin data: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}

	return ctx.GetStub().PutState(pathologistID, pathologistJSON)
}

func (s *SmartContract) GetPathologist(ctx contractapi.TransactionContextInterface, id string) (Pathologist, error) {
	// Retrieve the pathologist data from the world state using the provided pathologist ID
	pathologistJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Pathologist{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if pathologistJSON == nil {
		return Pathologist{}, fmt.Errorf("pathologist with ID %s does not exist", id)
	}

	// Unmarshal the JSON data into a Pathologist object
	var pathologist Pathologist
	err = json.Unmarshal(pathologistJSON, &pathologist)
	if err != nil {
		return Pathologist{}, fmt.Errorf("failed to unmarshal pathologist data: %v", err)
	}

	// Return the pathologist data
	return pathologist, nil
}

type MedicalResearchLab struct {
	LabID           string   `json:"labID"`
	Name            string   `json:"name"`
	LicenseID       int      `json:"licenseID"`
	ResearchArea    string   `json:"researchArea"`
	LabRating       float32  `json:"labRating"`
	IsAdded         bool     `json:"isAdded"`
	ImgUrl          []string `json:"imgUrl"` // MedicalResearchLabReports
	UserType        string   `json:"userType"`
	ProfilePic      string   `json:"profilePic"`
	Prescriptions          map[string][]string `json:"prescriptions"`
	EmailAddress    string   `json:"emailAddress"`
	AdminToMedRcLab []string `json:"adminToMedRcLab"` // vvv

}

func NewMedicalResearchLab(labID string, name string, licenseID int, researchArea string, labRating float32, emailAddress string) MedicalResearchLab {
	return MedicalResearchLab{
		LabID:           labID,
		Name:            name,
		LicenseID:       licenseID,
		ResearchArea:    researchArea,
		LabRating:       labRating,
		IsAdded:         true,
		UserType:        "medicalResearchLab",
		EmailAddress:    emailAddress,
		ImgUrl:          []string{}, // Default to an empty list        // Default to an empty string
		AdminToMedRcLab: []string{}, // Default to an empty list
	}
}

func (s *SmartContract) SetMedicalResearchLab(ctx contractapi.TransactionContextInterface, labID string, name string, licenseID int, researchArea string, labRating float32, emailAddress string) error {
	// Create a new MedicalResearchLab object
	medicalResearchLab := NewMedicalResearchLab(labID, name, licenseID, researchArea, labRating, emailAddress)

	// Marshal the medicalResearchLab object to JSON
	medicalResearchLabJSON, err := json.Marshal(medicalResearchLab)
	if err != nil {
		return fmt.Errorf("failed to marshal medical research lab: %v", err)
	}
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[labID] = string(TMedicalResearchLab)
	adminData, err := s.GetAllAdmindata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin data: %v", err)
	}
	adminData.TotalNumberOfMedicalResearchLab += 1
	adminDataJSON, err := json.Marshal(adminData)
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %v", err)
	}
	err = ctx.GetStub().PutState("user-analytics", adminDataJSON)
	if err != nil {
		return fmt.Errorf("failed to save admin data: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}
	return ctx.GetStub().PutState(labID, medicalResearchLabJSON)
}

func (s *SmartContract) GetMedicalResearchLab(ctx contractapi.TransactionContextInterface, id string) (MedicalResearchLab, error) {
	labJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return MedicalResearchLab{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if labJSON == nil {
		return MedicalResearchLab{}, fmt.Errorf("medical research lab with ID %d does not exist", id)
	}

	// Unmarshal the JSON data into a MedicalResearchLab object
	var lab MedicalResearchLab
	err = json.Unmarshal(labJSON, &lab)
	if err != nil {
		return MedicalResearchLab{}, fmt.Errorf("failed to unmarshal medical research lab data: %v", err)
	}

	// Return the medical research lab data
	return lab, nil
}

type PharmacyCompany struct {
	CompanyID          string   `json:"companyID"`
	Name               string   `json:"name"`
	LicenseID          int      `json:"licenseID"`
	ProductInformation string   `json:"productInformation"`
	PharmacyRating     float32  `json:"pharmacyRating"`
	IsAdded            bool     `json:"isAdded"`
	UserType           string   `json:"userType"`
	TopMedichine       []string `json:"topMedichine"`
	ProfilePic         string   `json:"profilePic"`
	EmailAddress       string   `json:"emailAddress"`
	Prescriptions          map[string][]string `json:"prescriptions"`
	AdminToPharmacy    []string `json:"adminToPharmacy"` // zzz

}

func NewPharmacyCompany(companyID string, name string, licenseID int, productInformation string, pharmacyRating float32, emailAddress string) PharmacyCompany {
	return PharmacyCompany{
		CompanyID:          companyID,
		Name:               name,
		LicenseID:          licenseID,
		ProductInformation: productInformation,
		PharmacyRating:     pharmacyRating,
		IsAdded:            true,
		UserType:           "pharmacyCompany",
		TopMedichine:       []string{}, // Default: empty list // Default: placeholder image
		EmailAddress:       emailAddress,
		AdminToPharmacy:    []string{}, // Default: empty list
	}
}
func (s *SmartContract) SetPharmacyCompany(ctx contractapi.TransactionContextInterface, companyID string, name string, licenseID int, productInformation string, pharmacyRating float32, emailAddress string) error {
	// Create a new PharmacyCompany object
	pharmacyCompany := NewPharmacyCompany(companyID, name, licenseID, productInformation, pharmacyRating, emailAddress)

	// Marshal the pharmacyCompany object to JSON
	pharmacyCompanyJSON, err := json.Marshal(pharmacyCompany)
	if err != nil {
		return fmt.Errorf("failed to marshal pharmacy company: %v", err)
	}
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	accounts[companyID] = string(TPharmacyCompany)
	adminData, err := s.GetAllAdmindata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get admin data: %v", err)
	}
	adminData.TotalNumberOfPharmacyCompany += 1
	adminDataJSON, err := json.Marshal(adminData)
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %v", err)
	}
	err = ctx.GetStub().PutState("user-analytics", adminDataJSON)
	if err != nil {
		return fmt.Errorf("failed to save admin data: %v", err)
	}
	if err := s.putState(ctx, "accounts", accounts); err != nil {
		return fmt.Errorf("failed to store  accounts update: %v", err)
	}
	return ctx.GetStub().PutState(companyID, pharmacyCompanyJSON)
}

func (s *SmartContract) GetPharmacyCompany(ctx contractapi.TransactionContextInterface, id string) (PharmacyCompany, error) {
	// Retrieve the pharmacy company data from the world state using the provided company ID
	companyJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return PharmacyCompany{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// If no data is found, return an error
	if companyJSON == nil {
		return PharmacyCompany{}, fmt.Errorf("pharmacy company with ID %s does not exist", id)
	}

	// Unmarshal the JSON data into a PharmacyCompany object
	var company PharmacyCompany
	err = json.Unmarshal(companyJSON, &company)
	if err != nil {
		return PharmacyCompany{}, fmt.Errorf("failed to unmarshal pharmacy company data: %v", err)
	}

	// Return the pharmacy company data
	return company, nil
}

func (s *SmartContract) ConnectedAccountType(ctx contractapi.TransactionContextInterface, userId string) (string, error) {
	// Fetch the user data from the ledger
	userBytes, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user with ID=%s: %v", userId, err)
	}
	if userBytes == nil {

		return string(TUnknown), nil
	}

	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return "", err
	}

	accountType, exists := accounts[userId]
	if !exists {

		return string(TUnknown), nil
	}

	switch EntityType(accountType) {
	case TDoctor, TPathologist, TAdmin, TPatient, TMedicalResearchLab, TPharmacyCompany:
		return accountType, nil
	default:
		return string(TUnknown), nil
	}
}

// func (s *SmartContract) ConnectedAccountType(ctx contractapi.TransactionContextInterface, userId string) (string, error) {
// 	// Fetch the user data from the ledger
// 	userBytes, err := ctx.GetStub().GetState(userId)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to retrieve user with ID=%s: %v", userId, err)
// 	}
// 	if userBytes == nil {
// 		return "", fmt.Errorf("user with ID=%s not found", userId)
// 	}
// 	accounts, err := s.GetAccounts(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	accountType, exists := accounts[userId]
// 	if !exists {
// 		return string(TUnknown), fmt.Errorf("account type for user ID=%s not found", userId)
// 	}

// 	// Validate the account type
// 	switch EntityType(accountType) {
// 	case TDoctor, TPathologist, TAdmin, TPatient, TMedicalResearchLab, TPharmacyCompany:
// 		return accountType, nil
// 	default:
// 		return string(TUnknown), fmt.Errorf("invalid or unknown user type")
// 	}
// }

type UserData struct {
	ImagesUrl []string `json:"imagesUrl"`
}

//	type ShareData struct {
//		SharedAllUsersAddress []string `json:"sharedAllUsersAddress"`
//	}
const ownerId string = "107657684236227940746"

type Transaction struct {
	From          string `json:"from"`
	To            string `json:"to"`
	Executed      bool   `json:"executed"`
	Confirmations int    `json:"confirmations"`
}

const (
	TransactionsKey = "transactions"
	PendingTxKey    = "pendingTx"
	AllAdminsKey    = "allAdmins"
	AllPatientsKey    = "allPatients"
	IsConfirmedKey  = "isConfirmed"
	PendingRequestedUserKey = "pendingRequestedUser"
	PendingRequesterUserKey ="pendingRequesterUser"
)

func (s *SmartContract) GetPendingTx(ctx contractapi.TransactionContextInterface) (map[string][]string, error) {
	// Retrieve the "pendingTx" state from the ledger
	pendingTxJSON, err := ctx.GetStub().GetState(PendingTxKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pending transactions from the ledger: %v", err)
	}

	// Check if "pendingTx" exists in the ledger
	if pendingTxJSON == nil {
		return make(map[string][]string), nil
	}

	// Deserialize the JSON into a map[string][]string
	var pendingTx map[string][]string
	err = json.Unmarshal(pendingTxJSON, &pendingTx)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize pending transactions: %v", err)
	}

	return pendingTx, nil
}

func (s *SmartContract) GetPendingRequestedUser(ctx contractapi.TransactionContextInterface) (map[string]map[string][]string, error) {
	// Retrieve the pending transactions from the ledger
	pendingRequestedUserJSON, err := ctx.GetStub().GetState(PendingRequestedUserKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pending requested users from the ledger: %v", err)
	}

	// If no data exists, return an empty map
	if pendingRequestedUserJSON == nil {
		return make(map[string]map[string][]string), nil
	}

	// Deserialize the JSON into a nested map
	var pendingRequestedUser map[string]map[string][]string
	err = json.Unmarshal(pendingRequestedUserJSON, &pendingRequestedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize pending requested users: %v", err)
	}

	return pendingRequestedUser, nil
}

func (s *SmartContract) GetPendingRequesterUser(ctx contractapi.TransactionContextInterface) (map[string]map[string][]string, error) {
	// Retrieve the pending transactions from the ledger
	pendingRequesterUserJSON, err := ctx.GetStub().GetState(PendingRequesterUserKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pending requester users from the ledger: %v", err)
	}

	// If no data exists, return an empty map
	if pendingRequesterUserJSON == nil {
		return make(map[string]map[string][]string), nil
	}

	// Deserialize the JSON into a nested map
	var pendingRequesterUser map[string]map[string][]string
	err = json.Unmarshal(pendingRequesterUserJSON, &pendingRequesterUser)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize pending requester users: %v", err)
	}

	return pendingRequesterUser, nil
}


func (s *SmartContract) GetTransactions(ctx contractapi.TransactionContextInterface) (map[string]Transaction, error) {
	// Retrieve the "transactions" state from the ledger
	transactionsJSON, err := ctx.GetStub().GetState(TransactionsKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions from the ledger: %v", err)
	}

	// Check if "transactions" exists in the ledger
	if transactionsJSON == nil {
		return make(map[string]Transaction), nil
	}

	// Deserialize the JSON into a map[string]Transaction
	var transactions map[string]Transaction
	err = json.Unmarshal(transactionsJSON, &transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transactions: %v", err)
	}

	return transactions, nil
}

func (s *SmartContract) GetAllAdmins(ctx contractapi.TransactionContextInterface) ([]string, error) {
	// Retrieve the "allAdmins" state from the ledger
	allAdminsJSON, err := ctx.GetStub().GetState(AllAdminsKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve allAdmins from the ledger: %v", err)
	}

	// Check if "allAdmins" exists in the ledger
	if allAdminsJSON == nil {
		return []string{}, nil
	}

	// Deserialize the JSON into a slice of strings
	var allAdmins []string
	err = json.Unmarshal(allAdminsJSON, &allAdmins)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize allAdmins: %v", err)
	}

	return allAdmins, nil
}

func (s *SmartContract) GetAllPatients(ctx contractapi.TransactionContextInterface) ([]string, error) {
	// Retrieve the "allPatients" state from the ledger
	allPatientsJSON, err := ctx.GetStub().GetState(AllPatientsKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve allPatients from the ledger: %v", err)
	}

	// Check if "allPatients" exists in the ledger
	if allPatientsJSON == nil {
		return []string{}, nil
	}

	// Deserialize the JSON into a slice of strings
	var allPatients []string
	err = json.Unmarshal(allPatientsJSON, &allPatients)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize allPatients: %v", err)
	}

	return allPatients, nil
}

func (s *SmartContract) GetIsConfirmed(ctx contractapi.TransactionContextInterface) (map[string]map[string]bool, error) {
	// Retrieve the "isConfirmed" state from the ledger
	isConfirmedJSON, err := ctx.GetStub().GetState(IsConfirmedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve isConfirmed from the ledger: %v", err)
	}

	// Check if "isConfirmed" exists in the ledger
	if isConfirmedJSON == nil {
		return make(map[string]map[string]bool), nil
	}

	// Deserialize the JSON into a map[string]map[string]bool
	var isConfirmed map[string]map[string]bool
	err = json.Unmarshal(isConfirmedJSON, &isConfirmed)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize isConfirmed: %v", err)
	}

	return isConfirmed, nil
}

// func onlyAdmin(userId int, adminMap map[int]Admin) error {
// 	admin, exists := adminMap[userId]
// 	if !exists || !admin.IsAdded {
// 		return errors.New("you must be an admin")
// 	}
// 	return nil
// }

type EntityType string

const (
	TUnknown            EntityType = "noExsistEntity"
	TDoctor             EntityType = "Doctor"
	TPathologist        EntityType = "Pathologist"
	TMedicalResearchLab EntityType = "MedicalResearchLab"
	TPharmacyCompany    EntityType = "PharmacyCompany"
	
	TPatient            EntityType = "Patient"
	TAdmin              EntityType = "Admin"
)

func (s *SmartContract) GiveConfirmation(ctx contractapi.TransactionContextInterface, userId string, adminUserId string) error {
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return err
	}

	if accounts[adminUserId] != string(TAdmin) {
		return fmt.Errorf("user ID=%d is not an admin", adminUserId)
	}
	pendingTx, err := s.GetPendingTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending transactions: %v", err)
	}
	if !contains(pendingTx[ownerId], userId) {
		return fmt.Errorf("user ID=%d does not have a pending transaction", userId)
	}
	isConfirmed, err := s.GetIsConfirmed(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve isConfirmed: %v", err)
	}
	if _, exists := isConfirmed[adminUserId]; !exists {
		isConfirmed[adminUserId] = make(map[string]bool)
	}
	if isConfirmed[adminUserId][userId] {
		return fmt.Errorf("user ID=%d has already confirmed transaction for admin ID=%d", userId, adminUserId)
	}

	admin, err := s.GetAdmin(ctx, adminUserId)
	if err != nil || !admin.IsAdded {
		return fmt.Errorf("Admin with ID=%s not found or not added.", userId)
	}
	owner, err := s.GetAdmin(ctx, ownerId)
	if err != nil || !owner.IsAdded {
		return fmt.Errorf("Admin with ID=%s not found or not added.", userId)
	}
	transactions, err := s.GetTransactions(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions: %v", err)
	}
	tx, txExists := transactions[userId]
	if !txExists {
		return fmt.Errorf("transaction for user ID=%d does not exist", userId)
	}

	if tx.Executed {
		return fmt.Errorf("transaction for user ID=%d is already confirmed", userId)
	}
	allAdmins, err := s.GetAllAdmins(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve all admins: %v", err)
	}
	totalAdmins := len(allAdmins)
	if tx.Confirmations == totalAdmins-1 {
		userAccountType := accounts[tx.To]

		if userAccountType == string(TMedicalResearchLab) {
			// Add admin to MedicalResearchLab
			medicalResearchLab, err := s.GetMedicalResearchLab(ctx, tx.To)
			if err != nil || !medicalResearchLab.IsAdded {
				return fmt.Errorf("MedicalResearchLab with ID=%s not found or not added.", userId)
			}

			medicalResearchLab.AdminToMedRcLab = append(medicalResearchLab.AdminToMedRcLab, ownerId)
			owner.SharedAllUsersAddress = append(owner.SharedAllUsersAddress, medicalResearchLab.LabID)
			pendingTx, err := s.GetPendingTx(ctx)
			if err != nil {
				return fmt.Errorf("failed to retrieve pending transactions: %v", err)
			}
			pendingTx[ownerId] = removeElement(pendingTx[ownerId], userId)
			if err := s.putState(ctx, PendingTxKey, pendingTx); err != nil {
				return fmt.Errorf("failed to store  pendingTx: %v", err)
			}
			medicalResearchLabData, err := json.Marshal(medicalResearchLab)
			if err != nil {
				return fmt.Errorf("error serializing medical research lab data: %v", err)
			}
			err = ctx.GetStub().PutState(tx.To, medicalResearchLabData)
			if err != nil {
				return fmt.Errorf("error saving medical research lab data to ledger: %v", err)
			}

		} else if userAccountType == string(TPharmacyCompany) {
			// Add admin to PharmacyCompany
			pharmacy, err := s.GetPharmacyCompany(ctx, tx.To)
			if err != nil || !pharmacy.IsAdded {
				return fmt.Errorf("PharmacyCompany with ID=%s not found or not added.", userId)
			}

			pharmacy.AdminToPharmacy = append(pharmacy.AdminToPharmacy, ownerId)
			owner.SharedAllUsersAddress = append(owner.SharedAllUsersAddress, pharmacy.CompanyID)

			pendingTx[ownerId] = removeElement(pendingTx[ownerId], userId)
			if err := s.putState(ctx, PendingTxKey, pendingTx); err != nil {
				return fmt.Errorf("failed to store  pendingTx: %v", err)
			}
			pharmacyData, err := json.Marshal(pharmacy)
			if err != nil {
				return fmt.Errorf("error serializing pharmacy company data: %v", err)
			}
			err = ctx.GetStub().PutState(tx.To, pharmacyData)
			if err != nil {
				return fmt.Errorf("error saving pharmacy company data to ledger: %v", err)
			}

		} else {
			return fmt.Errorf("only transactions for MedicalResearchLab or PharmacyCompany can be executed")
		}

		// Mark transaction as executed
		tx.Executed = true
		transactions[userId] = tx

		if err := s.putState(ctx, TransactionsKey, transactions); err != nil {
			return fmt.Errorf("failed to store  transactions: %v", err)
		}

		if err := s.putState(ctx, PendingTxKey, pendingTx); err != nil {
			return fmt.Errorf("failed to store pending transactions: %v", err)
		}
		adminJson, err := json.Marshal(admin)
		if err != nil {
			return fmt.Errorf("failed to marshal admin: %v", err)
		}
		err = ctx.GetStub().PutState(adminUserId, adminJson)
		if err != nil {
			return fmt.Errorf("failed to update admin in world state: %v", err)
		}
		ownerJson, err := json.Marshal(owner)
		if err != nil {
			return fmt.Errorf("failed to marshal owner: %v", err)
		}
		err = ctx.GetStub().PutState(ownerId, ownerJson)
		if err != nil {
			return fmt.Errorf("failed to update owner in world state: %v", err)
		}
		// Reset confirmation state
		// isConfirmed[adminUserId][userId] = false
	} else {
		// Increment the confirmations
		tx.Confirmations++
		isConfirmed[adminUserId][userId] = true
		transactions[userId] = tx
		if err := s.putState(ctx, IsConfirmedKey, isConfirmed); err != nil {
			return fmt.Errorf("failed to store  pendingTx: %v", err)
		}
		if err := s.putState(ctx, TransactionsKey, transactions); err != nil {
			return fmt.Errorf("failed to store  transactions update: %v", err)
		}

	}

	return nil
}

func (s *SmartContract) AcceptByPatient(ctx contractapi.TransactionContextInterface, userId string, requesterId string, disease string) error {
	// Retrieve pending requests
	pendingRequestedUsers, err := s.GetPendingRequestedUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending requested users: %v", err)
	}

	pendingRequesterUsers, err := s.GetPendingRequesterUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending requester users: %v", err)
	}

	// Check if the patient has pending requests for the specified disease
	requesters, exists := pendingRequestedUsers[userId][disease]
	if !exists || len(requesters) == 0 {
		return fmt.Errorf("no pending requests found for patient ID=%s and disease=%s", userId, disease)
	}

	// Check if the requester is valid
	if !contains(requesters, requesterId) {
		return fmt.Errorf("requester ID=%s has not requested disease=%s from patient ID=%s", requesterId, disease, userId)
	}

	// Retrieve patient data
	patient, err := s.GetPatient(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to retrieve patient data: %v", err)
	}

	// Get prescriptions for the requested disease
	prescriptions, hasPrescriptions := patient.Prescriptions[disease]
	if !hasPrescriptions || len(prescriptions) == 0 {
		return fmt.Errorf("patient ID=%s has no prescriptions for disease=%s", userId, disease)
	}

	// Retrieve accounts
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve accounts: %v", err)
	}

	// Validate the requester type and update their prescriptions
	userRole, exists := accounts[requesterId]
	if !exists {
		return fmt.Errorf("invalid requester ID=%s", requesterId)
	}

	switch userRole {
	case string(TPharmacyCompany):
		pharmacy, err := s.GetPharmacyCompany(ctx, requesterId)
		if err != nil {
			return fmt.Errorf("failed to retrieve pharmacy data: %v", err)
		}

		// Initialize the Prescriptions map if nil
		if pharmacy.Prescriptions == nil {
			pharmacy.Prescriptions = make(map[string][]string)
		}

		// Append prescriptions
		pharmacy.Prescriptions[disease] = append(pharmacy.Prescriptions[disease], prescriptions...)

		// Save updated PharmacyCompany data
		if err := s.putState(ctx, requesterId, pharmacy); err != nil {
			return fmt.Errorf("failed to update pharmacy data: %v", err)
		}

	case string(TMedicalResearchLab):
		lab, err := s.GetMedicalResearchLab(ctx, requesterId)
		if err != nil {
			return fmt.Errorf("failed to retrieve research lab data: %v", err)
		}

		// Initialize the Prescriptions map if nil
		if lab.Prescriptions == nil {
			lab.Prescriptions = make(map[string][]string)
		}

		// Append prescriptions
		lab.Prescriptions[disease] = append(lab.Prescriptions[disease], prescriptions...)

		// Save updated MedicalResearchLab data
		if err := s.putState(ctx, requesterId, lab); err != nil {
			return fmt.Errorf("failed to update research lab data: %v", err)
		}

	default:
		return fmt.Errorf("requester ID=%s is not a valid PharmacyCompany or MedicalResearchLab", requesterId)
	}

	// Remove the requester from pendingRequestedUsers
	pendingRequestedUsers[userId][disease] = removeUserFromList(requesters, requesterId)

	// If no more requesters left for this disease, remove the disease entry
	if len(pendingRequestedUsers[userId][disease]) == 0 {
		delete(pendingRequestedUsers[userId], disease)
	}

	// If no more diseases are left under the patient, remove the patient entry
	if len(pendingRequestedUsers[userId]) == 0 {
		delete(pendingRequestedUsers, userId)
	}

	// Remove patient from pendingRequesterUsers
	pendingRequesterUsers[requesterId][disease] = removeUserFromList(pendingRequesterUsers[requesterId][disease], userId)

	// If no more patients left under this disease, remove the disease entry
	if len(pendingRequesterUsers[requesterId][disease]) == 0 {
		delete(pendingRequesterUsers[requesterId], disease)
	}

	// If no more diseases are left for the requester, remove the requester entry
	if len(pendingRequesterUsers[requesterId]) == 0 {
		delete(pendingRequesterUsers, requesterId)
	}

	// Store updated states
	if err := s.putState(ctx, PendingRequestedUserKey, pendingRequestedUsers); err != nil {
		return fmt.Errorf("failed to update pending requested users: %v", err)
	}

	if err := s.putState(ctx, PendingRequesterUserKey, pendingRequesterUsers); err != nil {
		return fmt.Errorf("failed to update pending requester users: %v", err)
	}

	return nil
}



// Helper function to remove a user from a list
func removeUserFromList(users []string, userId string) []string {
	newUsers := []string{}
	for _, u := range users {
		if u != userId {
			newUsers = append(newUsers, u)
		}
	}
	return newUsers
}


func (s *SmartContract) RequestPatientData(ctx contractapi.TransactionContextInterface, userId string, disease string) error {
	// Retrieve all user accounts
	accounts, err := s.GetAccounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve accounts: %v", err)
	}

	// Validate user role
	userRole, exists := accounts[userId]
	if !exists || (userRole != string(TMedicalResearchLab) && userRole != string(TPharmacyCompany)) {
		return fmt.Errorf("user ID=%s is not authorized (must be a PharmacyCompany or MedicalResearchLab)", userId)
	}

	// Retrieve all patients
	patients, err := s.GetAllPatients(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve all patients: %v", err)
	}

	// Filter patients based on the requested disease
	var patientIDs []string
	for _, patientID := range patients {
		patient, err := s.GetPatient(ctx, patientID)
		if err != nil {
			continue 
		}

		if _, exists := patient.Prescriptions[disease]; exists {
			patientIDs = append(patientIDs, patient.PatientID) 
		}
	}

	if len(patientIDs) == 0 {
		return fmt.Errorf("no patients found with prescriptions for disease: %s", disease)
	}

	// Validate user role
	switch userRole {
	case string(TMedicalResearchLab):
		lab, err := s.GetMedicalResearchLab(ctx, userId)
		if err != nil || !lab.IsAdded {
			return fmt.Errorf("MedicalResearchLab with ID=%s not found or not added", userId)
		}

	case string(TPharmacyCompany):
		pharmacy, err := s.GetPharmacyCompany(ctx, userId)
		if err != nil || !pharmacy.IsAdded {
			return fmt.Errorf("PharmacyCompany with ID=%s not found or not added", userId)
		}

	default:
		return fmt.Errorf("only MedicalResearchLab or PharmacyCompany can execute this transaction")
	}

	// Retrieve existing pending requests
	pendingRequestedUsers, err := s.GetPendingRequestedUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending requested users: %v", err)
	}

	pendingRequesterUsers, err := s.GetPendingRequesterUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending requester users: %v", err)
	}

	// Ensure user entry exists in pendingRequesterUsers
	if _, exists := pendingRequesterUsers[userId]; !exists {
		pendingRequesterUsers[userId] = make(map[string][]string)
	}

	// Add only new patient IDs that are not already pending
	var newPatientIDs []string
	if existingPatients, exists := pendingRequesterUsers[userId][disease]; exists {
		for _, patientID := range patientIDs {
			if !contains(existingPatients, patientID) {
				newPatientIDs = append(newPatientIDs, patientID)
			}
		}
	} else {
		newPatientIDs = append(newPatientIDs, patientIDs...)
	}

	if len(newPatientIDs) == 0 {
		return fmt.Errorf("all selected patients for disease %s are already pending for user ID=%s", disease, userId)
	}

	// Update pendingRequesterUsers
	pendingRequesterUsers[userId][disease] = append(pendingRequesterUsers[userId][disease], newPatientIDs...)

	// Update pendingRequestedUsers
	for _, patientID := range newPatientIDs {
		if _, exists := pendingRequestedUsers[patientID]; !exists {
			pendingRequestedUsers[patientID] = make(map[string][]string)
		}
		pendingRequestedUsers[patientID][disease] = append(pendingRequestedUsers[patientID][disease], userId)
	}

	// Store the updated pending requests
	if err := s.putState(ctx, PendingRequestedUserKey, pendingRequestedUsers); err != nil {
		return fmt.Errorf("failed to store pending requested users: %v", err)
	}

	if err := s.putState(ctx, PendingRequesterUserKey, pendingRequesterUsers); err != nil {
		return fmt.Errorf("failed to store pending requester users: %v", err)
	}

	return nil
}

func main() {

	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(&SmartContract{})

	if err != nil {
		log.Panicf("error create  chaincode: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err := server.Start(); err != nil {
		log.Panicf("error starting  chaincode: %s", err)
	}
}
