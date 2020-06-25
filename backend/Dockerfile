# Utilizamos la imagen oficial de GO (golang)
FROM golang
# Es importante llevarnos los archivos .go para hacer correr el server
COPY . /backend
# Nos situamos en la carpeta de backend
WORKDIR /backend
# Hacemos público el puerto 5000 para poder mapearlo en el host
EXPOSE 5000
# Aquí vamos a instalar el comando whois debido a que el server lo usa, también
# importamos los módulos de GO que usa el server y finalmente buildeamos
RUN apt update && apt install whois && go get github.com/lib/pq && go get github.com/valyala/fasthttp && go get github.com/PuerkitoBio/goquery && go get github.com/buaazp/fasthttprouter && go build
# Su comando principal es estar ejecutando el binario correspondiente al servidor
CMD ["./backend"]