package spec

import "fmt"

// stringifyValidErr flattens *ValidateError into a plain error whose
// text no longer names the offending field. Callers that branch on
// cylinder_diameter_m / inlet_velocity_mps then see only a generic line.
func stringifyValidErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("算例不合法")
}
