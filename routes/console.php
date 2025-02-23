<?php

use Illuminate\Support\Facades\Schedule;
use App\Jobs\cleanExpiredShares;

Schedule::job(cleanExpiredShares::class)->daily();

//command to manually run the job
Artisan::command('shares:clean-expired', function () {
    cleanExpiredShares::dispatch();
})->purpose('Clean expired shares');
