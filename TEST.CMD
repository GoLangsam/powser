@Echo.

go generate  ./...
@Echo.
@Echo compare dch/dch.* - 4 Differences
@Echo.

@pause

go fmt  ./...
go vet  ./...
go test	./...
go doc -all	> godoc.md

@pause

cd cmd

cd ps
go build .
ps.exe
@pause
cd ..

cd ..

@Echo.
@Echo finished :-)

@pause
