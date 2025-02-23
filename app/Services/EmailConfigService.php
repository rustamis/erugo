<?php

namespace App\Services;

use Illuminate\Support\Facades\Mail;
use Illuminate\Mail\Mailable;

class EmailConfigService
{

  private $settingService;
  private $settings;

  public function __construct()
  {
    $this->settingService = new SettingsService();
    $this->settings = $this->settingService->getMailSettings();
  }

  public function configureMailer()
  {
    $settings = $this->settings;

    $config = [
      'transport' => 'smtp',
      'host' => $settings['host'],
      'port' => $settings['port'],
      'encryption' => $settings['encryption'],
      'username' => $settings['username'],
      'password' => $settings['password'],
      'timeout' => null,
    ];

    $mailer = Mail::build($config);

    return $mailer;
  }

  public function configureMailable(Mailable $mailable)
  {
    $mailable->from($this->settings['from_address'], $this->settings['from_name']);
    return $mailable;
  }
}
