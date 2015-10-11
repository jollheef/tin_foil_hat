#!/usr/bin/python3
# Пример сервиса
#
###### Обычное взаимодействие #####
## Регистрация
# REG
# логин
# пароль
#
## Положить флаг
# PUT
# логин
# пароль
# флаг
#
## Получить флаг
# GET
# логин
# пароль
# флаг
#
###### Тестовые функции #####
## Очистить флаги (для проверки CORRUPT)
# CLEAR
#
## Нарушить логику регистрации (MUMBLE)
# REGFAIL
## Вернуть логику регистрации (MUMBLE)
# REGOK
## Нарушить логику хранения данных (MUMBLE, CORRUPT для get)
# DATAFAIL
#
## Перестать отвечать на запросы в течении 10 секунд (MUMBLE)
# SLEEP
#
###### Памятка по статусам #####
# OK -- сервис онлайн, обрабатывает запросы, получает и отдает флаги.
# MUMBLE -- сервис онлайн, но некорректно работает
# CORRUPT -- сервис онлайн, но установленные флаги невозможно получить.
# DOWN -- сервис оффлайн.

from socket import socket, SO_REUSEADDR
from time import sleep
from signal import signal, SIGTERM
from sys import argv

BUFFER_SIZE = 1024
BACKLOG_SIZE = 10

users = dict()
database = dict()
regfail = False
datafail = False

sock = socket(SO_REUSEADDR)

sock.bind(('', int(argv[1])))

sock.listen(BACKLOG_SIZE)

def gracefully_exit(*args):
    global sock
    sock.close()

signal(SIGTERM, gracefully_exit)

def Handler(conn, addr):
    global users, database, regfail, datafail

    cmd = conn.recv(BUFFER_SIZE)

    if b'SLEEP\n' == cmd:
        sleep(10)

    if b'CLEAR\n' == cmd:
        database = dict()

    elif b'REGFAIL\n' == cmd:
        regfail = True

    elif b'REGOK\n' == cmd:
        regfail = False

    elif b'DATAFAIL\n' == cmd:
        datafail = False

    elif b'REG\n' == cmd:
        if regfail:
            conn.close()
            return

        conn.send(b'OK\n')

        login = conn.recv(BUFFER_SIZE)
        if login in users:
            conn.send(b'EXIST\n')
            conn.close()
            return
        else:
            conn.send(b'OK\n')

        password = conn.recv(BUFFER_SIZE)

        users[login] = password

        conn.send(b'OK\n')


    elif b'PUT\n' == cmd or b'GET\n' == cmd:
        conn.send(b'OK\n')

        login = conn.recv(BUFFER_SIZE)

        if login in users:
            conn.send(b'OK\n')
        else:
            conn.send(b'INCORRECT\n')
            conn.close()
            return

        password = conn.recv(BUFFER_SIZE)
        if password != users[login]:
            conn.send(b'INCORRECT\n')
            conn.close()
            return
        else:
            conn.send(b'OK\n')

        if b'PUT\n' == cmd:
            data = conn.recv(BUFFER_SIZE)
            database[login] = data

            if datafail:
                database[login] += '_'

            conn.send(b'OK\n')

        elif b'GET\n' == cmd:
            if login in database:
                sleep(0.1)
                conn.send(database[login] + b'\nOK\n')
            else:
                conn.send(b'NOOK\n')
                conn.close()
                return

    else:
        conn.send(b'NOOK\n')

    conn.close()


try:
    while True:
        conn, addr = sock.accept()
        Handler(conn, addr)
except Exception as e:
    print(e)
    sock.close()
