// Copyright 31-May-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
	"github.com/dedeme/ktlib/tcp"
	"github.com/dedeme/ktlib/thread"
)

// Process to run in server
func process(conn tcp.ConnT) bool {
	tx, _ := tcp.Read(conn, 10000)
	if tx == "end" {
		tcp.Write(conn, "Closing server")
		tcp.CloseConnection(conn)
		return true
	}
	tcp.Write(conn, "Send from server: "+tx)
	tcp.CloseConnection(conn)
	return false
}

func main() {

	// Launch server.
	th1 := thread.Start(func() {
		sv := tcp.Server(23344)

		for {
			conn, err := tcp.Accept(sv, 0)
			if err != nil {
				panic(err)
			}
			if process(conn) {
				break
			}
		}

		tcp.CloseServer(sv)
	})

	// Three connection from client.
	for i := 0; i < 3; i++ {
		conn, err := tcp.Dial("localhost:23344", 0)
		if err != nil {
			panic(err)
		}
		sys.Println(str.Fmt("Sending 'abc%d'", i))
		tcp.Write(conn, str.Fmt("abc%d", i))
    bs, _ := tcp.Read(conn, 10000)
		sys.Println(bs)
		tcp.CloseConnection(conn)
	}

	// Ending server
	conn, err := tcp.Dial("localhost:23344", 0)
	if err != nil {
		panic(err)
	}
	tcp.Write(conn, "end")
  bs, _ := tcp.Read(conn, 10000)
	sys.Println(bs)
	tcp.CloseConnection(conn)

	// Wait until server is ended.
	thread.Join(th1)
}
