package util

import (
	"os"
	"os/exec"
	"path/filepath"
)

func getVideoKeyframe(videoPath string) ([]byte, error) {
	// 创建一个临时目录
	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)
	// 使用 FFmpeg 提取关键帧
	imagePath := filepath.Join(tempDir, "cover.jpg")
	cmd := exec.Command("ffmpeg", "-ss", "00:00:00.5", "-i", videoPath, "-vframes", "1", "-q:v", "2", imagePath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 读取关键帧
	keyframe, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}
	return keyframe, nil
}

func Cover(videoPath, imagePath string) error {
	// 获取视频关键帧
	keyframe, err := getVideoKeyframe(videoPath)
	if err != nil {
		return err
	}
	// 将关键帧保存为图片
	err = os.WriteFile(imagePath, keyframe, 0644)
	//0644是一个文件权限掩码，用于表示文件的权限。它通常被用于Unix和类Unix系统中，例如Linux和Mac OS X。
	//
	//它由四个数字组成，每个数字代表文件所有者、文件所属组和其他用户的权限。每个数字的值分别是0-7，表示读、写和执行权限。
	//
	//在0644中，第一位数字0表示这是一个常规文件类型，后面三个数字表示文件权限。6表示文件所有者具有读和写权限（4+2），4表示文件所属组有读权限，其他用户也有读权限，但没有写权限和执行权限。
	//
	//因此，0644表示这个文件的权限为 -rw-r--r--，即文件所有者可以读写，文件所属组和其他用户可以读取，但不能写入或执行。
	if err != nil {
		return err
	}
	return nil
}
