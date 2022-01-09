package services

type GetAllConfig struct {
	Limit          uint
	Offset         uint
	IncludeDeleted bool
	OnlyDeleted    bool
	Expand         []string
	Order          []string
	// op is one of  =, !=, <, <=, >, >=
	Query []string
}
