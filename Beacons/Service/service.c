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

#include "wsbind.h"
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
