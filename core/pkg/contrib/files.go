package contrib

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const FILE_FIELD = "file"

func parseFormFail(context *gin.Context) {
	context.JSON(http.StatusBadRequest, gin.H{
		"message": "Can not parse form",
	})
}

func getFileLists(filepath string) []os.FileInfo {
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		log.Print(err)
	}
	return files
}
