package service

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/pkg/constant"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FilesService struct {
	ctx *gin.Context
}

var (
	maxSize int64 = 20 << 20 // 20MB
	// 图片类型
	imageExtRange = []string{"jpg", "jpeg", "png", "gif", "bmp", "tif", "webp", "psd", "ai"}
	// 媒体类型
	mediaExtRange = []string{"mp4", "wmv", "rmvb", "3gp", "mov", "m4v", "avi", "mkv", "flv", "mpg", "mp3", "wma", "aac", "flac", "ape"}
	// 压缩文件
	compressExtRange = []string{"zip", "7z", "rar", "bz2", "gz", "tar"}
	// 生产力类型文件
	productivityExtRange = []string{"xmind", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "rp", "pdf"}
	// 其它类型
	otherExtRange = []string{"txt", "json"}
)

func NewFilesService(ctx *gin.Context) *FilesService {
	return &FilesService{
		ctx: ctx, // 传递上下文
	}
}

// UploadFile 上传文件
func (receiver *FilesService) UploadFile(file *multipart.FileHeader, fileType string) (*dto.FileVo, error) {
	if fileType == "" {
		fileType = "all"
	}
	// 文件大小限制
	if file.Size > maxSize {
		return nil, exception.NewException(response.FilesLimitExceeded)
	}
	// 获取文件后缀
	extName := path.Ext(file.Filename)

	// 检查文件后缀是否合规
	if !slice.Contain(receiver.SuffixSelect(fileType), strings.ToLower(strings.TrimPrefix(extName, "."))) {
		return nil, exception.NewException(response.FilesSuffixError)
	}

	// 创建上传目录
	savePath := filepath.Join("./uploads", time.Now().Format(constant.DateNoSeparationFormat))
	if !fileutil.IsExist(savePath) {
		if err := os.MkdirAll(savePath, 0666); err != nil {
			return nil, err
		}
	}

	// 生成文件名
	filename := cryptor.Md5String(file.Filename + strconv.FormatInt(file.Size, 10) + strconv.FormatInt(time.Now().UnixNano(), 10))
	saveFile := filepath.Join(savePath, filename+extName)
	if err := receiver.SaveUploadedFile(file, saveFile); err != nil {
		return nil, err
	}

	return &dto.FileVo{
		Name: filename + extName, // 新文件名
		Url:  "/" + saveFile,     // 文件路径，返回时在前面加上 /
		Tag:  file.Filename,      // 一般是文件的原始名称
		Ext:  extName,            // 文件后缀
		Size: file.Size,          // 文件大小
	}, nil
}

func (receiver *FilesService) SuffixSelect(typeName string) []string {
	switch typeName {
	case "image":
		return imageExtRange
	case "media":
		return mediaExtRange
	case "compress":
		return compressExtRange
	case "productivity":
		return productivityExtRange
	case "other":
		return otherExtRange
	case "all":
		tmp := append(imageExtRange, mediaExtRange...)
		tmp = append(tmp, compressExtRange...)
		tmp = append(tmp, productivityExtRange...)
		return append(tmp, otherExtRange...)
	default:
		return make([]string, 0)
	}
}

// SaveUploadedFile uploads the form file to specific dst.
// Copy
func (receiver *FilesService) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, src)
	return err
}
