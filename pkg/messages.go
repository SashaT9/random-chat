package protocol

func RegionRequest(region string) ([]byte, error) {
	payload := RegionRequestPayload{Region: region}
	return Marshal(TypeRegionRequest, payload)
}

func RegionCount(count int) ([]byte, error) {
	payload := RegionCountPayload{Count: count}
	return Marshal(TypeRegionCount, payload)
}
