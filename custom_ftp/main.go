package main

import (
	"flag"
	"fmt"
	"log"
	"main/client"
	"main/server"
)

/*
Задание:
    - [x] Вам необходимо разработать протокол передачи произвольного файла с одного компьютера на другой.
          Написать клиент и сервер, реализующие этот протокол.

    - [ ] Сервер также должен выводить скорость приёма данных от клиента.
    - [x] Серверу передаётся в параметрах номер порта, на котором он будет ждать входящих соединений от клиентов.

    Клиенту передаётся в параметрах относительный или абсолютный путь к файлу, который нужно отправить.
    Длина имени файла не превышает 4096 байт в кодировке UTF-8. Размер файла не более 1 терабайта.
    Клиенту также передаётся в параметрах DNS-имя (или IP-адрес) и номер порта сервера.
    Клиент передаёт серверу имя файла в кодировке UTF-8, размер файла и его содержимое.

    Для передачи используется TCP.
    Протокол передачи придумайте сами (т.е. программы разных студентов могут оказаться несовместимы).

    Сервер сохраняет полученный файл в поддиректорию uploads своей текущей директории.
    Имя файла, по возможности, совпадает с именем, которое передал клиент.
    Сервер никогда не должен писать за пределы директории uploads.

    В процессе приёма данных от клиента, сервер должен раз в 3 секунды выводить в консоль мгновенную скорость приёма и среднюю скорость за сеанс.
    Скорости выводятся отдельно для каждого активного клиента.
    Если клиент был активен менее 3 секунд, скорость всё равно должна быть выведена для него один раз.
    Под скоростью здесь подразумевается количество байтов переданных за единицу времени.

    После успешного сохранения всего файла сервер проверяет, совпадает ли размер полученных данных с размером, переданным клиентом,
    и сообщает клиенту об успехе или неуспехе операции, после чего закрывает соединение.

    Клиент должен вывести на экран сообщение о том, успешной ли была передача файла.
    Все используемые ресурсы ОС должны быть корректно освобождены, как только они больше не нужны.

    Сервер должен уметь работать параллельно с несколькими клиентами.
    Для этого необходимо использовать треды (POSIX threads или их аналог в вашей ОС).
    Сразу после приёма соединения от одного клиента, сервер ожидает следующих клиентов.
    В случае ошибки сервер должен разорвать соединение с клиентом. При этом он должен продолжить обслуживать остальных клиентов.
*/

type flags struct {
	host     string
	port     int
	server   bool
	client   bool
	filename string
}

var (
	f = flags{
		host:     "",
		port:     0,
		server:   false,
		client:   false,
		filename: "",
	}
)

func init() {
	prepareFlags()
}

func prepareFlags() {
	flag.StringVar(&f.host, "host", "localhost", "Хост. ")
	flag.IntVar(&f.port, "port", 8080, "Порт")
	flag.BoolVar(&f.server, "server", false, "Запуск программы в режиме сервера.")
	flag.BoolVar(&f.client, "client", false, "Запуск программы в режиме клиента.")
	flag.StringVar(&f.filename, "filename", "", "Имя передаваемого файла")
}

func main() {
	flag.Parse()

	if f.server && f.client {
		fmt.Println("Ошибка: Нельзя указать и -server, и -client одновременно.")
		return
	} else if !f.server && !f.client {
		fmt.Println("Ошибка: Вы должны указать -server или -client.")
		return
	}

	if f.server {
		fmt.Printf("Запуск сервера на %s:%d\n", f.host, f.port)
		s := server.NewServer(f.host, f.port, 120)
		log.Fatal(s.ListenAndServe())
	}

	if f.client {
		fmt.Printf("Запуск клиента. Запросы на сервер %s:%d\n", f.host, f.port)
		c := client.NewClient(f.host, f.port)
		log.Fatal(c.Send(f.filename))
	}
}
