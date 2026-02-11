package shared

import (
	"fmt"
	"path/filepath"
	"time"
)

func CreateFileName(imageName string) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(imageName))
}
