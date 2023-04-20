// gcc critical.c -o critical.exe

#include <windows.h>
#include <stdio.h>

typedef VOID (_stdcall *RtlSetProcessIsCritical) (IN BOOLEAN NewValue, OUT PBOOLEAN OldValue /* optional */, IN BOOLEAN IsWinlogon);

// require UAC elevation
int main(void) {
    HANDLE hToken;
	LUID seDebugName;
	TOKEN_PRIVILEGES tkp;
	if (!OpenProcessToken(GetCurrentProcess(), TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY, &hToken))
	{
	    printf("OpenProcessToken failed\n");
		return 1;
	}
	if (!LookupPrivilegeValue(NULL, SE_DEBUG_NAME, &seDebugName))
	{
	    printf("LookupPrivilegeValue failed\n");
		CloseHandle(hToken);
		return 1;
	}
	tkp.PrivilegeCount = 1;
	tkp.Privileges[0].Luid = seDebugName;
	tkp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;
	if (!AdjustTokenPrivileges(hToken, FALSE, &tkp, sizeof(tkp), NULL, NULL))
	{
	    printf("AdjustTokenPrivileges failed\n");
		return 1;
	}
	CloseHandle(hToken);

    HMODULE hNtdll = GetModuleHandle(TEXT("ntdll.dll"));
    RtlSetProcessIsCritical pRtlSetProcessIsCritical = (RtlSetProcessIsCritical)GetProcAddress(hNtdll, "RtlSetProcessIsCritical");
    if (pRtlSetProcessIsCritical == 0) {
        printf("GetProcAddress failed\n");
        return 1;
    }

    printf("Going to sleep...\n");
    pRtlSetProcessIsCritical(1, 0, 0);
    HANDLE hThread = OpenThread(THREAD_ALL_ACCESS, FALSE, GetCurrentThreadId());
    SuspendThread(hThread);

    printf("Awaken!\n");
    CloseHandle(hThread);
    return 0;
}
