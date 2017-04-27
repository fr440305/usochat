#!/bin/bash

#You can use this bash file if you wanna run it safely.

#Usage:
#	<runcmd>;
#
#	<runcmd> = "urun.sh" <command> <flags>.
#	<command> = "start" | "qstart".


printHelp () {
	echo "";
	echo "USO is a simple chatting website, and this file,";
	echo "uso.sh, is the booter of this website's server.";
	echo "";
	echo "Usage:";
	echo "        bash uso.sh <command> <flags>";
	echo "";
	echo "If you wanna run it, type \`uso.sh start\`.";
	echo "";
};

build () {
	rm ./u.out
	go build -o ./u.out ./*.go;
};

createSpecPage () {
	html_page_header='
<!doctype html>
<html><head><title>USO Project Specification</title></head>
<body> ';
	html_page_tail='</body></html>';
	echo "" > ./frontend/spec.html;
	echo $html_page_header >> ./frontend/spec.html;
	echo $html_page_tail >> ./frontend/spec.html;
};
#The following function runs the uso.out.
#Usage:
run () {
	#echo $1;
	if [[ $1 == quite ]]; then {
		#echo $run_mode;
		#echo "--quite";
		./u.out 1>./u.std.log 2>./u.err.log;
	}; elif [[ $1 == noise ]]; then {
		#echo "--noise";
		./u.out;
	}; fi;
};

#The following function extracts the arguments of this bash,
#analize it, and ...
parseArgs () {
	command=$1;
	arg=$2;
	if [[ $command == "" ]]; then {
		#$ urun.sh;
		printHelp;
	}; else {
		build;
		echo The command is $command;
		if [[ $command == start ]]; then {
			if [[ $arg == quietly ]]; then {
				run quite;
			} elif [[ $arg == noisely || $arg == "" ]]; then {
				run noise;
			}; else {
				echo "type: $ bash uso.sh start [quitely | noisely];";
			}; fi;
		}; elif [[ $command == clean ]]; then {
			rm ./*.log;
			rm ./*.out;
			rm ./*.swp;
			echo "The work directory has been cleaned.";
			exit 0;
		}; elif [[ $command == doc ]]; then {
			echo TODO;
		}; elif [[ $command == spec ]]; then {
			# count the scale of this project.
			# TODO
			createSpecPage;
			echo TODO;
		}; else {
			echo "$command is an invalid command!";
		}; fi;
	}; fi;
};

parseArgs $1 $2 $3 $4 $5 $6;
