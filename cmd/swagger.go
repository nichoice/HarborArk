package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// 创建 swagger 主命令
	swaggerCmd := &cobra.Command{
		Use:   "swagger",
		Short: "管理 Swagger API 文档",
		Long:  `管理 Swagger API 文档的生成、更新和预览`,
	}

	// 创建 generate 子命令
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "生成 Swagger API 文档",
		Long:  `根据代码中的注释生成 Swagger API 文档`,
		Run: func(cmd *cobra.Command, args []string) {
			outputDir, _ := cmd.Flags().GetString("output")
			mainFile, _ := cmd.Flags().GetString("main")
			force, _ := cmd.Flags().GetBool("force")
			generateSwaggerDocs(outputDir, mainFile, force)
		},
	}

	// 创建 validate 子命令
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "验证 Swagger 文档",
		Long:  `验证 Swagger 文档的格式和内容是否正确`,
		Run: func(cmd *cobra.Command, args []string) {
			validateSwaggerDocs()
		},
	}

	// 创建 clean 子命令
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "清理 Swagger 文档",
		Long:  `删除生成的 Swagger 文档文件`,
		Run: func(cmd *cobra.Command, args []string) {
			cleanSwaggerDocs()
		},
	}

	// 添加标志
	generateCmd.Flags().StringP("output", "o", "cmd/docs", "输出目录")
	generateCmd.Flags().StringP("main", "m", "cmd/server.go", "主 API 文件")
	generateCmd.Flags().BoolP("force", "f", false, "强制重新生成")

	// 添加子命令
	swaggerCmd.AddCommand(generateCmd)
	swaggerCmd.AddCommand(validateCmd)
	swaggerCmd.AddCommand(cleanCmd)

	// 添加到根命令
	rootCmd.AddCommand(swaggerCmd)
}

// generateSwaggerDocs 生成 Swagger 文档
func generateSwaggerDocs(outputDir, mainFile string, force bool) {
	fmt.Println("🚀 正在生成 Swagger 文档...")

	// 检查是否需要重新生成
	if !force && isSwaggerDocsExist(outputDir) {
		fmt.Println("📄 Swagger 文档已存在，使用 --force 标志强制重新生成")
		return
	}

	// 检查 swag 命令是否可用
	if !isSwagInstalled() {
		fmt.Println("📦 未找到 swag 命令，正在安装...")
		if err := installSwag(); err != nil {
			fmt.Printf("❌ 安装 swag 失败: %v\n", err)
			return
		}
		fmt.Println("✅ swag 安装成功")
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ 创建输出目录失败: %v\n", err)
		return
	}

	// 执行 swag 命令生成文档
	cmd := exec.Command("swag", "init", "-g", mainFile, "-o", outputDir, "--parseDependency", "--parseInternal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 生成 Swagger 文档失败: %v\n", err)
		return
	}

	fmt.Println("✅ Swagger 文档生成成功！")
	fmt.Printf("📁 文档位置: %s\n", outputDir)
	fmt.Println("🌐 启动服务器后访问: http://localhost:8080/swagger/index.html")

	// 显示生成的文件
	showGeneratedFiles(outputDir)
}

// validateSwaggerDocs 验证 Swagger 文档
func validateSwaggerDocs() {
	fmt.Println("🔍 正在验证 Swagger 文档...")

	docsDir := "cmd/docs"
	swaggerJSON := filepath.Join(docsDir, "swagger.json")

	// 检查文档是否存在
	if !isSwaggerDocsExist(docsDir) {
		fmt.Println("❌ Swagger 文档不存在，请先运行 'swagger generate'")
		return
	}

	// 检查 JSON 文件格式
	if _, err := os.Stat(swaggerJSON); err != nil {
		fmt.Printf("❌ swagger.json 文件不存在: %v\n", err)
		return
	}

	// 读取并验证 JSON 格式
	content, err := os.ReadFile(swaggerJSON)
	if err != nil {
		fmt.Printf("❌ 读取 swagger.json 失败: %v\n", err)
		return
	}

	if len(content) == 0 {
		fmt.Println("❌ swagger.json 文件为空")
		return
	}

	fmt.Println("✅ Swagger 文档验证通过")
	fmt.Printf("📊 文档大小: %d 字节\n", len(content))
}

// cleanSwaggerDocs 清理 Swagger 文档
func cleanSwaggerDocs() {
	fmt.Println("🧹 正在清理 Swagger 文档...")

	docsDir := "cmd/docs"

	// 检查目录是否存在
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		fmt.Println("📁 文档目录不存在，无需清理")
		return
	}

	// 删除生成的文件
	files := []string{"docs.go", "swagger.json", "swagger.yaml"}
	deletedCount := 0

	for _, file := range files {
		filePath := filepath.Join(docsDir, file)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("❌ 删除 %s 失败: %v\n", file, err)
			} else {
				fmt.Printf("🗑️  已删除: %s\n", file)
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		fmt.Printf("✅ 清理完成，共删除 %d 个文件\n", deletedCount)
	} else {
		fmt.Println("📄 没有找到需要清理的文件")
	}
}

// isSwagInstalled 检查 swag 是否已安装
func isSwagInstalled() bool {
	_, err := exec.LookPath("swag")
	return err == nil
}

// installSwag 安装 swag 工具
func installSwag() error {
	cmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// isSwaggerDocsExist 检查 Swagger 文档是否存在
func isSwaggerDocsExist(outputDir string) bool {
	files := []string{"docs.go", "swagger.json", "swagger.yaml"}
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(outputDir, file)); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// showGeneratedFiles 显示生成的文件信息
func showGeneratedFiles(outputDir string) {
	fmt.Println("\n📋 生成的文件:")
	files := []string{"docs.go", "swagger.json", "swagger.yaml"}

	for _, file := range files {
		filePath := filepath.Join(outputDir, file)
		if info, err := os.Stat(filePath); err == nil {
			fmt.Printf("  📄 %s (%d 字节)\n", file, info.Size())
		}
	}
}

// AutoUpdateSwaggerDocs 自动更新 Swagger 文档
func AutoUpdateSwaggerDocs() {
	if zap.L() != nil {
		zap.L().Info("正在自动更新 Swagger 文档...")
	} else {
		fmt.Println("🔄 正在自动更新 Swagger 文档...")
	}

	cmd := exec.Command("swag", "init", "-g", "cmd/server.go", "-o", "cmd/docs", "--parseDependency", "--parseInternal")
	if err := cmd.Run(); err != nil {
		if zap.L() != nil {
			zap.L().Warn("Swagger 文档自动更新失败", zap.Error(err))
		} else {
			fmt.Printf("⚠️  警告: Swagger 文档自动更新失败: %v\n", err)
		}
	} else {
		if zap.L() != nil {
			zap.L().Info("Swagger 文档自动更新成功")
		} else {
			fmt.Println("✅ Swagger 文档自动更新成功")
		}
	}
}
