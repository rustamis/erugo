export default class SettingsCollection {
  static settings = []

  constructor(settings) {
    this.settings = settings
  }

  addSetting(setting) {
    this.settings.push(setting)
  }

  updateSetting(key, value) {
    const setting = this.getSetting(key)
    if (setting) {
      setting.value = value + ''
    }
  }

  getSetting(key) {
    return this.settings.find((setting) => setting.key === key)
  }

  getSettings() {
    return this.settings
  }

  getSettingValue(key) {
    const setting = this.getSetting(key)
    return setting ? setting.value : null
  }

  getSettingValueAsType(key) {
    const settingValue = this.getSettingValue(key)
    //figure out the type of the setting
    if (settingValue === 'true' || settingValue === 'false') {
      return settingValue === 'true'
    }
    //check if the setting is a number
    if (!isNaN(settingValue)) {
      return parseInt(settingValue)
    }
    //check if the setting is a float
    if (!isNaN(parseFloat(settingValue))) {
      return parseFloat(settingValue)
    }
    //if it's not a number, return the string
    return settingValue
  }
}
