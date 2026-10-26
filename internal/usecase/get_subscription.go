package usecase

type xuiService interface {
	GetSubscriptionsByTgID(telegramID int64) ([]string, error)
}

type GetSubscriptionUseCase struct {
	xuiService xuiService
}

// Execute implements [handlers.GetSubscriptionUseCase].
func (g *GetSubscriptionUseCase) Execute(telegramID int64) ([]string, error) {
	return g.xuiService.GetSubscriptionsByTgID(telegramID)
}

func NewGetSubscriptionUseCase(xuiService xuiService) *GetSubscriptionUseCase {
	return &GetSubscriptionUseCase{
		xuiService: xuiService,
	}
}
