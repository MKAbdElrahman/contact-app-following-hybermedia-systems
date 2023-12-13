package handler

import (
	"app/view/blog"

	"github.com/labstack/echo/v4"
)

type BlogHandler struct {
	Blog *blog.Blog
}

func NewBlogHandler(blog *blog.Blog) *BlogHandler {

	return &BlogHandler{
		Blog: blog,
	}
}

func (h *BlogHandler) HandleGetIndexPage(ctx echo.Context) error {

	return h.Blog.RenderIndexPage(ctx, blog.IndexPageData{})
}

func (h *BlogHandler) HandleGetContactPage(ctx echo.Context) error {

	return h.Blog.RenderContactPage(ctx, blog.ContactPageData{})
}
