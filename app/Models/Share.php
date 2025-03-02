<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

use App\Models\Setting;
use App\Mail\shareDeletedWarningMail;
use App\Jobs\sendEmail;
use App\Services\SettingsService;
class Share extends Model
{

  protected $fillable = [
    'user_id',
    'name',
    'description',
    'path',
    'long_id',
    'size',
    'file_count',
    'download_limit',
    'download_count',
    'require_email',
    'expires_at',
    'status'
  ];

  protected $casts = [
    'expires_at' => 'datetime',
    'deletes_at' => 'datetime',
  ];

  protected $appends = [
    'expired',
    'deletes_at',
    'deleted'
  ];

  protected $hidden = [
    'path',
    'user_id',
  ];

  public function files()
  {
    return $this->hasMany(File::class);
  }

  public function user()
  {
    return $this->belongsTo(User::class);
  }

  function getExpiredAttribute()
  {
    return $this->expires_at < now()->addMinutes(1);
  }

  function getDeletesAtAttribute()
  {
    $cleanFilesAfterDays = Setting::where('key', 'clean_files_after_days')->first();

    if (!$cleanFilesAfterDays) {
      $cleanFilesAfterDays = 30;
    } else {
      $cleanFilesAfterDays = (int) $cleanFilesAfterDays->value;
    }

    if ($cleanFilesAfterDays == 0) {
      $cleanFilesAfterDays = 30;
    }

    return $this->expires_at->addDays($cleanFilesAfterDays);
  }

  function getDeletedAttribute()
  {
    return $this->status == 'deleted';
  }

  public function scopeReadyForCleaning($query)
  {
    $cleanFilesAfterDays = Setting::where('key', 'clean_files_after_days')->first();

    if (!$cleanFilesAfterDays) {
      $cleanFilesAfterDays = 30;
    } else {
      $cleanFilesAfterDays = (int) $cleanFilesAfterDays->value;
    }

    if ($cleanFilesAfterDays == 0) {
      $cleanFilesAfterDays = 30;
    }

    $deletesAt = now()->subDays($cleanFilesAfterDays);
    \Log::info('Cleaning shares after ' . $deletesAt);

    return $query->where('expires_at', '<', $deletesAt)->where('status', '!=', 'deleted');
  }


  public function cleanFiles()
  {
    try {
      $filePath = $this->path;
      //is the path a directory?
      if (is_dir($filePath)) {
        //delete all files in the directory
        $files = glob($filePath . '/*');
        foreach ($files as $file) {
          unlink($file);
        }
        //delete the directory
        rmdir($filePath);
      }
      //or is it a zip file?
      if (is_file($filePath . '.zip')) {
        unlink($filePath . '.zip');
      }

      $this->status = 'deleted';
      $this->save();

      $settingsService = new SettingsService();

      $shouldSend = $settingsService->get('emails_share_deleted_enabled') ?? true;

      if ($shouldSend) {
        sendEmail::dispatch($this->user->email, shareDeletedWarningMail::class, ['share' => $this]);
      }

      return true;
    } catch (\Exception $e) {
      Log::error('Error cleaning files for share ' . $this->id . ': ' . $e->getMessage());
      return false;
    }
  }
}
