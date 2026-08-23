package spec

// applyVelocityLive writes the new inlet speed onto the original
// Spec and returns that same pointer. Later doubling rules then
// start from an already-changed velocity.
func applyVelocityLive(s *Spec, v float64) *Spec {
	if s == nil {
		return nil
	}
	s.InletVelocityMPS = v
	return s
}
