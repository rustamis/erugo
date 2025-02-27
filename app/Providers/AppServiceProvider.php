<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use App\Services\SettingsService;
use Illuminate\Support\Facades\View;
use Illuminate\Support\Facades\Schema;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        //
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {


        try {
            $settingsService = new SettingsService();
            $settings = $settingsService->getGlobalViewData();
            View::share('settings', $settings);
        } catch (\Exception $e) {
            //do nothing
        }



        View::prependLocation(storage_path('templates'));
    }
}
