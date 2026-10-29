package worker

import (
	quotespb "PlataTest/contracts"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupClient() (quotespb.QuotesClient, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := quotespb.NewQuotesClient(conn)
	return client, nil
}
