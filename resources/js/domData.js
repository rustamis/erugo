export const domData = () => {
  const data = {
    version: document.body.getAttribute('data-version'),
    ...getSettings()
  }
  return data
}

const getSettings = () => {
  const body = document.body
  if(!body) {
    return {}
  }
  const settings = JSON.parse(body.getAttribute('data-settings'))
  return settings
}