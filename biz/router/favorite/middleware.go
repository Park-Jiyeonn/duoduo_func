// Code generated by hertz generator.

package Favorite

import (
	"github.com/cloudwego/hertz/pkg/app"
	"simple_tiktok/biz/mw"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _douyinMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoriteMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _actionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _like_ctionMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{mw.JwtMiddleware()}
}

func _listMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getlikelistMw() []app.HandlerFunc {
	// your code...
	return nil
}