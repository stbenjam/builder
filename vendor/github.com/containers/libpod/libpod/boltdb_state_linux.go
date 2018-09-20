// +build linux

package libpod

import (
	"github.com/sirupsen/logrus"
)

// replaceNetNS handle network namespace transitions after updating a
// container's state.
func replaceNetNS(netNSPath string, ctr *Container, newState *containerState) error {
	if netNSPath != "" {
		// Check if the container's old state has a good netns
		if ctr.state.NetNS != nil && netNSPath == ctr.state.NetNS.Path() {
			newState.NetNS = ctr.state.NetNS
		} else {
			// Close the existing namespace.
			// Whoever removed it from the database already tore it down.
			if err := ctr.runtime.closeNetNS(ctr); err != nil {
				return err
			}

			// Open the new network namespace
			ns, err := joinNetNS(netNSPath)
			if err == nil {
				newState.NetNS = ns
			} else {
				logrus.Errorf("error joining network namespace for container %s", ctr.ID())
				ctr.valid = false
			}
		}
	} else {
		// The container no longer has a network namespace
		// Close the old one, whoever removed it from the DB should have
		// cleaned it up already.
		if err := ctr.runtime.closeNetNS(ctr); err != nil {
			return err
		}
	}
	return nil
}

// getNetNSPath retrieves the netns path to be stored in the database
func getNetNSPath(ctr *Container) string {
	if ctr.state.NetNS != nil {
		return ctr.state.NetNS.Path()
	}
	return ""
}
