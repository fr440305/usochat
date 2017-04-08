#urun.sh
#You can use this bash file if you wanna run it safely.

#Usage:
#
#	"urun" <command> <flags>.
#


echo "$BASH , $BASH_ARGC , ${BASH_ARGV[0]}, ${BASH_ARGV[1]}"

go build -o ./uso.out ./*.go
#( ./uso.out 1>./log/std.log 2>./log/err.log )


