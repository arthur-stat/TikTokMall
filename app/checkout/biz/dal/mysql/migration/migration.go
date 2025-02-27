package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Migration 表示一个数据库迁移
type Migration struct {
	ID      int
	Name    string
	Content string
}

// RunMigrations 运行所有迁移脚本
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// 创建迁移记录表 (如果不存在)
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("创建迁移表失败: %w", err)
	}

	// 获取所有已执行的迁移
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("获取已应用的迁移失败: %w", err)
	}

	// 加载所有迁移文件
	migrations, err := loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("加载迁移文件失败: %w", err)
	}

	// 执行尚未应用的迁移
	for _, migration := range migrations {
		if _, ok := applied[migration.ID]; !ok {
			// 开始事务
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("开始事务失败: %w", err)
			}

			// 执行迁移
			if _, err := tx.Exec(migration.Content); err != nil {
				tx.Rollback()
				return fmt.Errorf("执行迁移 %s 失败: %w", migration.Name, err)
			}

			// 记录迁移已执行
			if _, err := tx.Exec("INSERT INTO migrations (id, name) VALUES (?, ?)", migration.ID, migration.Name); err != nil {
				tx.Rollback()
				return fmt.Errorf("记录迁移失败: %w", err)
			}

			// 提交事务
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("提交事务失败: %w", err)
			}

			fmt.Printf("已应用迁移: %s\n", migration.Name)
		}
	}

	return nil
}

// 创建迁移记录表
func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		)
	`)
	return err
}

// 获取所有已执行的迁移
func getAppliedMigrations(db *sql.DB) (map[int]struct{}, error) {
	rows, err := db.Query("SELECT id FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]struct{})
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		applied[id] = struct{}{}
	}

	return applied, rows.Err()
}

// 加载所有迁移文件
func loadMigrations(dir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// 解析文件名格式 (例如 001_init_schema.sql)
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}

		id := 0
		fmt.Sscanf(parts[0], "%d", &id)
		if id == 0 {
			continue
		}

		// 读取迁移文件内容
		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			ID:      id,
			Name:    parts[1][:len(parts[1])-4], // 移除 .sql 后缀
			Content: string(content),
		})
	}

	return migrations, nil
}
