package example

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func TestCkInsert(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	if err := conn.PingContext(ctx); err != nil {
		log.Printf("%v", err)
	}
	fmt.Println("Connected to ClickHouse!")
}

func TestCreateDatabase(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS test_db")
	if err != nil {
		fmt.Printf("创建数据库失败: %v\n", err)
		return
	}
	fmt.Println("数据库创建成功")
}

func TestCreateTable(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "test_db",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_table (
		id UInt64,
		name String,
		created_at DateTime
	) ENGINE = MergeTree()
	ORDER BY id`)
	if err != nil {
		fmt.Printf("创建表失败: %v\n", err)
		return
	}
	fmt.Println("数据表创建成功")
}

func TestInsertData(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "test_db",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, "INSERT INTO test_table (id, name, created_at) VALUES (1, '测试数据', now())")
	if err != nil {
		fmt.Printf("插入数据失败: %v\n", err)
		return
	}
	fmt.Println("数据插入成功")
}

func TestUpdateData(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "test_db",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, "ALTER TABLE test_table UPDATE name = '更新数据' WHERE id = 1")
	if err != nil {
		fmt.Printf("更新数据失败: %v\n", err)
		return
	}
	fmt.Println("数据更新成功")
}

func TestDeleteData(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "test_db",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, "ALTER TABLE test_table DELETE WHERE id = 1")
	if err != nil {
		fmt.Printf("删除数据失败: %v\n", err)
		return
	}
	fmt.Println("数据删除成功")
}

func TestQueryData(t *testing.T) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"192.168.2.200:9000"},
		Auth: clickhouse.Auth{
			Database: "test_db",
			Username: "user1",
			Password: "user1user1",
		},
		DialTimeout: 10 * time.Second,
	})

	ctx := context.Background()
	rows, err := conn.QueryContext(ctx, "SELECT * FROM test_table")
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("查询结果:")
	for rows.Next() {
		var (
			id        uint64
			name      string
			createdAt time.Time
		)
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			fmt.Printf("解析行失败: %v\n", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s, CreatedAt: %s\n", id, name, createdAt)
	}
}
