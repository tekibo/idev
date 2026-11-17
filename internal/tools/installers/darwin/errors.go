package darwin

import "errors"

var (
	// ErrPayloadMissing reports that idev setup has not installed the Apple tools payload.
	ErrPayloadMissing = errors.New("Apple tools payload not installed; run idev setup")
	// ErrSDKMissing reports that xtool has not created the private Darwin SDK.
	ErrSDKMissing = errors.New("private Darwin SDK not found; run idev sdk install <Xcode.xip>")
	// ErrSwiftMissing reports that the private Swift toolchain is absent.
	ErrSwiftMissing = errors.New("private Swift toolchain not installed; run idev setup --swift")
	// ErrPythonMissing reports that python3 is absent.
	ErrPythonMissing = errors.New("python3 not found; run idev system plan")
	// ErrToolMissing reports a missing build product.
	ErrToolMissing = errors.New("Apple tool build output is missing")
)
