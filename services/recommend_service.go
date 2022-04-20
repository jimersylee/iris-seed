package services

var RecommendService = newRecommendService()

func newRecommendService() *recommendService {
	return &recommendService{}
}

type recommendService struct {
}

//
// Action
// @Description: 针对文章的操作,用来给推荐系统喂数据
// @receiver this
// @param actionName 操作类型,如:like,collect,comment,dislike,view,share
// @param articleId 文章id
// @param userId 用户id
// @param fingerprint 浏览器指纹
//
func (this *recommendService) Action(actionName string, articleId int, userId int64, fingerprint string) {

	return
}

//
// GetRecommendArticles
// @Description: 获取推荐文章
// @receiver this
// @param lastId 分页的最后一个id
// @param size  分页的大小
//
func (this *recommendService) GetRecommendArticles(lastId int, size int) {

}
