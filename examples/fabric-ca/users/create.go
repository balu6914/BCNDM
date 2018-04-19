package users

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

// Add new user to fabric network
func Create(bc *BcNetwork) (response string, err error) {
	sdk := bc
	println(sdk)

	ctxProvider := bc.sdk.Context()

	mspClient, err := msp.New(bc)

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name: "test",
		// Affiliation is mandatory. "org1" and "org2" are hardcoded as CA defaults
		// See https://github.com/hyperledger/fabric-ca/blob/release/cmd/fabric-ca-server/config.go
		Affiliation: "org1",
		Secret:      "12345",
	})

	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Enroll the new user
	err = mspClient.Enroll(username, msp.WithSecret(enrollmentSecret))

	if err != nil {
		t.Fatalf("Enroll failed: %v", err)
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		t.Fatalf("GetSigningIdentity failed: %v", err)
	}

	return "Create callback", u
}
