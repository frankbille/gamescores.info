package service

import (
	"fmt"
	gin "github.com/gamescores/gin"
	"math"
	"strconv"
	"src/domain"
)

func addPaginationLinks(halResource domain.HalResource, baseURL string, currentPage, recordsPerPage, totalRecordCount int) {
	var firstPage = 1

	var lastPage = int(math.Ceil(float64(totalRecordCount) / float64(recordsPerPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	var prevPage = currentPage - 1
	var nextPage = currentPage + 1

	halResource.AddLink(domain.RelSelf, generatePaginationURL(baseURL, currentPage))

	if firstPage != lastPage {
		halResource.AddLink(domain.RelFirst, generatePaginationURL(baseURL, firstPage))
		halResource.AddLink(domain.RelLast, generatePaginationURL(baseURL, lastPage))
	}
	if prevPage >= firstPage && prevPage < lastPage {
		halResource.AddLink(domain.RelPrev, generatePaginationURL(baseURL, prevPage))
	}
	if nextPage <= lastPage {
		halResource.AddLink(domain.RelNext, generatePaginationURL(baseURL, nextPage))
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
