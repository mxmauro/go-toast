//go:build cgo

package toast_test

import (
	"testing"

	"github.com/mxmauro/go-toast"
)

// -----------------------------------------------------------------------------

func TestToast(t *testing.T) {
	err := toast.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer toast.Finalize()

	err = toast.Show(toast.Options{
		toast.WindowsOptionApplicationID: "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe",
		toast.WindowsOptionTitle:         "Demo",
		toast.WindowsOptionMessage:       "Hello World!",
		toast.WindowsOptionAudio:         toast.WindowsAudioLoopingAlarm,
		toast.WindowsOptionButtons: []toast.WindowsButton{
			{
				Label:     "Press me",
				Arguments: "",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
