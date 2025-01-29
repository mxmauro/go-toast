//go:build cgo

package toast

// #cgo LDFLAGS: -lruntimeobject
// #include "toast_windows.h"
import "C"
import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"text/template"
	"unsafe"
)

// -----------------------------------------------------------------------------

const (
	WindowsOptionApplicationID       = "appID"
	WindowsOptionActivationType      = "activationType"
	WindowsOptionActivationArguments = "activationArgs"
	WindowsOptionDuration            = "duration"
	WindowsOptionTitle               = "title"
	WindowsOptionMessage             = "message"
	WindowsOptionIconPath            = "iconPath"
	WindowsOptionAudio               = "audio"
	WindowsOptionAudioLoop           = "audioLoop"
	WindowsOptionButtons             = "buttons"
)

const (
	WindowsDurationShort = "short"
	WindowsDurationLong  = "long"
)

const (
	WindowsAudioSilent         = "silent"
	WindowsAudioDefault        = "ms-winsoundevent:Notification.Default"
	WindowsAudioIM             = "ms-winsoundevent:Notification.IM"
	WindowsAudioMail           = "ms-winsoundevent:Notification.Mail"
	WindowsAudioReminder       = "ms-winsoundevent:Notification.Reminder"
	WindowsAudioSMS            = "ms-winsoundevent:Notification.SMS"
	WindowsAudioLoopingAlarm   = "ms-winsoundevent:Notification.Looping.Alarm"
	WindowsAudioLoopingAlarm2  = "ms-winsoundevent:Notification.Looping.Alarm2"
	WindowsAudioLoopingAlarm3  = "ms-winsoundevent:Notification.Looping.Alarm3"
	WindowsAudioLoopingAlarm4  = "ms-winsoundevent:Notification.Looping.Alarm4"
	WindowsAudioLoopingAlarm5  = "ms-winsoundevent:Notification.Looping.Alarm5"
	WindowsAudioLoopingAlarm6  = "ms-winsoundevent:Notification.Looping.Alarm6"
	WindowsAudioLoopingAlarm7  = "ms-winsoundevent:Notification.Looping.Alarm7"
	WindowsAudioLoopingAlarm8  = "ms-winsoundevent:Notification.Looping.Alarm8"
	WindowsAudioLoopingAlarm9  = "ms-winsoundevent:Notification.Looping.Alarm9"
	WindowsAudioLoopingAlarm10 = "ms-winsoundevent:Notification.Looping.Alarm10"
	WindowsAudioLoopingCall    = "ms-winsoundevent:Notification.Looping.Call"
	WindowsAudioLoopingCall2   = "ms-winsoundevent:Notification.Looping.Call2"
	WindowsAudioLoopingCall3   = "ms-winsoundevent:Notification.Looping.Call3"
	WindowsAudioLoopingCall4   = "ms-winsoundevent:Notification.Looping.Call4"
	WindowsAudioLoopingCall5   = "ms-winsoundevent:Notification.Looping.Call5"
	WindowsAudioLoopingCall6   = "ms-winsoundevent:Notification.Looping.Call6"
	WindowsAudioLoopingCall7   = "ms-winsoundevent:Notification.Looping.Call7"
	WindowsAudioLoopingCall8   = "ms-winsoundevent:Notification.Looping.Call8"
	WindowsAudioLoopingCall9   = "ms-winsoundevent:Notification.Looping.Call9"
	WindowsAudioLoopingCall10  = "ms-winsoundevent:Notification.Looping.Call10"
)

// -----------------------------------------------------------------------------

type WindowsButton struct {
	ActivationType string
	Label          string
	Arguments      string
}

// -----------------------------------------------------------------------------

var xmlTmpl *template.Template

// -----------------------------------------------------------------------------

func toastInit() error {
	var err error

	xmlTmpl, err = template.New("toast").Parse(`<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
	<visual>
		<binding template="ToastGeneric">
			{{if .IconPath}}
			<image placement="appLogoOverride" src="{{.IconPath}}" />
			{{end}}
			{{if .Title}}
			<text>{{.Title}}</text>
			{{end}}
			{{if .Message}}
			<text>{{.Message}}</text>
			{{end}}
		</binding>
	</visual>
	{{if ne .Audio "silent"}}
	<audio src="{{.Audio}}" loop="{{.AudioLoop}}" />
	{{else}}
	<audio silent="true" />
	{{end}}
	{{if .Buttons}}
	<actions>
		{{range .Buttons}}
		<action activationType="{{.ActivationType}}" content="{{.Label}}" arguments="{{.Arguments}}" />
		{{end}}
	</actions>
	{{end}}
</toast>`)
	if err != nil {
		return err
	}

	// Done
	return nil
}

func toastDone() {
}

func toastShow(opts Options) error {
	type toastXmlParameters struct {
		ActivationType      string
		ActivationArguments string
		Duration            string
		Title               string
		Message             string
		IconPath            string
		Audio               string
		AudioLoop           string
		Buttons             []WindowsButton
	}
	var appIdStrPtr *uint16
	var xmlTemplateStrPtr *uint16
	var toastXml bytes.Buffer

	// Get Application ID
	vString, err := getStringOption(opts, WindowsOptionApplicationID, true, nil)
	if err != nil {
		return err
	}
	if vString == nil || len(*vString) == 0 {
		return errors.New("invalid '" + WindowsOptionApplicationID + "' option")
	}
	appIdStrPtr, err = syscall.UTF16PtrFromString(*vString)
	if err != nil {
		return err
	}

	// Parse the rest of options
	tmplParams := toastXmlParameters{}

	// ActivationType
	vString, err = getStringOption(opts, WindowsOptionActivationType, false, addressOfString("protocol"))
	if err != nil {
		return err
	}
	tmplParams.ActivationType = strings.ToLower(*vString)
	if tmplParams.ActivationType != "foreground" && tmplParams.ActivationType != "background" && tmplParams.ActivationType != "protocol" {
		return errors.New("invalid '" + WindowsOptionActivationType + "' option")
	}

	// ActivationArguments
	vString, err = getStringOption(opts, WindowsOptionActivationArguments, false, addressOfString(""))
	if err != nil {
		return err
	}
	tmplParams.ActivationArguments, err = escapeXmlText(*vString)
	if err != nil {
		return err
	}

	// Duration
	vString, err = getStringOption(opts, WindowsOptionDuration, false, addressOfString("long"))
	if err != nil {
		return err
	}
	tmplParams.Duration = strings.ToLower(*vString)
	if tmplParams.Duration != "long" && tmplParams.Duration != "short" {
		return errors.New("invalid '" + WindowsOptionDuration + "' option")
	}

	// Title
	vString, err = getStringOption(opts, WindowsOptionTitle, false, addressOfString(""))
	if err != nil {
		return err
	}
	tmplParams.Title, err = escapeXmlText(*vString)
	if err != nil {
		return err
	}

	// Message
	vString, err = getStringOption(opts, WindowsOptionMessage, false, addressOfString(""))
	if err != nil {
		return err
	}
	tmplParams.Message, err = escapeXmlText(*vString)
	if err != nil {
		return err
	}

	// IconPath
	vString, err = getStringOption(opts, WindowsOptionIconPath, false, addressOfString(""))
	if err != nil {
		return err
	}
	tmplParams.IconPath, err = escapeXmlText(*vString)
	if err != nil {
		return err
	}

	// Audio
	vString, err = getStringOption(opts, WindowsOptionAudio, false, addressOfString("silent"))
	if err != nil {
		return err
	}
	tmplParams.Audio, err = escapeXmlText(*vString)
	if err != nil {
		return err
	}

	// AudioLoop
	if tmplParams.Audio != "silent" {
		vString, err = getStringOption(opts, WindowsOptionAudioLoop, false, addressOfString("false"))
		if err != nil {
			return err
		}
		tmplParams.AudioLoop = strings.ToLower(*vString)
		switch tmplParams.AudioLoop {
		case "1":
			fallthrough
		case "yes":
			tmplParams.AudioLoop = "true"
		case "true":

		case "0":
			fallthrough
		case "no":
			tmplParams.AudioLoop = "false"
		case "false":

		default:
			return errors.New("invalid '" + WindowsOptionAudioLoop + "' option")
		}
	} else {
		tmplParams.AudioLoop = "false"
	}

	tmplParams.Buttons, err = getButtonsArrayOption(opts)
	if err != nil {
		return err
	}

	// Build XML
	err = xmlTmpl.Execute(&toastXml, tmplParams)
	if err != nil {
		return err
	}
	xmlTemplateStrPtr, err = syscall.UTF16PtrFromString(toastXml.String())
	if err != nil {
		return err
	}

	// Show the toast
	hr := C.int_toast_show((*C.wchar_t)(unsafe.Pointer(appIdStrPtr)), (*C.wchar_t)(unsafe.Pointer(xmlTemplateStrPtr)))
	if hr != 0 {
		return fmt.Errorf("toastShow failed with HRESULT=0x%X", uint(hr))
	}

	// Done
	return nil
}

func escapeXmlText(input string) (string, error) {
	var escaped bytes.Buffer

	err := xml.EscapeText(&escaped, []byte(input))
	if err != nil {
		return "", err
	}
	return escaped.String(), nil
}

func getStringOption(opts Options, key string, mustExist bool, defaultValue *string) (*string, error) {
	v, ok := opts[key]
	if !ok {
		if mustExist {
			return nil, fmt.Errorf("missing option '%s'", key)
		}
		return defaultValue, nil
	}
	s, ok2 := v.(string)
	if !ok2 {
		return nil, fmt.Errorf("option '%s' is not a string", key)
	}
	return &s, nil
}

func getButtonsArrayOption(opts Options) ([]WindowsButton, error) {
	var src []WindowsButton
	var err error

	v, ok := opts[WindowsOptionButtons]
	if !ok {
		return nil, nil
	}
	switch v2 := v.(type) {
	case WindowsButton:
		src = []WindowsButton{v2}
	case []WindowsButton:
		src = v2
	default:
		return nil, errors.New("option '" + WindowsOptionButtons + "' is not an button struct or array of them")
	}

	ret := make([]WindowsButton, 0)
	for _, action := range src {
		fixedAction := WindowsButton{}

		// ActivationType
		if len(action.ActivationType) > 0 {
			fixedAction.ActivationType = strings.ToLower(action.ActivationType)
		} else {
			fixedAction.ActivationType = "protocol"
		}
		if fixedAction.ActivationType != "foreground" && fixedAction.ActivationType != "background" && fixedAction.ActivationType != "protocol" {
			return nil, errors.New("invalid ActivationType value in '" + WindowsOptionButtons + "' option")
		}

		// Arguments
		fixedAction.Arguments, err = escapeXmlText(action.Arguments)
		if err != nil {
			return nil, err
		}

		// Label
		fixedAction.Label, err = escapeXmlText(action.Label)
		if err != nil {
			return nil, err
		}

		ret = append(ret, fixedAction)
	}
	return ret, nil
}
