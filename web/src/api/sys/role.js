import service from '@/utils/request'

// 查询角色列表
export function listRole(query) {
  return service({
    url: '/v1/role',
    method: 'get',
    params: query
  })
}

// 查询角色详细
export function getRole(roleId) {
  return service({
    url: '/v1/role/' + roleId,
    method: 'get'
  })
}

// 新增角色
export function addRole(data) {
  return service({
    url: '/v1/role',
    method: 'post',
    data: data
  })
}

// 修改角色
export function updateRole(data, roleId) {
  return service({
    url: '/v1/role/' + roleId,
    method: 'put',
    data: data
  })
}

// 角色数据权限
export function dataScope(data) {
  return service({
    url: '/v1/roledatascope',
    method: 'put',
    data: data
  })
}

// 角色状态修改
export function changeRoleStatus(roleId, status) {
  const data = {
    roleId,
    status
  }
  return service({
    url: '/v1/role-status',
    method: 'put',
    data: data
  })
}

// 删除角色
export function delRole(roleId) {
  return service({
    url: '/v1/role',
    method: 'delete',
    data: roleId
  })
}

export function getListrole(id) {
  return service({
    url: '/v1/menu/role/' + id,
    method: 'get'
  })
}

export function getRoutes() {
  return service({
    url: '/v1/menurole',
    method: 'get'
  })
}

// export function getNames() {
//   return request({
//     url: '/v1/menuids',
//     method: 'get'
//   })
// }
