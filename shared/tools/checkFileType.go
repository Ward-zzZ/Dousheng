package tools

import (
	"fmt"
	"github.com/h2non/filetype"
)

// checkFileType 函数会判断给定的文件是否是视频文件并返回文件类型
func checkFileType(fileData []byte) (bool, string, error) {
	kind, err := filetype.Match(fileData)
	if err != nil {
		return false, "", fmt.Errorf("failed to match file type: %v", err)
	}

	// 检查 MIME 类型是否为 video
	isVideo := kind.MIME.Type == "video"
	return isVideo, kind.Extension, nil
}
