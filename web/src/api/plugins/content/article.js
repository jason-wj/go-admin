import service from '@/utils/request'

// 查询Article列表
export function listArticle(query) {
  return service({
    url: '/v1/content/article',
    method: 'get',
    params: query
  })
}

// 查询Article详细
export function getArticle(id) {
  return service({
    url: '/v1/content/article/' + id,
    method: 'get'
  })
}

// 新增Article
export function addArticle(data) {
  return service({
    url: '/v1/content/article',
    method: 'post',
    data: data
  })
}

// 修改Article
export function updateArticle(data) {
  return service({
    url: '/v1/content/article/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除Article
export function delArticle(data) {
  return service({
    url: '/v1/content/article',
    method: 'delete',
    data: data
  })
}

// 导出Article列表
export function exportArticle(query) {
  return service({
    url: '/v1/article/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

