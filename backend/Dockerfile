FROM golang

COPY . /backend

WORKDIR /backend

EXPOSE 5000

RUN apt update && apt install whois && go get github.com/lib/pq && go get github.com/valyala/fasthttp && go get github.com/PuerkitoBio/goquery && go get github.com/buaazp/fasthttprouter && go build

CMD ["./backend"]