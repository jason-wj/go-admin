import service from '@/utils/request'

// 查询App列表
export function listApp(query) {
  return service({
    url: '/v1/filemgr/app',
    method: 'get',
    params: query
  })
}

// 查询App详细
export function getApp(id) {
  return service({
    url: '/v1/filemgr/app/' + id,
    method: 'get'
  })
}

// 新增App
export function addApp(data) {
  return service({
    url: '/v1/filemgr/app',
    method: 'post',
    data: data
  })
}

// 修改App
export function updateApp(data) {
  return service({
    url: '/v1/filemgr/app/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除App
export function delApp(data) {
  return service({
    url: '/v1/filemgr/app',
    method: 'delete',
    data: data
  })
}

// 导出App列表
export function exportApp(query) {
  return service({
    url: '/v1/filemgr/app/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

