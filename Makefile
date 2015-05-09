test: targa .PHONY
	./anytotarga < basn3p04.png | tgatoppm | pnmtopng > test.png


targa: main.go
	go fmt
	go build

.PHONY:
