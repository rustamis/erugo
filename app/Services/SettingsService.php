<?php

namespace App\Services;

use App\Models\Setting;

class SettingsService
{
  public function get($key)
  {
    $setting = Setting::where('key', $key)->first();
    if (!$setting) {
      return null;
    }
    return $setting->value;
  }

  public function getMaxUploadSize()
  {
    $max_upload_size = $this->get('max_share_size');
    $max_upload_size_unit = $this->get('max_share_size_unit');
    if ($max_upload_size_unit && $max_upload_size) {
      if ($max_upload_size_unit == 'MB') {
        return $max_upload_size * 1024 * 1024;
      } else if ($max_upload_size_unit == 'GB') {
        return $max_upload_size * 1024 * 1024 * 1024;
      }
    }
    return null;
  }
}
