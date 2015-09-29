package context

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"math"
	"strconv"
)

func addPaginationLinks(halResource HalResource, baseURL string, currentPage, recordsPerPage, totalRecordCount int) {
	var firstPage = 1

	var lastPage = int(math.Ceil(float64(totalRecordCount) / float64(recordsPerPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	var prevPage = currentPage - 1
	var nextPage = currentPage + 1

	halResource.AddLink(relSelf, generatePaginationURL(baseURL, currentPage))

	if firstPage != lastPage {
		halResource.AddLink(relFirst, generatePaginationURL(baseURL, firstPage))
		halResource.AddLink(relLast, generatePaginationURL(baseURL, lastPage))
	}
	if prevPage >= firstPage && prevPage < lastPage {
		halResource.AddLink(relPrev, generatePaginationURL(baseURL, prevPage))
	}
	if nextPage <= lastPage {
		halResource.AddLink(relNext, generatePaginationURL(baseURL, nextPage))
	}
}

func generatePaginationURL(baseURL string, page int) string {
	if page > 1 {
		return fmt.Sprintf("%s?page=%d", baseURL, page)
	}
	return baseURL
}

func getStartRecord(currentPage, recordsPerPage int) int {
	return (currentPage - 1) * recordsPerPage
}

func getCurrentPage(c *gin.Context) int {
	pageString := c.Request.URL.Query().Get("page")

	var currentPage = 1

	pageValue, err := strconv.ParseInt(pageString, 10, 32)
	if err == nil {
		currentPage = int(pageValue)
	}
	if currentPage < 1 {
		currentPage = 1
	}
	return currentPage
}