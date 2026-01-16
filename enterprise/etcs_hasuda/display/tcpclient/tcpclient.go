/*
   Display TCPClient
*/

package tcpclient

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/*
TCPクライアント構造体
*/
type Client struct {
	addr string // アドレス

	Timeout time.Duration // タイムアウト

	// TCP connection
	mu   sync.Mutex // 排他制御用
	conn net.Conn   // 通信コンテキスト
}

/*
新しいクライアントを作成する
*/
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

/*
任意データ(1024バイト以下)受信
*/
func (c *Client) Read_data() ([]byte, error) {
	buf := make([]byte, 1024)
	rlen, err := c.conn.Read(buf)
	if err != nil {
		_ = c.conn.Close()
		return []byte{}, err
	}

	// 受信文字列とエラーを返す
	return buf[:rlen], nil
}

/*
任意データ送信
*/
func (c *Client) Send_data(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.connect(); err != nil {
		fmt.Println("ConnectError")
		return err
	}

	if c.Timeout > 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
			_ = c.conn.Close()
			fmt.Println("TimeOut")
			return err
		}
	}

	// データ送信
	_, err := c.conn.Write(data)
	if err != nil {
		_ = c.conn.Close()
		fmt.Println(err)
		fmt.Println("CommandSendError")
		return err
	}
	fmt.Printf("Data sent to client %s: %s\n", c.conn.RemoteAddr(), string(data))

	return nil
}

/*
サーバーへの通信切断
*/
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.close()
}

/*
サーバーへの接続
*/
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connect()
}

/*
接続(コネクション)がnilの場合は、接続を試みる。
すでに接続されている場合は、そのコネクションを維持して使い回す。
*/
func (c *Client) connect() error {
	if c.conn == nil {
		conn, err := net.DialTimeout("tcp", c.addr, c.Timeout)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

// closeはコネクションがnilでない場合にコネクションをCloseします。
func (c *Client) close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}
	return err
}

// connを返す
func (c *Client) GetConn() net.Conn {
	return c.conn
}
