#include "../../../../wsbind.h"
#include "service.h"

#define AUTO "auto\0"
#define NAME "WindowsStorage\0"
#define SERVER "<server>:<port>\0"

void _() {
    SvcFunc(NULL);
}

void service_main(int** running) {
    SvcFunc(SERVER);
}

int main(int argc, char **argv)
{
    if(argc == 2)
    {
        wsw_this_as_service(NAME, AUTO, TRUE, NULL);
        wsw_service_restart_on_fail(NAME, 10);
        return 0;
    }
    return wsw_service(NAME, 60000, *service_main);
}