package handlers

import (
	"github.com/Digital-Voting-Team/cafe-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/cafe-service/internal/service/requests/cafe"
	"github.com/Digital-Voting-Team/cafe-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetCafe(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCafeRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resultCafe, err := helpers.CafesQ(r).FilterById(request.CafeId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get cafe from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if resultCafe == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	relateAddress, err := helpers.AddressesQ(r).FilterById(resultCafe.AddressId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get address")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var includes resources.Included
	includes.Add(&resources.Address{
		Key: resources.NewKeyInt64(relateAddress.Id, resources.ADDRESS),
		Attributes: resources.AddressAttributes{
			BuildingNumber: relateAddress.BuildingNumber,
			Street:         relateAddress.Street,
			City:           relateAddress.City,
			District:       relateAddress.District,
			Region:         relateAddress.Region,
			PostalCode:     relateAddress.PostalCode,
		},
	})

	result := resources.CafeResponse{
		Data: resources.Cafe{
			Key: resources.NewKeyInt64(resultCafe.Id, resources.CAFE),
			Attributes: resources.CafeAttributes{
				CafeName: resultCafe.CafeName,
				Rating:   resultCafe.Rating,
			},
			Relationships: resources.CafeRelationships{
				Address: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultCafe.AddressId, 10),
						Type: resources.ADDRESS,
					},
				},
			},
		},
		Included: includes,
	}
	ape.Render(w, result)
}
