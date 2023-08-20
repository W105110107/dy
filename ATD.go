package main

import (
	m "dy/controller"
	"fmt"
)

// 添加测试数据
func Add_Test_Data() {
	//创建数据库表
	err := m.DB.AutoMigrate(&m.User{}, &m.Video{}, &m.Comment{}, &m.Token{}, &m.PeopleMessage{})
	if err != nil {
		fmt.Println(err)
	}
	// 添加测试数据
	var userList []m.User
	//m.DB.Find(&userList).Delete(&userList)
	userList = []m.User{
		{Avatar: m.URL + "测试1.jpg", BackgroundImage: m.URL + "测试1H.jpg", Signature: "测试人员1的个人简介", TotalFavorite: 0, ID: 1, Name: "测试人员1", FollowList: []int64{}, FollowerList: []int64{}, FavoriteList: []int64{}, Videos: []m.Video{}},
		{Avatar: m.URL + "测试2.jpg", BackgroundImage: m.URL + "测试2H.jpg", Signature: "测试人员2的个人简介", TotalFavorite: 0, ID: 2, Name: "测试人员2", FollowList: []int64{}, FollowerList: []int64{}, FavoriteList: []int64{}, Videos: []m.Video{}},
		{Avatar: m.URL + "测试3.jpg", BackgroundImage: m.URL + "测试3H.jpg", Signature: "测试人员3的个人简介", TotalFavorite: 0, ID: 3, Name: "测试人员3", FollowList: []int64{}, FollowerList: []int64{}, FavoriteList: []int64{}, Videos: []m.Video{}},
	}
	m.DB.Create(&userList)
	//"2006-01-02 15:04:05"
	var VideoList []m.Video
	//m.DB.Find(&VideoList).Delete(&VideoList)
	VideoList = []m.Video{
		{
			CoverURL:     "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
			ID:           1,
			PlayURL:      "https://www.w3schools.com/html/movie.mp4",
			Title:        "熊",
			CreateTime:   "2023-08-01 12:00:00",
			UserID:       1,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/664bc4e86cfae46338056e7ec016555e.jpg",
			ID:           2,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/664bc4e86cfae46338056e7ec016555e.mp4",
			Title:        "生活",
			CreateTime:   "2023-08-01 13:00:00",
			UserID:       1,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/b31512aa2916f3ab484f8a4be569a0fa.jpg",
			ID:           3,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/b31512aa2916f3ab484f8a4be569a0fa.mp4",
			Title:        "鞭炮",
			CreateTime:   "2023-08-02 12:00:00",
			UserID:       1,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/a1da2b17a124aa88d0134d75d3983b5f.jpg",
			ID:           4,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/a1da2b17a124aa88d0134d75d3983b5f.mp4",
			Title:        "沙滩",
			CreateTime:   "2023-08-02 13:00:00",
			UserID:       1,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/dad18708c59804da1f3abb996cb56770.jpg",
			ID:           5,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/dad18708c59804da1f3abb996cb56770.mp4",
			Title:        "科技视频",
			CreateTime:   "2023-08-03 12:00:00",
			UserID:       1,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675682358.jpg",
			ID:           6,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675682358.mp4",
			Title:        "听海",
			CreateTime:   "2023-08-03 13:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/bea309f6840bee5d95c233616b3f1f34.jpg",
			ID:           7,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/bea309f6840bee5d95c233616b3f1f34.mp4",
			Title:        "手柄",
			CreateTime:   "2023-08-04 12:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/dbd19a6cba6bcf02027613b4caefdce8.jpg",
			ID:           8,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/dbd19a6cba6bcf02027613b4caefdce8.mp4",
			Title:        "行贿",
			CreateTime:   "2023-08-04 13:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/b29672e25430abc8eb31daecda52b8cb.jpg",
			ID:           9,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/b29672e25430abc8eb31daecda52b8cb.mp4",
			Title:        "高铁",
			CreateTime:   "2023-08-05 12:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1671fadf2ac23fed56996c3dc935ce92.jpg",
			ID:           10,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1671fadf2ac23fed56996c3dc935ce92.mp4",
			Title:        "火烈鸟",
			CreateTime:   "2023-08-05 13:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/V30203-123620.jpg",
			ID:           11,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/V30203-123620.mp4",
			Title:        "阳台",
			CreateTime:   "2023-08-06 12:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/安欣霸凌高启强-哔哩哔哩_302651035.jpg",
			ID:           12,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/%E5%AE%89%E6%AC%A3%E9%9C%B8%E5%87%8C%E9%AB%98%E5%90%AF%E5%BC%BA-%E5%93%94%E5%93%A9%E5%93%94%E5%93%A9_302651035.mp4",
			Title:        "安欣霸凌",
			CreateTime:   "2023-08-06 13:00:00",
			UserID:       2,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675658873.jpg",
			ID:           13,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675658873.mp4",
			Title:        "玻利维亚",
			CreateTime:   "2023-08-07 12:00:00",
			UserID:       3,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
		{
			CoverURL:     "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675685888.jpg",
			ID:           14,
			PlayURL:      "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/1675685888.mp4",
			Title:        "电梯",
			CreateTime:   "2023-08-07 13:00:00",
			UserID:       3,
			Comments:     []m.Comment{},
			FavoriteList: []int64{},
		},
	}
	m.DB.Create(&VideoList)

	var tokenList []m.Token
	//m.DB.Find(&tokenList).Delete(&tokenList)
	tokenList = []m.Token{
		{Token: "测试人员1-123456a", ID: 1},
		{Token: "测试人员2-123456b", ID: 2},
		{Token: "测试人员3-123456c", ID: 3},
	}
	m.DB.Create(&tokenList)
}
