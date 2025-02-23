<?php

namespace App\Jobs;

use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use App\Services\EmailConfigService;
use Illuminate\Mail\Mailable;
class sendEmail implements ShouldQueue
{
    use Queueable;
    public $to;
    public $mailableClass;
    public $mailableData;
    /**
     * Create a new job instance.
     */
    public function __construct(String $to, String $mailableClass, $mailableData = [])
    {
        $this->to = $to;
        $this->mailableClass = $mailableClass;
        $this->mailableData = $mailableData;
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        $emailService = new EmailConfigService();
        $mailer = $emailService->configureMailer();
        $mailable = new $this->mailableClass(...$this->mailableData);
        $mailable = $emailService->configureMailable($mailable);
        $mailer->to($this->to)->send($mailable);
    }
}
