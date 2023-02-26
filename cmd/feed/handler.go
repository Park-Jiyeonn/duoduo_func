package main

import (
	"context"
	feed "simple_tiktok/kitex_gen/feed"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// GetVideoList implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideoList(ctx context.Context, request *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideoListById implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideoListById(ctx context.Context, request *feed.GetByIDRequest) (resp *feed.FeedResponse, err error) {
	// TODO: Your code here...
	return
}
