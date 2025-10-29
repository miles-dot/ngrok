package server

import (
        "database/sql"
        "fmt"
        "ngrok/log"
        "time"

        _ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接
var DB *sql.DB

// TokenInfo 存储token的详细信息
type TokenInfo struct {
        ID        int
        Name      string
        Phone     string
        Token     string
        Remark    string
        ExpiresAt time.Time
        CreatedAt time.Time
}

// InitDB 初始化数据库连接
func InitDB(dbConfig *DatabaseConfig) error {
        // 如果数据库功能未启用，直接返回
        if !dbConfig.Enabled {
                log.Info("数据库认证功能未启用")
                return nil
        }

        var err error
        DB, err = sql.Open(dbConfig.Driver, dbConfig.DSN)
        if err != nil {
                return fmt.Errorf("数据库连接失败: %v", err)
        }

        // 设置连接池参数
        DB.SetMaxOpenConns(dbConfig.MaxOpenConns)
        DB.SetMaxIdleConns(dbConfig.MaxIdleConns)
        DB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

        // 测试连接
        if err = DB.Ping(); err != nil {
                return fmt.Errorf("数据库连接测试失败: %v", err)
        }

        log.Info("数据库连接成功，已启用token认证")
        return nil
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenStr string) (*TokenInfo, error) {
        if DB == nil {
                return nil, fmt.Errorf("数据库未连接")
        }

        token := &TokenInfo{}
        query := `SELECT id, name, phone, token, remark, expires_at, created_at FROM tokens WHERE (phone = ? OR token = ?) AND expires_at > NOW()`

        err := DB.QueryRow(query, tokenStr, tokenStr).Scan(
                &token.ID,
                &token.Name,
                &token.Phone,
                &token.Token,
                &token.Remark,
                &token.ExpiresAt,
                &token.CreatedAt,
        )

        if err == sql.ErrNoRows {
                return nil, nil // token不存在或已过期
        } else if err != nil {
                return nil, fmt.Errorf("查询token失败: %v", err)
        }

        return token, nil
}