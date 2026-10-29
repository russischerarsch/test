package worker

import (
	quotespb "PlataTest/contracts"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupClient() {
	grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	quotespb.NewQuotesClient()
}
