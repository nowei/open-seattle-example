package testutils

import "github.com/nowei/open-seattle-example/server/internal/api"

func CreateTestDonationRegistration(name string, quantity int, distributionType api.DonationType, description *string) api.DonationRegistration {
	return api.DonationRegistration{
		Description: description,
		Name:        name,
		Quantity:    quantity,
		Type:        distributionType,
	}
}

func CreateTestDonationDistribution(donationId int, quantity int, distributionType api.DonationType, description *string) api.DonationDistribution {
	return api.DonationDistribution{
		DonationId:  donationId,
		Description: description,
		Quantity:    quantity,
		Type:        distributionType,
	}
}
