package models

type GetURLRequest struct {
	ID     *ResourceID         `urlparam:"linkID"`
}
func (gr GetURLRequest) VerifyValues() error {
	return VerifyStruct(gr)
}

type GetLatestURLsRequest struct {
	Offset *StringifiedInt `urlparam:"offset"`
	Limit  *StringifiedInt `urlparam:"limit"`
}

func (grLatest GetLatestURLsRequest) VerifyValues() error {
	return VerifyStruct(grLatest)
}

type DeleteURLRequest struct {
	// UserID *UserID `json:"userID"`
	LinkID *ResourceID `json:"linkID"`
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
	LinkIDs     []*ResourceID    `json:"linkIDs"`
	Description *Description `json:"description"`
	UserID      *UserID      `json:"userID"`
}

func (pcr PostCollectionRequest) VerifyValues() error {
	return VerifyStruct(pcr)
}

type PutInCollectionRequest struct {
	LinkID *ResourceID `json:"linkID"`
	CollID *ResourceID `json:"collID"`
}

func (putcr PutInCollectionRequest) VerifyValues() error {
	return VerifyStruct(putcr)
}

type GetCollectionRequest struct {
	ID *ResourceID `urlparam:"collID"`
}

func (gcr GetCollectionRequest) VerifyValues() error {
	return VerifyStruct(gcr)
}
