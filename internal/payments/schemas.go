package payments

type Plan int

const (
	PlanSupport Plan = iota
	Plan30
)

var planName = map[Plan]string{
	PlanSupport: "support",
	Plan30:      "30_days",
}

func (p Plan) String() string {
	return planName[p]
}

type InvoiceRequest struct {
	Plan Plan `json:"plan"`
}

type InvoiceResponse struct {
	InvoiceLink string `json:"invoice_link"`
}
