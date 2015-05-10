test.png: anytotarga
	./anytotarga < basn3p04.png | tgatoppm | pnmtopng > test.png


anytotarga: main.go
	go fmt
	go build
