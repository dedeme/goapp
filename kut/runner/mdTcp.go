// Copyright 14-Mar-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"bufio"
	"fmt"
	"github.com/dedeme/kut/builtin/bfail"
	"github.com/dedeme/kut/builtin/bfunction"
	"github.com/dedeme/kut/expression"
	"io"
	"net"
	"strconv"
	"time"
)

type tcpConnT struct {
	conn net.Conn
	scan *bufio.Scanner
}

// \<tpcServer>, i -> [<tcpConn>, s]
func tcpAccept(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case net.Listener:
		switch tm := (args[1].Value).(type) {
		case int64:
			if tm < 0 {
				tm = 0
			}
			conn, er := v.Accept()
			conEx := expression.MkFinal(conn)
			errs := ""
			if er != nil {
				errs = "tcp.Accept ERROR:\n" + er.Error()
				conEx = expression.MkEmpty()
			} else if tm > 0 {
				err = conn.SetDeadline(time.Now().Add(time.Duration(tm) * time.Millisecond))
			}
			errEx := expression.MkFinal(errs)
			ex = expression.MkFinal([]*expression.T{conEx, errEx})
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<tcpServer>")
	}
	return
}

// \<tpcConn> -> ()
func tcpCloseConnection(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case net.Conn:
		err = v.Close()
	default:
		err = bfail.Type(args[0], "<tcpConnection>")
	}
	return
}

// \<tpcServer> -> ()
func tcpCloseServer(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case net.Listener:
		err = v.Close()
	default:
		err = bfail.Type(args[0], "<tcpServer>")
	}
	return
}

// \s, i -> [<tcpConn>, s]
func tcpDial(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case string:
		switch tm := (args[1].Value).(type) {
		case int64:
			if tm < 0 {
				tm = 0
			}
			conn, er := net.Dial("tcp", v)
			conEx := expression.MkFinal(conn)
			errs := ""
			if er != nil {
				errs = "tcp.Accept ERROR:\n" + er.Error()
				conEx = expression.MkEmpty()
			} else if tm > 0 {
				err = conn.SetDeadline(time.Now().Add(time.Duration(tm) * time.Millisecond))
			}
			errEx := expression.MkFinal(errs)
			ex = expression.MkFinal([]*expression.T{conEx, errEx})
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "string")
	}
	return
}

// \<tcpConn>, i -> [s, s]
func tcpRead(args []*expression.T) (ex *expression.T, err error) {
	switch cn := (args[0].Value).(type) {
	case net.Conn:
		switch lim := (args[1].Value).(type) {
		case int64:
			if lim < 1 {
				err = bfail.Mk("Connection limit less than 1")
				return
			}
			bs := make([]byte, lim+1)
			var n int
			n, err0 := cn.Read(bs)
			if err0 != nil {
				if err0 == io.EOF {
					ex = expression.MkFinal([]*expression.T{
						expression.MkFinal(""),
						expression.MkFinal(""),
					})
				} else {
					ex = expression.MkFinal([]*expression.T{
						expression.MkFinal(""),
						expression.MkFinal(err0.Error()),
					})
				}
				return
			}
			n2 := int64(n)
			if n2 > lim {
				ex = expression.MkFinal([]*expression.T{
					expression.MkFinal(""),
					expression.MkFinal(fmt.Sprintf("Bytes read out of limit (%v)", lim)),
				})
				return
			}
			bs2 := make([]byte, n)
			for i := 0; i < n; i++ {
				bs2[i] = bs[i]
			}
			ex = expression.MkFinal([]*expression.T{
				expression.MkFinal(string(bs2)),
				expression.MkFinal(""),
			})
		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<tcpConnection>")
	}
	return
}

// \<tcpConn>, i -> [s, <bytes>]
func tcpReadBin(args []*expression.T) (ex *expression.T, err error) {
	switch cn := (args[0].Value).(type) {
	case net.Conn:
		switch lim := (args[1].Value).(type) {
		case int64:
			if lim < 1 {
				err = bfail.Mk("Connection limit less than 1")
				return
			}
			bs := make([]byte, lim+1)
			n, err0 := cn.Read(bs)
			if err0 != nil {
				if err0 == io.EOF {
					ex = expression.MkFinal([]*expression.T{
						expression.MkFinal(""),
						expression.MkFinal(""),
					})
				} else {
					ex = expression.MkFinal([]*expression.T{
						expression.MkFinal(""),
						expression.MkFinal(err0.Error()),
					})
				}
				return
			}
			n2 := int64(n)
			if n2 > lim {
				ex = expression.MkFinal([]*expression.T{
					expression.MkFinal(""),
					expression.MkFinal(fmt.Sprintf("Bytes read out of limit (%v)", lim)),
				})
				return
			}
			bs2 := make([]byte, n)
			for i := 0; i < n; i++ {
				bs2[i] = bs[i]
			}
			ex = expression.MkFinal([]*expression.T{
				expression.MkFinal(bs2),
				expression.MkFinal(""),
			})

		default:
			err = bfail.Type(args[1], "int")
		}
	default:
		err = bfail.Type(args[0], "<tcpConnection>")
	}
	return
}

// \i -> <tpcServer>
func tcpServer(args []*expression.T) (ex *expression.T, err error) {
	switch v := (args[0].Value).(type) {
	case int64:
		port := ":" + strconv.FormatInt(v, 10)
		var sv net.Listener
		sv, err = net.Listen("tcp", port)
		if err == nil {
			ex = expression.MkFinal(sv)
		}
	default:
		err = bfail.Type(args[0], "int")
	}
	return
}

// \<tcpConn>, s -> [s]|[]
func tcpWrite(args []*expression.T) (ex *expression.T, err error) {
	switch cn := (args[0].Value).(type) {
	case net.Conn:
		switch s := (args[1].Value).(type) {
		case string:
			_, err0 := fmt.Fprintf(cn, s)
			if err0 != nil {
				ex = expression.MkFinal([]*expression.T{
					expression.MkFinal(err0.Error()),
				})
			} else {
				ex = expression.MkFinal([]*expression.T{})
			}
			return
		default:
			err = bfail.Type(args[1], "string")
		}
	default:
		err = bfail.Type(args[0], "<tcpConnection>")
	}
	return
}

// \<tcpConn>, <bytes> -> [s]|[]
func tcpWriteBin(args []*expression.T) (ex *expression.T, err error) {
	switch cn := (args[0].Value).(type) {
	case net.Conn:
		switch bs := (args[1].Value).(type) {
		case []byte:
			_, err0 := cn.Write(bs)
			if err0 != nil {
				ex = expression.MkFinal([]*expression.T{
					expression.MkFinal(err0.Error()),
				})
			} else {
				ex = expression.MkFinal([]*expression.T{})
			}
		default:
			err = bfail.Type(args[1], "<byte>")
		}
	default:
		err = bfail.Type(args[0], "<tcpConnection>")
	}
	return
}

func tcpGet(fname string) (fn *bfunction.T, ok bool) {
	ok = true
	switch fname {
	case "accept":
		fn = bfunction.New(2, tcpAccept)
	case "closeConnection":
		fn = bfunction.New(1, tcpCloseConnection)
	case "closeServer":
		fn = bfunction.New(1, tcpCloseServer)
	case "dial":
		fn = bfunction.New(2, tcpDial)
	case "read":
		fn = bfunction.New(2, tcpRead)
	case "readBin":
		fn = bfunction.New(2, tcpRead)
	case "server":
		fn = bfunction.New(1, tcpServer)
	case "write":
		fn = bfunction.New(2, tcpWrite)
	case "writeBin":
		fn = bfunction.New(2, tcpWriteBin)
	default:
		ok = false
	}

	return
}
