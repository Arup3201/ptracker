package domain

type PrivateProjectListed struct {
	*ProjectSummary
	Role string
}

type PublicProjectListed struct {
	*Project
	Role string
}

type JoinRequestListed struct {
	*JoinRequest
	*Member
}
