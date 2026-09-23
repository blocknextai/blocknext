package contract

import (
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type RuntimeConfig = triggersDomainTriggers.RuntimeConfig
type Trigger = triggersDomainTriggers.Trigger

const TriggerTypeSchedule = triggersDomainTriggers.TriggerTypeSchedule
const TriggerTypeWebhook = triggersDomainTriggers.TriggerTypeWebhook

type TriggerService = triggersApplicationTriggers.TriggerService
