package service

import (
	"Order_Service/internal/entity"
	_interface "Order_Service/internal/interface"

	"github.com/google/uuid"
)

type PartnerCommissionService struct {
	commissionRepository _interface.PartnerCommissionRepository
}

func NewPartnerCommissionService(commissionRepository _interface.PartnerCommissionRepository) *PartnerCommissionService {
	return &PartnerCommissionService{
		commissionRepository: commissionRepository,
	}
}

func (s *PartnerCommissionService) GetCommissionByPartnerID(partnerID uuid.UUID) (*[]entity.PartnerCommissionRespone, error) {
	return s.commissionRepository.GetCommissionByPartnerID(partnerID)
}

func (s *PartnerCommissionService) CreateCommission(commission *entity.PartnerCommission) (*entity.PartnerCommission, error) {
	createdCommission, err := s.commissionRepository.CreateCommission(commission)
	if err != nil {
		return nil, err
	}
	return createdCommission, nil
}
func (s *PartnerCommissionService) GetCommissionByOrderItemID(orderItemID uuid.UUID) (*entity.PartnerCommissionRespone, error) {
	return s.commissionRepository.GetCommissionByOrderItemID(orderItemID)
}
