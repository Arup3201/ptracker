package domain

type PrivateProjectListed struct {
	*ProjectSummary
	Role string
}

type PublicProjectListed struct {
	*Project
	Role string
}
