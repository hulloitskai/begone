package loader

// ChromeDriver URLs.
const (
	baseDriverURL  = "https://chromedriver.storage.googleapis.com/2.43/chromedriver_"
	LinuxDriverURL = baseDriverURL + "linux64.zip"
	WinDriverURL   = baseDriverURL + "win32.zip"
	MacDriverURL   = baseDriverURL + "mac64.zip"
)

// Verification hashes.
const (
	LinuxDriverHash = "1a67148288f4320e5125649f66e02962"
	WinDriverHash   = "d238c157263ec7f668e0ea045f29f1b7"
	MacDriverHash   = "84244e590d866294c1eaa5fa65ad51f1b0d2cc4fc1595aa3e414f5aac4da60ef"
)

// Loader constants.
const (
	DriverName    = "chromedriver"
	WinDriverName = "chromedriver.exe"
)
