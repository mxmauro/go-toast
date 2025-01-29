# go-toast

A Golang cross-platform library for sending desktop notifications. It requires [CGO](https://go.dev/wiki/cgo).

**NOTE**: Currently, only Microsoft Windows is supported. MacOS soon.

## Example

```golang
package main

import (
    "github.com/mxmauro/go-toast"
)

func main() {
    // Initialize the toast library
    err := toast.Initialize()
    if err != nil {
        panic(err)
    }
    defer toast.Finalize()

    // Show a toast
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
        panic(err)
    }
}
```

#### Toast options are platform dependant.

## LICENSE

[MIT](/LICENSE)
