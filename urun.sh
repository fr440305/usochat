#!/bin/bash

#urun.sh
#You can use this bash file if you wanna run it safely.

#Usage:
#	<runcmd>;
#
#	<runcmd> = "urun.sh" <command> <flags>.
#	<command> = "start" | "qstart".


printHelp () {
	echo "";
	echo "USO is a simple chatting website, and this file,";
	echo "urun.sh, is the booter of this website's server.";
	echo "";
	echo "Usage:";
	echo "        [./]urun.sh <command> <flags>";
	echo "";
	echo "If you wanna run it, type \`urun.sh start\`.";
	echo "";
};

build () {
	go build -o ./uso.out ./*.go;
};

#The following function runs the uso.out.
#Usage:
run () {
	#echo $1;
	if [[ $1 == quite ]]; then {
		#echo $run_mode;
		#echo "--quite";
		./uso.out 1>./u.std.log 2>./u.err.log;
	}; elif [[ $1 == noise ]]; then {
		#echo "--noise";
		./uso.out;
	}; fi;
};

#The following function extracts the arguments of this bash,
#analize it, and ...
parseArgs () {
	if [[ $BASH_ARGC -eq 0 ]]; then {
		#$ urun.sh;
		printHelp;
	}; else {
		build;
		command=${BASH_ARGV[$(($BASH_ARGC-1))]};
		echo The command is $command;
		if [[ $command == start ]]; then {
			run noise;
		}; elif [[ $command == qstart ]]; then {
			run quite;
		}; elif [[ $command == clean ]]; then {
			rm ./*.log;
			rm ./*.out;
			rm ./*.swp;
			echo "The work directory has been cleaned.";
			exit 0;
		}; else {
			echo $command is an invalid command!;
		}; fi;
	}; fi;
};

parseArgs;
