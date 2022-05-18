import service from '@/utils/request'

// 查询Post列表
export function listPost(query) {
  return service({
    url: '/v1/sys/post',
    method: 'get',
    params: query
  })
}

// 查询Post详细
export function getPost(postId) {
  return service({
    url: '/v1/sys/post/' + postId,
    method: 'get'
  })
}

// 新增Post
export function addPost(data) {
  return service({
    url: '/v1/sys/post',
    method: 'post',
    data: data
  })
}

// 修改Post
export function updatePost(data) {
  return service({
    url: '/v1/sys/post/' + data.postId,
    method: 'put',
    data: data
  })
}

// 删除Post
export function delPost(data) {
  return service({
    url: '/v1/sys/post',
    method: 'delete',
    data: data
  })
}

// 导出Post列表
export function exportPost(query) {
  return service({
    url: '/v1/sys/post/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

