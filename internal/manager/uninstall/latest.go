package uninstall

// LatestReceipt returns the newest applied uninstall receipt.
func (s *Service) LatestReceipt() (Receipt, error) {
	return LatestReceipt()
}
