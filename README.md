# USO - Rewrite - prelease.

Uso is a simple online chatting platfrom.
You can go to (here)[http://45.118.135.218:9999/]
to try, but it does not work all the time.

## Prog Construction

I will describe the structure of this platform
in three perspectives: feature, modules,
and status.
(用例功能、模块结构、状态转换)
 
### Overview

This platform is developed by using
HTML/JavaScript as front-end and Golang
as server.

The front-end and server will be connected
by using WebSocket.

### Features

```
[Client] <-- JSON --> [Node]
[Client] <-- JSON --> [Node]

      ... ... ... ...

[Client] <--(JSON)--> [Node]
                       |  ^
                (Msg)  |  |  (Msg)
                       V  |
                     [Center]
```

### Modules

### Status

### Dependencies

The dependent package(s) is/are:

-	github.com/gorilla/websocket
--	An implementation of WebSocket in Golang.


## Code Style

The Types will be named in a single word with Capital Initial.
(So the name of the type must has only one word.)
eg:

```Go
type Node struct {...}
type Center struct {...}
```

The methods and the constructors will be coded in camelStyle.
Use as more than one word as possible.
eg:
```Go
func newCenter () *Center
func (C *Center) newNode () *Node
func (c *Center) boardcast () error
func (n *Node) sendMsgToCenter () error
```

The tool functions (utilities) will be named in canmelStyle,
with a underscore prefix (stands for its father-class is ```_```,
except for main function).
eg:
```Go
func _strToMsg (source string) ([]byte, err)
```

The variables inside the function or method will be named
in underscore style. Only lowercase.
eg:
```Go
func (S *Sample) showExample () {
	var sample_var = 0
	var no_uppercase_letter = 9
	fmt.Prinln(sample_var, no_uppercase_letter)
}
```

The constants will be named ALL_UPPERCASE style. Use more than
one word as possible.
eg:
```Go
const ONEWORD = 0
const ALL_UPPERCASE = 8
```


## Author

FireRain :=: 火雨
email: fr440305@gmail.com

## Bugs

- Sometimes, it will crash.
Maybe is has infinite loop somewhere. (apr, 1, 2017)
