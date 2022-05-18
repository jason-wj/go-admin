import service from '@/utils/request'

// 查询Api列表
export function listApi(query) {
  return service({
    url: '/v1/sys-api',
    method: 'get',
    params: query
  })
}

// 查询Api详细
export function getApi(id) {
  return service({
    url: '/v1/sys-api/' + id,
    method: 'get'
  })
}

// 修改Api
export function updateApi(data) {
  return service({
    url: '/v1/sys-api/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除Api
export function delApi(data) {
  return service({
    url: '/v1/sys-api',
    method: 'delete',
    data: data
  })
}

// 导出Api列表
export function exportApi(query) {
  return service({
    url: '/v1/sys-api/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

