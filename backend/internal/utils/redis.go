package utils

func GetBucketKey(userId string) string {
	return "bucket:user:" + userId
}

func GetAccessTokenKey(sessionId string) string {
	return "access_token:session" + sessionId
}

func GetVerifierKey(state string) string {
	return "pkce_verifier:state:" + state
}
