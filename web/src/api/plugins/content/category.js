import service from '@/utils/request'

// 查询Category列表
export function listCategory(query) {
  return service({
    url: '/v1/content/category',
    method: 'get',
    params: query
  })
}

// 查询Category详细
export function getCategory(id) {
  return service({
    url: '/v1/content/category/' + id,
    method: 'get'
  })
}

// 新增Category
export function addCategory(data) {
  return service({
    url: '/v1/content/category',
    method: 'post',
    data: data
  })
}

// 修改Category
export function updateCategory(data) {
  return service({
    url: '/v1/content/category/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除Category
export function delCategory(data) {
  return service({
    url: '/v1/content/category',
    method: 'delete',
    data: data
  })
}

// 导出Category列表
export function exportCategory(query) {
  return service({
    url: '/v1/category/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

