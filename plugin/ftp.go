package plugin

// FtpBurst burst info for ftp
var FtpBurst BurstCell

func burstFtp()

func init() {
	FtpBurst = NewBurstCell("ftp")
}
