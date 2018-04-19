package server

// type conn struct {
// 	net.Conn
// }

// func (c *conn) write(p []byte) (n int, err error) {
// 	c.updateDeadline()
// 	n, err = c.Conn.Write(p)
// 	return
// }

// func (c *conn) read(b []byte) (n int, err error) {
// 	c.updateDeadline()
// 	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
// 	n, err = r.Read(b)
// 	return
// }

// func (c *conn) close() (err error) {
// 	err = c.Conn.Close()
// 	return
// }

// func (c *conn) updateDeadline() {
// 	idleDeadline := time.Now().Add(c.IdleTimeout)
// 	c.Conn.SetDeadline(idleDeadline)
// }
