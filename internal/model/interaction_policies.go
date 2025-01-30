package model

type DefaultInteractionPolicies struct {
	Direct   InteractionPolicy `json:"direct"`
	Private  InteractionPolicy `json:"private"`
	Public   InteractionPolicy `json:"public"`
	Unlisted InteractionPolicy `json:"unlisted"`
}

type InteractionPolicy struct {
	CanFavourite PolicyRules `json:"can_favourite"`
	CanReblog    PolicyRules `json:"can_reblog"`
	CanReply     PolicyRules `json:"can_reply"`
}

type PolicyRules struct {
	Always       []string `json:"always"`
	WithApproval []string `json:"with_approval"`
}
