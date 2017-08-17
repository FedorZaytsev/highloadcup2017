# Наследуемся от CentOS 7
FROM centos:7

# Выбираем рабочую папку
WORKDIR /root

# Устанавливаем wget и скачиваем Go
RUN yum install -y wget && \
    wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz

# Устанавливаем Go, создаем workspace и папку проекта
RUN tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg && \
    mkdir go/src/dumb

RUN yum install -y epel-release

# Redis
RUN yum install -y redis

# Redis start
#Run "redis-server"
#CMD ["redis-cli", "ping"]

# Задаем переменные окружения для работы Go
ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go

# Копируем наш исходный main.go внутрь контейнера, в папку go/src/dumb
ADD . / go/src/dumb/

WORKDIR /root/go/src/dumb

# Компилируем и устанавливаем наш сервер
RUN go build

# Открываем 80-й порт наружу
EXPOSE 80

# Запускаем наш сервер
CMD redis-server & ./dumb