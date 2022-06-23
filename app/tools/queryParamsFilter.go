package tools

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

const (
	NameLimitParam = "limit"
	NameDescParam  = "desc"
	NameSinceParam = "since"
	NameSortParam = "sort"
	NameRelatedParam = "related"
)

const (
	LimitParamDefault = 100
	SinceParamDefault = ""
	SortParamDefault  = "asc"
	SortParamTrue     = "desc"
	SortParamTree     = "tree"
	SortParamParentTree = "parent_tree"
	SortParamFlatDefault = "flat"
)

type FilterThread struct {
	Limit int
	Sort  string
	Since string
}

type FilterPosts struct {
	Limit int
	Sort string
	Since string
	Desc string
}

type FilterUser struct {
	Limit int
	Since string
	Desc string
}

type FilterOnePost struct {
	User bool
	Forum bool
	Thread bool
}

func ParseQueryFilterThread(ctx echo.Context) FilterThread {
	var result FilterThread
	queryParam := ctx.QueryParams()

	limit := queryParam.Get(NameLimitParam)
	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			result.Limit = 100
		} else {
			result.Limit = int(limitInt)
		}
	} else {
		result.Limit = LimitParamDefault
	}

	sort := queryParam.Get(NameDescParam)
	if sort == "true" {
		result.Sort = SortParamTrue
	} else {
		result.Sort = SortParamDefault
	}

	since := queryParam.Get(NameSinceParam)
	result.Since = since

	return result
}

func ParseQueryFilterPost (ctx echo.Context) FilterPosts{
	var result FilterPosts
	queryParam := ctx.QueryParams()

	limit := queryParam.Get(NameLimitParam)
	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			result.Limit = 100
		} else {
			result.Limit = int(limitInt)
		}
	} else {
		result.Limit = LimitParamDefault
	}

	switch queryParam.Get(NameSortParam) {
	case SortParamTree:
		result.Sort = SortParamTree
	case SortParamParentTree:
		result.Sort = SortParamParentTree
	default:
		result.Sort = SortParamFlatDefault
	}

	sort := queryParam.Get(NameDescParam)
	if sort == "true" {
		result.Desc = SortParamTrue
	} else {
		result.Desc = SortParamDefault
	}

	since := queryParam.Get(NameSinceParam)
	result.Since = since

	return result
}

func ParseQueryFilterUser (ctx echo.Context) FilterUser{
	var result FilterUser
	queryParam := ctx.QueryParams()

	limit := queryParam.Get(NameLimitParam)
	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			result.Limit = 100
		} else {
			result.Limit = int(limitInt)
		}
	} else {
		result.Limit = LimitParamDefault
	}

	sort := queryParam.Get(NameDescParam)
	if sort == "true" {
		result.Desc = SortParamTrue
	} else {
		result.Desc = SortParamDefault
	}

	since := queryParam.Get(NameSinceParam)
	result.Since = since

	return result
}

func ParseQueryFilterOnePost (ctx echo.Context) FilterOnePost {
	var result FilterOnePost
	queryParam := ctx.QueryParams()

	related := queryParam.Get(NameRelatedParam)
	if related != "" {
		if strings.Contains(related, "user") {
			result.User = true
		}
		if strings.Contains(related, "thread") {
			result.Thread = true
		}
		if strings.Contains(related, "forum") {
			result.Forum = true
		}
	}

	return result
}
