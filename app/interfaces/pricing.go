package interfaces

import "majesticcoding.com/app/models"

type PricingService interface {
	GetAll() ([]models.Pricing, error)
}
