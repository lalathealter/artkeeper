package models

type GetURLRequest struct {
	ID     *LinkID         `urlparam:"id"`
	Offset *StringifiedInt `urlparam:"offset"`
	Limit  *StringifiedInt `urlparam:"limit"`
	// client string
}

type DeleteURLRequest struct {
	// UserID *UserID `json:"userID"`
	LinkID *LinkID `json:"linkID"`
}

func (dr DeleteURLRequest) VerifyValues() error {
	return VerifyStruct(dr)
}

type PostURLRequest struct {
	Link        *InputLink   `json:"link"`
	Description *Description `json:"description"`
	UserID      *UserID      `json:"userID"`
}

func (pr PostURLRequest) VerifyValues() error {
	return VerifyStruct(pr)
}

type PostCollectionRequest struct {
	LinkIDs     []*LinkID    `json:"linkIDs"`
	Description *Description `json:"description"`
	UserID      *UserID      `json:"userID"`
}

func (pcr PostCollectionRequest) VerifyValues() error {
	return VerifyStruct(pcr)
}
