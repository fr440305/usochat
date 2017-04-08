go build -o ./uso.out ./*.go
( ./uso.out 1>./log/std.log 2>./log/err.log )
