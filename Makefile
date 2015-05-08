test: targa .PHONY
	./targa | tgatoppm | pnmtopng > test.png

targa: main.go
	go build

.PHONY:
