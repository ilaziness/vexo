package sftp

import "testing"

func TestDecideResume(t *testing.T) {
	tests := []struct {
		name       string
		dest, src  int64
		wantAction ResumeAction
		wantOff    int64
	}{
		{"no dest", 0, 100, ResumeCopy, 0},
		{"missing dest", -1, 100, ResumeCopy, 0},
		{"partial", 40, 100, ResumeContinue, 40},
		{"complete", 100, 100, ResumeSkip, 100},
		{"dest larger", 120, 100, ResumeCopy, 0},
		{"empty source complete", 0, 0, ResumeCopy, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			act, off := DecideResume(tc.dest, tc.src)
			if act != tc.wantAction || off != tc.wantOff {
				t.Fatalf("DecideResume(%d,%d)=(%v,%d), want (%v,%d)",
					tc.dest, tc.src, act, off, tc.wantAction, tc.wantOff)
			}
		})
	}
}
