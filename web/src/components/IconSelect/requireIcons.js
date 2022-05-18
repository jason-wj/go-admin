
const svgs = import.meta.globEager('./icons/svg/*.svg')
// const requireAll = requireContext => requireContext.keys()

// const re = /\.\/(.*)\.svg/

/* const icons = svgs.map(i => {
  return i.match(re)[1]
})*/

const icons = []

Object.keys(svgs).forEach((key) => {
  icons.push(key.replace(/(\.\/icons\/svg\/|\.svg)/g, ''))
})

export default icons
