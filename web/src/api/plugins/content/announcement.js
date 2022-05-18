import service from '@/utils/request'

// 查询Announcement列表
export function listAnnouncement(query) {
  return service({
    url: '/v1/content/announcement',
    method: 'get',
    params: query
  })
}

// 查询Announcement详细
export function getAnnouncement(id) {
  return service({
    url: '/v1/content/announcement/' + id,
    method: 'get'
  })
}

// 新增Announcement
export function addAnnouncement(data) {
  return service({
    url: '/v1/content/announcement',
    method: 'post',
    data: data
  })
}

// 修改Announcement
export function updateAnnouncement(data) {
  return service({
    url: '/v1/content/announcement/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除Announcement
export function delAnnouncement(data) {
  return service({
    url: '/v1/content/announcement',
    method: 'delete',
    data: data
  })
}

// 导出Announcement列表
export function exportAnnouncement(query) {
  return service({
    url: '/v1/content/announcement/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

