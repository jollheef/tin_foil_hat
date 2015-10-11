#!/usr/bin/python3

from sys import stderr

STATUS_CHECKER_ERROR = 1
STATUS_SERVICE_MUMBLE = 2
STATUS_SERVICE_CORRUPT = 3
STATUS_SERVICE_DOWN = 4

class ServiceMumbleException(Exception):
    pass

class ServiceCorruptException(Exception):
    pass

class ServiceDownException(Exception):
    pass

class NonImplementedException(Exception):
    pass

def error(s):
    print(s, file=stderr)

class Checker(object):

    def usage(self):
        error("Usage:")
        error("\tput HOST PORT FLAG\tПоложить флаг в сервис. Возвращает состояние.")
        error("\tget HOST PORT STATE\tПолучить флаг из сервиса для состояния.")
        error("\tchk HOST PORT\tПроверить доступность и целостность сервиса.")

    def __init__(self, argv):
        if len(argv) < 3:
            self.usage()
            exit(STATUS_CHECKER_ERROR)

        try:
            error(argv)
            cmd = argv[1]
            host = argv[2]
            port = int(argv[3])

            if "put" == cmd:
                if len(argv) < 5:
                    error("Недостаточно аргументов.")
                    exit(STATUS_CHECKER_ERROR)
                flag = argv[4]
                error('Put flag \'' + str(flag) + '\' to '
                      + str(host) + ':' + str(port))
                print(self.put(host, port, flag))

            elif "get" == cmd:
                if len(argv) < 5:
                    error("Недостаточно аргументов.")
                    exit(STATUS_CHECKER_ERROR)
                state = argv[4]
                error('Get flag from ' + str(host) + ':' + str(port))
                print(self.get(host, port, state))

            elif "chk" == cmd:
                error('Check status of ' + str(host) + ':' + str(port))
                self.chk(host, port)

            else:
                self.usage()
                exit(STATUS_CHECKER_ERROR)

        except ServiceMumbleException:
            exit(STATUS_SERVICE_MUMBLE)

        except ServiceCorruptException:
            exit(STATUS_SERVICE_CORRUPT)

        except ServiceDownException:
            exit(STATUS_SERVICE_DOWN)

    """
    Положить флаг в сервис

    @param host адрес хоста
    @param port порт сервиса
    @param flag флаг
    @return состояние, необходимое для получения флага
    """
    def put(self, host, port, flag):
        raise NonImplementedException()

    """
    Получить флаг из сервиса

    @param host адрес хоста
    @param port порт сервиса
    @param state состояние
    @return флаг
    """
    def get(self, host, port, state):
        raise NonImplementedException()

    """
    Проверить состояние сервиса

    @param host адрес хоста
    @param port порт сервиса
    """
    def chk(self, host, port):
        raise NonImplementedException()

if __name__ == '__main__':
    error("Interactive run not supported")
