package archivercmd

import (
	"fmt"

	archiver "github.com/saricon83-sudo/Tyr365AdminCli/m365Archiver"
	"github.com/spf13/cobra"
)

// withArchiverClient instantiates an Archiver client and passes it to the handler.
// Keeps command RunE functions succinct while providing consistent error formatting.
func withArchiverClient(_ *cobra.Command, handler func(*archiver.ArchiverClient) error) error {
	client, err := archiver.NewArchiverClient()
	if err != nil {
		return fmt.Errorf("failed to create Archiver client: %w", err)
	}

	return handler(client)
}
