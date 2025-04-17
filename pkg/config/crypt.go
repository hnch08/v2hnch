package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"v2hnch/pkg/logger"
)

const (
	// 加密密钥，建议使用环境变量或其他安全的方式存储
	encryptionKey = "your-secret-key-here-32-bytes-long!"
)

// 生成加密密钥
func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:32] // 返回32字节的密钥
}

// Encrypt 加密数据
func Encrypt(data []byte) ([]byte, error) {
	logger.Info("开始加密数据")

	// 创建cipher
	key := deriveKey(encryptionKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("创建cipher失败: %v", err)
		return nil, fmt.Errorf("创建cipher失败: %v", err)
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("创建GCM失败: %v", err)
		return nil, fmt.Errorf("创建GCM失败: %v", err)
	}

	// 创建nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error("生成nonce失败: %v", err)
		return nil, fmt.Errorf("生成nonce失败: %v", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	logger.Info("数据加密成功")
	return ciphertext, nil
}

// Decrypt 解密数据
func Decrypt(data []byte) ([]byte, error) {
	logger.Info("开始解密数据")

	// 创建cipher
	key := deriveKey(encryptionKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("创建cipher失败: %v", err)
		return nil, fmt.Errorf("创建cipher失败: %v", err)
	}

	// 创建GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("创建GCM失败: %v", err)
		return nil, fmt.Errorf("创建GCM失败: %v", err)
	}

	// 提取nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		logger.Error("加密数据长度不足")
		return nil, fmt.Errorf("加密数据长度不足")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error("解密数据失败: %v", err)
		return nil, fmt.Errorf("解密数据失败: %v", err)
	}

	logger.Info("数据解密成功")
	return plaintext, nil
}

// WriteEncrypted 将数据加密后写入文件
func WriteEncrypted(filename string, data []byte) error {
	logger.Info("开始加密写入文件")

	// 加密数据
	encrypted, err := Encrypt(data)
	if err != nil {
		return err
	}

	// 写入文件
	err = os.WriteFile(filename, encrypted, 0644)
	if err != nil {
		logger.Error("写入加密文件失败: %v", err)
		return fmt.Errorf("写入加密文件失败: %v", err)
	}

	logger.Info("加密文件写入成功")
	return nil
}

// ReadEncrypted 从文件读取并解密数据
func ReadEncrypted(filename string) ([]byte, error) {
	logger.Info("开始读取加密文件")

	// 读取文件
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("读取加密文件失败: %v", err)
		return nil, fmt.Errorf("读取加密文件失败: %v", err)
	}

	// 解密数据
	decrypted, err := Decrypt(data)
	if err != nil {
		return nil, err
	}

	logger.Info("加密文件读取成功")
	return decrypted, nil
}
