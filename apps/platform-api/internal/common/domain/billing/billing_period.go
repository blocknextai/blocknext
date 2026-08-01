package billing

type BillingPeriod string

const (
	BillingPeriodLifetime BillingPeriod = "lifetime"
	BillingPeriodMonthly  BillingPeriod = "monthly"
	BillingPeriodAnnually BillingPeriod = "annually"
)

var (
	BillingPeriods = map[BillingPeriod]struct{}{
		BillingPeriodLifetime: {},
		BillingPeriodMonthly:  {},
		BillingPeriodAnnually: {},
	}
)

func (o BillingPeriod) String() string {
	return string(o)
}

func (o BillingPeriod) IsValid() bool {
	_, ok := BillingPeriods[o]
	return ok
}

func (o BillingPeriod) GetDuration() (years int, months int, isRecurring bool) {
	switch o {
	case BillingPeriodMonthly:
		return 0, 1, true
	case BillingPeriodAnnually:
		return 1, 0, true
	case BillingPeriodLifetime:
		return 0, 0, false
	default:
		return 0, 0, false
	}
}
