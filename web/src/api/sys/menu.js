import service from '@/utils/request'

export const asyncMenu = () => {
  return service({
    url: '/v1/menurole',
    method: 'get'
  })
}

// 查询菜单列表
export function listMenu(query) {
  return service({
    url: '/v1/menu',
    method: 'get',
    params: query
  })
}

// 查询菜单详细
export function getMenu(menuId) {
  return service({
    url: '/v1/menu/' + menuId,
    method: 'get'
  })
}

// 新增菜单
export function addMenu(data) {
  return service({
    url: '/v1/menu',
    method: 'post',
    data: data
  })
}

// 修改菜单
export function updateMenu(data, id) {
  return service({
    url: '/v1/menu/' + id,
    method: 'put',
    data: data
  })
}

// 删除菜单
export function delMenu(data) {
  return service({
    url: '/v1/menu',
    method: 'delete',
    data: data
  })
}

// 新
// 根据角色ID查询菜单下拉树结构
export function roleMenuTreeselect(roleId) {
  return service({
    url: '/v1/sys/roleMenuTreeselect/' + roleId,
    method: 'get'
  })
}
