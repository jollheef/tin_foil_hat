#!/usr/bin/python3
###### Памятка по статусам #####
# OK -- сервис онлайн, обрабатывает запросы, получает и отдает флаги.
# MUMBLE -- сервис онлайн, но некорректно работает
# CORRUPT -- сервис онлайн, но установленные флаги невозможно получить.
# DOWN -- сервис оффлайн.

from sys import argv
from socket import socket, AF_INET, SOCK_STREAM
from string import ascii_letters
from random import randint, shuffle
from time import sleep

from tinfoilhat import Checker, \
    ServiceMumbleException,     \
    ServiceCorruptException,    \
    ServiceDownException

class DummyChecker(Checker):

    BUFSIZE = 1024

    """
    Сгенерировать логин

    @return строка логина из 10 символов английского алфавита
    """
    def random_login(self):
        symbols = list(ascii_letters)
        shuffle(symbols)
        return ''.join(symbols[0:10])

    """
    Сгенерировать пароль

    @return строка пароля из 10 цифр
    """
    def random_password(self):
        return str(randint(100500**2, 100500**3))[0:10]

    """
    Отправить логин и пароль сервису.

    @param sock сокет
    @param login логин
    @param password пароль
    """
    def send_cred(self, s, login, password):
        s.send(login.encode('utf-8'))
        if b'OK\n' != s.recv(self.BUFSIZE):
            raise ServiceMumbleException()
        s.send(password.encode('utf-8'))
        if b'OK\n' != s.recv(self.BUFSIZE):
            raise ServiceMumbleException()

    """
    Положить флаг в сервис

    @param host адрес хоста
    @param port порт сервиса
    @param flag флаг
    @return состояние, необходимое для получения флага
    """
    def put(self, host, port, flag):
        try:
            s = socket(AF_INET, SOCK_STREAM)

            s.connect((host, port))

            s.send(b'REG\n')

            if b'OK\n' != s.recv(self.BUFSIZE):
                raise ServiceMumbleException()

            login = self.random_login()
            password = self.random_password()

            self.send_cred(s, login, password)

            s.close()

            s = socket(AF_INET, SOCK_STREAM)

            s.connect((host, port))

            s.send(b'PUT\n')

            if b'OK\n' != s.recv(self.BUFSIZE):
                raise ServiceMumbleException()

            self.send_cred(s, login, password)

            s.send(flag.encode('utf-8'))

            if b'OK\n' != s.recv(self.BUFSIZE):
                raise ServiceMumbleException()

            return login + ":" + password

        except OSError as e:
            if e.errno == 111:  # ConnectionRefusedError
                raise ServiceDownException()
            else:
                raise ServiceMumbleException()

    """
    Получить флаг из сервиса

    @param host адрес хоста
    @param port порт сервиса
    @param state состояние
    @return флаг
    """
    def get(self, host, port, state):
        login, password = state.split(':')

        s = socket(AF_INET, SOCK_STREAM)

        s.connect((host, port))

        s.send(b'GET\n')

        if b'OK\n' != s.recv(self.BUFSIZE):
            raise ServiceMumbleException()

        try:
            self.send_cred(s, login, password)
        except ServiceMumbleException:
            raise ServiceCorruptException()

        try:
            flag, ret = s.recv(self.BUFSIZE).split()
            return flag.decode('utf-8')
        except ValueError:
            raise ServiceCorruptException()

    """
    Проверить состояние сервиса

    @param host адрес хоста
    @param port порт сервиса
    """
    def chk(self, host, port):
        # Так как сервис реализует только логику хранилища,
        # её и проверяем.
        # Это отличается от put и get тем, что происходит в один момент,
        # тем самым наличие данных по прошествии времени не проверяется.
        data = self.random_password()
        state = self.put(host, port, data)
        new_data = self.get(host, port, state)

        if data != new_data:
            raise ServiceMumbleException()

if __name__ == '__main__':
    DummyChecker(argv)
