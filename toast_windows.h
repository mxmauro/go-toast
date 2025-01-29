#include <windows.h>
#include <roapi.h>
#include <initguid.h>
#include <roapi.h>
#include <winstring.h>
#include <Windows.ui.notifications.h>
#include <stdio.h>

// -----------------------------------------------------------------------------

// ABI.Windows.UI.Notifications.IToastNotificationManagerStatics / 50ac103f-d235-4598-bbef-98fe4d1a3ad4
static const GUID UIID_IToastNotificationManagerStatics = { 0x50AC103F, 0xD235, 0x4598, { 0xBB, 0xEF, 0x98, 0xFE, 0x4D, 0x1A, 0x3A, 0xD4 } };

// ABI.Windows.Notifications.IToastNotificationFactory / 04124b20-82c6-4229-b109-fd9ed4662b53
static const GUID UIID_IToastNotificationFactory = { 0x04124B20, 0x82C6, 0x4229, { 0xB1, 0x09, 0xFD, 0x9E, 0xD4, 0x66, 0x2B, 0x53 } };

// ABI.Windows.Data.Xml.Dom.IXmlDocument / f7f3a506-1e87-42d6-bcfb-b8c809fa5494
static const GUID UIID_IXmlDocument = { 0xF7F3A506, 0x1E87, 0x42D6, { 0xBC, 0xFB, 0xB8, 0xC8, 0x09, 0xFA, 0x54, 0x94 } };

// ABI.Windows.Data.Xml.Dom.IXmlDocumentIO / 6cd0e74e-ee65-4489-9ebf-ca43e87ba637
static const GUID UIID_IXmlDocumentIO = { 0x6CD0E74E, 0xEE65, 0x4489, { 0x9E, 0xBF, 0xCA, 0x43, 0xE8, 0x7B, 0xA6, 0x37 } };

// -----------------------------------------------------------------------------

static int int_toast_initialize()
{
    return (int)RoInitialize(RO_INIT_MULTITHREADED);
}

static void int_toast_finalize()
{
    RoUninitialize();
    return;
}

static HRESULT int_toast_createXmlDocumentFromString(const wchar_t* xmlString, __x_ABI_CWindows_CData_CXml_CDom_CIXmlDocument** doc)
{
	HSTRING_HEADER docStrHdr, xmlStrHdr;
	HSTRING docStr, xmlStr;
	IInspectable* inspectable = NULL;
	__x_ABI_CWindows_CData_CXml_CDom_CIXmlDocumentIO* docIO = NULL;
	HRESULT hr;

	*doc = NULL;
	hr = WindowsCreateStringReference(RuntimeClass_Windows_Data_Xml_Dom_XmlDocument,
	                                  (UINT32)wcslen(RuntimeClass_Windows_Data_Xml_Dom_XmlDocument),
	                                  &docStrHdr, &docStr);
	if (FAILED(hr))
	{
		goto done;
	}
	if (docStr == NULL)
	{
		hr = E_POINTER;
		goto done;
	}

	hr = RoActivateInstance(docStr, &inspectable);
	if (SUCCEEDED(hr))
	{
		hr = inspectable->lpVtbl->QueryInterface(inspectable, &UIID_IXmlDocument, (void**)doc);
	}
	if (FAILED(hr))
	{
		goto done;
	}

	(*doc)->lpVtbl->QueryInterface((*doc), &UIID_IXmlDocumentIO, (void**)&docIO);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = WindowsCreateStringReference(xmlString, (UINT32)wcslen(xmlString), &xmlStrHdr, &xmlStr);
	if (FAILED(hr))
	{
		goto done;
	}
	if (xmlStr == NULL)
	{
		hr = E_POINTER;
		goto done;
	}

	hr = docIO->lpVtbl->LoadXml(docIO, xmlStr);

done:
	if (docIO != NULL)
	{
		docIO->lpVtbl->Release(docIO);
	}
	if (inspectable != NULL)
	{
		inspectable->lpVtbl->Release(inspectable);
	}

	if (FAILED(hr) && (*doc) != NULL)
	{
		(*doc)->lpVtbl->Release(*doc);
		*doc = NULL;
	}
	return hr;
}

static HRESULT int_toast_show(const wchar_t *szAppIdW, const wchar_t *szXmlTemplateW)
{
	HSTRING_HEADER appIdStrHdr, toastNotifMgrStrHdr, toastNotifStrHdr;
	HSTRING appIdStr, toastNotifMgrStr, toastNotifStr;
	__x_ABI_CWindows_CData_CXml_CDom_CIXmlDocument* inputXmlDoc = NULL;
	__x_ABI_CWindows_CUI_CNotifications_CIToastNotificationManagerStatics* toastStatics = NULL;
	__x_ABI_CWindows_CUI_CNotifications_CIToastNotifier* notifier = NULL;
	__x_ABI_CWindows_CUI_CNotifications_CIToastNotificationFactory* notifFactory = NULL;
	__x_ABI_CWindows_CUI_CNotifications_CIToastNotification* toast = NULL;
	HRESULT hr;

	hr = (int)RoInitialize(RO_INIT_MULTITHREADED);
	if (FAILED(hr))
	{
		return hr;
	}

	hr = WindowsCreateStringReference(szAppIdW, (UINT32)wcslen(szAppIdW), &appIdStrHdr, &appIdStr);
	if (FAILED(hr))
	{
		goto done;
	}
	if (appIdStr == NULL)
	{
		hr = E_POINTER;
		goto done;
	}

	hr = int_toast_createXmlDocumentFromString(szXmlTemplateW, &inputXmlDoc);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = WindowsCreateStringReference(RuntimeClass_Windows_UI_Notifications_ToastNotificationManager,
	                                  (UINT32)wcslen(RuntimeClass_Windows_UI_Notifications_ToastNotificationManager),
	                                  &toastNotifMgrStrHdr, &toastNotifMgrStr);
	if (FAILED(hr))
	{
		goto done;
	}
	if (toastNotifMgrStr == NULL)
	{
		hr = E_POINTER;
		goto done;
	}

	hr = RoGetActivationFactory(toastNotifMgrStr, &UIID_IToastNotificationManagerStatics, (void**)&toastStatics);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = toastStatics->lpVtbl->CreateToastNotifierWithId(toastStatics, appIdStr, &notifier);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = WindowsCreateStringReference(RuntimeClass_Windows_UI_Notifications_ToastNotification,
	                                  (UINT32)wcslen(RuntimeClass_Windows_UI_Notifications_ToastNotification),
	                                  &toastNotifStrHdr, &toastNotifStr);
	if (FAILED(hr))
	{
		goto done;
	}
	if (toastNotifStr == NULL)
	{
		hr = E_POINTER;
		goto done;
	}

	hr = RoGetActivationFactory(toastNotifStr, &UIID_IToastNotificationFactory, (void**)&notifFactory);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = notifFactory->lpVtbl->CreateToastNotification(notifFactory, inputXmlDoc, &toast);
	if (FAILED(hr))
	{
		goto done;
	}

	hr = notifier->lpVtbl->Show(notifier, toast);
	if (FAILED(hr))
	{
		goto done;
	}

	// Let's wait a bit until COM thread delivers the notification
	Sleep(1);

done:
	if (toast != NULL)
	{
		toast->lpVtbl->Release(toast);
	}
	if (notifFactory != NULL)
	{
		notifFactory->lpVtbl->Release(notifFactory);
	}
	if (notifier != NULL)
	{
		notifier->lpVtbl->Release(notifier);
	}
	if (toastStatics != NULL)
	{
		toastStatics->lpVtbl->Release(toastStatics);
	}
	if (inputXmlDoc != NULL)
	{
		inputXmlDoc->lpVtbl->Release(inputXmlDoc);
	}

	RoUninitialize();
	return hr;
}
