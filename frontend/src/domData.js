export const domData = () => {
  const data = {
    version: document.body.getAttribute('data-version'),
    ...getSettings()
  }
  console.log(data)
  return data
}

const getSettings = () => {
  //   <body data-version="v0.0.1" data-setting-application_name="eurgo File Transfers" data-setting-css_accent_color="#1BC7B1" data-setting-css_accent_color_light="#E7E7E7" data-setting-css_primary_color="#FFB243" data-setting-css_secondary_color="#222200" data-setting-logo_width="100">
  const settings = {}
  //get any data-setting-<key>="value" attributes
  const settingAttributes = document.body.attributes
  for (let i = 0; i < settingAttributes.length; i++) {
    const attribute = settingAttributes[i]
    if (attribute.name.startsWith('data-setting-')) {
      settings[attribute.name.replace('data-setting-', '')] = attribute.value
    }
  }
  return settings
}