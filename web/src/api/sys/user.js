import service from '@/utils/request'
// @Summary 用户登录
export const login = (data) => {
  return service({
    url: '/v1/login',
    method: 'post',
    data: data
  })
}

// @Summary 获取验证码
export const captcha = (data) => {
  return service({
    url: '/v1/captcha',
    method: 'get',
    data: data
  })
}

// @Summary 用户注册
export const register = (data) => {
  return service({
    url: '/user/admin_register',
    method: 'post',
    data: data
  })
}

// @Summary 修改密码
export const changePassword = (data) => {
  return service({
    url: '/v1/user/pwd/set',
    method: 'put',
    data: data
  })
}

// 查询用户列表
export function listUser(query) {
  return service({
    url: '/v1/sys-user',
    method: 'get',
    params: query
  })
}

// 查询用户详细
export function getUser(userId) {
  return service({
    url: '/v1/sys-user/' + userId,
    method: 'get'
  })
}

// 删除用户
export function delUser(data) {
  return service({
    url: '/v1/sys-user',
    method: 'delete',
    data: data
  })
}

// 新增用户
export function addUser(data) {
  return service({
    url: '/v1/sys-user',
    method: 'post',
    data: data
  })
}

// 修改用户
export function updateUser(data) {
  return service({
    url: '/v1/sys-user',
    method: 'put',
    data: data
  })
}

// 用户密码重置
export function resetUserPwd(userId, password) {
  const data = {
    userId,
    password
  }
  return service({
    url: '/v1/user/pwd/reset',
    method: 'put',
    data: data
  })
}

// @Tags User
// @Summary 设置用户权限
export const setUserAuthority = (data) => {
  return service({
    url: '/user/setUserAuthority',
    method: 'post',
    data: data
  })
}

// @Tags SysUser
// @Summary 删除用户
export const deleteUser = (data) => {
  return service({
    url: '/user/deleteUser',
    method: 'delete',
    data: data
  })
}

// @Tags SysUser
// @Summary 设置用户信息
export const setUserInfo = (data) => {
  return service({
    url: '/user/setUserInfo',
    method: 'put',
    data: data
  })
}

// @Tags SysUser
// @Summary 设置用户信息
export const setSelfInfo = (data) => {
  return service({
    url: '/v1/sys-user',
    method: 'put',
    data: data
  })
}

export const updateSelfPhone = (data) => {
  return service({
    url: '/v1/sys-user/updateSelfPhone',
    method: 'put',
    data: data
  })
}

export const updateSelfNickName = (data) => {
  return service({
    url: '/v1/sys-user/updateSelfNickName',
    method: 'put',
    data: data
  })
}

export const updateSelfEmail = (data) => {
  return service({
    url: '/v1/sys-user/updateSelfEmail',
    method: 'put',
    data: data
  })
}

// @Tags User
// @Summary 设置用户权限
export const setUserAuthorities = (data) => {
  return service({
    url: '/user/setUserAuthorities',
    method: 'post',
    data: data
  })
}

// @Tags User
// @Summary 获取用户信息
export const getUserInfo = () => {
  return service({
    url: '/v1/user/profile',
    method: 'get'
  })
}

export const resetPassword = (data) => {
  return service({
    url: '/user/resetPassword',
    method: 'post',
    data: data
  })
}
