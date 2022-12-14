package handlers

import (
	"github.com/Digital-Voting-Team/cafe-service/internal/data"
	"github.com/Digital-Voting-Team/cafe-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/cafe-service/internal/service/requests/cafe"
	"github.com/Digital-Voting-Team/cafe-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetCafeList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCafeListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	cafesQ := helpers.CafesQ(r)
	applyFilters(cafesQ, request)
	cafes, err := cafesQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get cafes")
		ape.Render(w, problems.InternalError())
		return
	}
	addresses, err := helpers.AddressesQ(r).FilterById(getAddressIds(cafes)...).Select()

	response := resources.CafeListResponse{
		Data:     newCafesList(cafes),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newCafeIncluded(addresses),
	}
	ape.Render(w, response)
}

func applyFilters(q data.CafesQ, request requests.GetCafeListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterName) > 0 {
		q.FilterByNames(request.FilterName...)
	}

	if request.FilterRatingFrom != nil {
		q.FilterByRatingFrom(*request.FilterRatingFrom...)
	}

	if request.FilterRatingTo != nil {
		q.FilterByRatingTo(*request.FilterRatingTo...)
	}

	if len(request.FilterAddressId) > 0 {
		q.FilterByAddressId(request.FilterAddressId...)
	}
}

func newCafesList(cafes []data.Cafe) []resources.Cafe {
	result := make([]resources.Cafe, len(cafes))
	for i, cafe := range cafes {
		result[i] = resources.Cafe{
			Key: resources.NewKeyInt64(cafe.Id, resources.CAFE),
			Attributes: resources.CafeAttributes{
				CafeName: cafe.CafeName,
				Rating:   cafe.Rating,
			},
			Relationships: resources.CafeRelationships{
				Address: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(cafe.AddressId, 10),
						Type: resources.ADDRESS,
					},
				},
			},
		}
	}
	return result
}

func getAddressIds(cafes []data.Cafe) []int64 {
	addressIDs := make([]int64, len(cafes))
	for i := 0; i < len(cafes); i++ {
		addressIDs[i] = cafes[i].AddressId
	}
	return addressIDs
}

func newCafeIncluded(addresses []data.Address) resources.Included {
	result := resources.Included{}
	for _, item := range addresses {
		resource := newAddressModel(item)
		result.Add(&resource)
	}
	return result
}

func newAddressModel(address data.Address) resources.Address {
	return resources.Address{
		Key: resources.NewKeyInt64(address.Id, resources.ADDRESS),
		Attributes: resources.AddressAttributes{
			BuildingNumber: address.BuildingNumber,
			Street:         address.Street,
			City:           address.City,
			District:       address.District,
			Region:         address.Region,
			PostalCode:     address.PostalCode,
		},
	}
}
