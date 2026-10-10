package sftp

// ResumeAction describes how to handle a partial transfer.
type ResumeAction int

const (
	ResumeCopy     ResumeAction = iota // transfer from offset 0 (truncate/create)
	ResumeContinue                     // seek to offset and continue
	ResumeSkip                         // already complete
)

// DecideResume returns action and byte offset for a single-file transfer.
// destSize is the size of the partial destination (remote for upload, local for download).
// sourceSize is the full source size.
func DecideResume(destSize, sourceSize int64) (ResumeAction, int64) {
	if sourceSize < 0 {
		sourceSize = 0
	}
	if destSize <= 0 {
		return ResumeCopy, 0
	}
	if destSize == sourceSize {
		return ResumeSkip, destSize
	}
	if destSize < sourceSize {
		return ResumeContinue, destSize
	}
	return ResumeCopy, 0
}
